package mws

import (
	"encoding/xml"
	"errors"
)

var (
	Err429 = errors.New("oh no, too many requests")
	Err404 = errors.New("oh no, record not found")
)

type BrowseList struct {
	XMLName xml.Name     `xml:"Result"`
	Result  []BrowseNode `xml:"Node"`
}

type BrowseNode struct {
	BrowseNodeId           string `xml:"browseNodeId"`
	BrowseNodeName         string `xml:"browseNodeName"`
	BrowsePathById         string `xml:"browsePathById"`
	BrowsePathByName       string `xml:"browsePathByName"`
	HasChildren            bool   `xml:"hasChildren"`
	ProductTypeDefinitions string `xml:"productTypeDefinitions"`
}

type ProductTypeDefinitions struct {
	Schema struct {
		Link struct {
			Resource string `json:"resource"`
		} `json:"link"`
	} `json:"schema"`
}

type Report struct {
	ReportType       string   `json:"reportType"`
	ProcessingStatus string   `json:"processingStatus"`
	MarketplaceIds   []string `json:"MarketplaceIds"`
	ReportId         string   `json:"reportId"`
	ReportDocumentId string   `json:"reportDocumentId"`
}

type ReportDocument struct {
	ReportDocumentId     string `json:"reportDocumentId"`
	CompressionAlgorithm string `json:"compressionAlgorithm"`
	Url                  string `json:"url"`
}
