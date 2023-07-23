package mws

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sync/errgroup"
)

// BrowseService
type BrowseService struct{}

// Read
func (b *BrowseService) Read(p string) (*BrowseList, error) {
	file, err := os.Open(p)
	if err != nil {
		return nil, fmt.Errorf("ls: %w", err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("rd: %w", err)
	}

	var l BrowseList
	if err := xml.Unmarshal(data, &l); err != nil {
		return nil, fmt.Errorf("um: %w", err)
	}
	return &l, nil
}

// Flat
func (b *BrowseService) Flat(l *BrowseList, target string) error {

	var d = make(map[string]BrowseList)
	for _, v := range l.Result {
		s := strings.Split(v.BrowsePathById, ",")
		if _, ok := d[s[0]]; !ok {
			d[s[0]] = BrowseList{Result: []BrowseNode{v}}
		} else if ok {
			d[s[0]] = BrowseList{Result: append(d[s[0]].Result, v)}
		}
	}

	p, err := os.OpenFile(
		filepath.Join(target, "nodes.csv"),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return err
	}
	defer p.Close()

	var eg errgroup.Group
	eg.SetLimit(len(d))

	for i, v := range d {
		eg.Go(func(k string, s BrowseList) func() error {
			return func() error {
				log.Printf("Writing %s...\n", k)

				d, err := xml.MarshalIndent(s, "", "  ")
				if err != nil {
					return fmt.Errorf("xml:%w", err)
				}

				if err := os.WriteFile(
					filepath.Join(target, k+".xml"),
					d,
					0644,
				); err != nil {
					return fmt.Errorf("xml:%w", err)
				}

				t := s.Result[0].BrowseNodeName + " (" + k + ")"
				if _, err := p.WriteString(t + "\n"); err != nil {
					return fmt.Errorf("xml:%w", err)
				}
				return nil
			}
		}(i, v))
	}
	return eg.Wait()
}
