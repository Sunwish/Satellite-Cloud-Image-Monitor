package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
)

func getImageList(url string) (ImageList, error) {
	// 创建结构体实例
	var imageList ImageList

	// 发送 HTTP 请求
	resp, err := http.Get(url)
	if err != nil {
		return imageList, fmt.Errorf("Error fetching XML: %s", err.Error())
	}
	defer resp.Body.Close()

	// 读取 HTTP 响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return imageList, fmt.Errorf("Error reading response body: %s", err)
	}

	// 解析 XML
	err = xml.Unmarshal(body, &imageList)
	if err != nil {
		return imageList, fmt.Errorf("Error parsing XML: %s", err)
	}

	return imageList, nil
}

func downloadImage(url string, dir string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
	}

	// Create the "archived" directory if it doesn't exist
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	// Create the file with the same name as in the URL
	filePath := path.Join(dir, path.Base(url))
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Copy the content of the response body to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("Downloaded %s\n", filePath)
	return nil
}

func notify(baseUrl string, key string, title string, content string) {
	if baseUrl == "" {
		fmt.Println("消息推送失败：notifyBaseUrl为空")
		return
	}

	if key == "" {
		fmt.Println("消息推送失败：notifyKey为空")
		return
	}

	// 发起GET请求
	fullUrl := baseUrl + "push?pushkey=" + key + "&text=" + url.QueryEscape(title) + "&desp=" + url.QueryEscape(content) + "&type=markdown"
	response, err := http.Get(fullUrl)
	if err != nil {
		fmt.Println("消息推送失败：", err)
		return
	}
	defer response.Body.Close()
}
