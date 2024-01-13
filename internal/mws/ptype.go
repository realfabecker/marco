package mws

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rafaelbeecker/mwskit/internal/mws/signer"
	"golang.org/x/sync/errgroup"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type PtypeService struct{}

func (s *PtypeService) GetProductTypeDefSchemaUrl(sellerId string, productType string) (string, error) {
	url := `https://sellingpartnerapi-na.amazon.com/definitions/2020-09-01/productTypes/` + productType
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("GetProductTypeDefSchemaUrl (%s): %w", productType, err)
	}

	q := req.URL.Query()
	q.Add("sellerId", sellerId)
	q.Add("marketplaceIds", "A2Q3Y263D00KWC")
	q.Add("requirements", "LISTING")
	q.Add("locale", "pt_BR")
	req.URL.RawQuery = q.Encode()

	req.Header.Set("host", "sellingpartnerapi-na.amazon.com")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("x-amz-access-token", os.Getenv("AWS_ACCESS_TOKEN"))
	req.Header.Set("user-agent", "App 1.0 (Language=Golang/1.18);")

	req2 := signer.Sign4(req, signer.Credentials{
		AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		SecurityToken:   os.Getenv("AWS_SESSION_TOKEN"),
		Region:          "us-east-1",
		Service:         "execute-api",
	})

	client := http.Client{}
	resp, err := client.Do(req2)
	if resp != nil && resp.StatusCode == 429 {
		return "", Err429
	} else if resp != nil && resp.StatusCode == 404 {
		return "", Err404
	} else if err != nil {
		return "", fmt.Errorf("GetProductTypeDefSchemaUrl (%s): %w", productType, err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println(err)
		}
	}()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("GetProductTypeDefSchemaUrl (%s): %w", productType, err)
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("GetProductTypeDefSchemaUrl (%s): %d", productType, resp.StatusCode)
	}

	payload := ProductTypeDefinitions{}
	if err := json.Unmarshal(data, &payload); err != nil {
		return "", fmt.Errorf("GetProductTypeDefSchemaUrl (%s): %w", productType, err)
	}
	return payload.Schema.Link.Resource, nil
}

func (s *PtypeService) DownloadProductTypeDef(dest string, link string) error {
	request, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		return fmt.Errorf("DownloadProductTypeDef: %w", err)
	}

	file, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("DownloadProductTypeDef: %w", err)
	}

	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("DownloadProductTypeDef: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println(err)
		}
	}()

	if _, err := io.Copy(file, resp.Body); err != nil {
		return fmt.Errorf("DownloadProductTypeDef: %w", err)
	}
	return nil
}

func (s *PtypeService) DownloadBatchTypeDef(marketplace string, productList string, target string) error {
	file, err := os.Open(productList)
	if err != nil {
		return err
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Println(err)
		}
	}()

	reader := csv.NewReader(file)
	data, err := reader.ReadAll()
	if err != nil {
		return err
	}

	var eg errgroup.Group
	eg.SetLimit(50)

	var retry time.Time
	var mtx sync.Mutex
	for i, v := range data {
		eg.Go(func(t string, i int) func() error {
			return func() error {
				mtx.Lock()
				r := retry
				mtx.Unlock()

				if time.Now().Before(r) {
					time.Sleep(time.Until(r))
				}

				dest := filepath.Join(target, t+".json")
				f, err := os.Stat(dest)
				if f != nil {
					return nil
				} else if !errors.Is(err, os.ErrNotExist) {
					return err
				}

				log.Printf("downloading schema %s\n", t)
				link, err := s.GetProductTypeDefSchemaUrl(marketplace, t)
				if err != nil && errors.Is(err, Err429) {
					log.Printf("%s: 429\n", t)
					mtx.Lock()
					retry = time.Now().Add(time.Second * 10)
					mtx.Unlock()
					return nil
				} else if err != nil && errors.Is(err, Err404) {
					log.Printf("%s: 404\n", t)
					return nil
				} else if err != nil {
					return err
				}

				if err := s.DownloadProductTypeDef(dest, link); err != nil {
					return err
				}

				log.Printf("schema downloaded at %s\n", dest)
				return nil
			}
		}(v[0], i))
	}
	return eg.Wait()
}
