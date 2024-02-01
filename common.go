package main

import (
	"encoding/xml"
)

// 定义 XML 结构体
type ImageList struct {
	XMLName xml.Name `xml:"imagelist"`
	Name    string   `xml:"name,attr"`
	Tag     string   `xml:"tag,attr"`
	Images  []Image  `xml:"image"`
}

type Image struct {
	Time string `xml:"time,attr"`
	Desc string `xml:"desc,attr"`
	URL  string `xml:"url,attr"`
}

type DownloadInfo struct {
	URL      string
	LocalDir string
}
