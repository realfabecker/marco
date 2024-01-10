package mws

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rafaelbeecker/mwskit/internal/mws/signer"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type ReportService struct{}

func (r *ReportService) CreateReport(reportType string) (string, error) {
	body := map[string]interface{}{
		"marketplaceIds": []string{"A2Q3Y263D00KWC"},
		"reportType":     reportType,
	}
	data, err := json.Marshal(&body)
	if err != nil {
		return "", fmt.Errorf("marshal: %w", err)
	}

	url := `https://sellingpartnerapi-na.amazon.com/reports/2021-06-30/reports`
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("GetProductTypeDefSchemaUrl: %w", err)
	}

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
	if err != nil {
		return "", fmt.Errorf("CreateReport: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println(err)
		}
	}()
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("CreateReport: %w", err)
	}
	if resp.StatusCode != 202 {
		return "", fmt.Errorf("CreateReport: %d", resp.StatusCode)
	}
	var rep struct {
		ReportId string `json:"reportId"`
	}
	if err := json.Unmarshal(data, &rep); err != nil {
		return "", err
	}
	return rep.ReportId, nil
}

func (r *ReportService) GetReport(reportId string) (*Report, error) {
	url := `https://sellingpartnerapi-na.amazon.com/reports/2021-06-30/reports/` + reportId
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("GetProductTypeDefSchemaUrl: %w", err)
	}

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
	if err != nil {
		return nil, fmt.Errorf("CreateReport: %w", err)
	}
	defer resp.Body.Close()
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("CreateReport: %w", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("CreateReport: %d", resp.StatusCode)
	}
	var rep Report
	if err := json.Unmarshal(d, &rep); err != nil {
		return nil, fmt.Errorf("GetReport: %w", err)
	}
	return &rep, nil
}

func (r *ReportService) GetReportDocument(reportDocumentId string) (*ReportDocument, error) {
	url := `https://sellingpartnerapi-na.amazon.com//reports/2021-06-30/documents/` + reportDocumentId
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("GetProductTypeDefSchemaUrl: %w", err)
	}

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
	if err != nil {
		return nil, fmt.Errorf("CreateReport: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println(err)
		}
	}()
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("CreateReport: %w", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("CreateReport: %d", resp.StatusCode)
	}
	var rep ReportDocument
	if err := json.Unmarshal(d, &rep); err != nil {
		return nil, fmt.Errorf("GetReport: %w", err)
	}
	return &rep, nil
}

func (r *ReportService) DownloadReportDocument(name string, dest string, url string) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("GetProductTypeDefSchemaUrl: %w", err)
	}
	file, err := os.Create(filepath.Join(dest, name+".gz"))
	if err != nil {
		return fmt.Errorf("DownloadReportDocument: %w", err)
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("DownloadReportDocument: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println(err)
		}
	}()
	if _, err := io.Copy(file, resp.Body); err != nil {
		return fmt.Errorf("DownloadReportDocument: %w", err)
	}
	return nil
}
