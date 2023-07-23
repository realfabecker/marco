package mws

import "encoding/xml"

// BrowseList mws browse node report list
type BrowseList struct {
	XMLName xml.Name     `xml:"Result"`
	Result  []BrowseNode `xml:"Node"`
}

// BrowseNode mws browse node definition
type BrowseNode struct {
	BrowseNodeId           string `xml:"browseNodeId"`
	BrowseNodeName         string `xml:"browseNodeName"`
	BrowsePathById         string `xml:"browsePathById"`
	BrowsePathByName       string `xml:"browsePathByName"`
	HasChildren            bool   `xml:"hasChildren"`
	ProductTypeDefinitions string `xml:"productTypeDefinitions"`
}

// ProductTypeDefinitions download specs
type ProductTypeDefinitions struct {
	Schema struct {
		Link struct {
			Resource string `json:"resource"`
		} `json:"link"`
	} `json:"schema"`
}
