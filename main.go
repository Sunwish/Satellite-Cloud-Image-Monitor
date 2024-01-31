package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"
)

var (
	baseURL        = flag.String("baseURL", "https://img.nsmc.org.cn/CLOUDIMAGE/FY4B/AGRI/GCLR/SEC/china/", "Image Base URL")
	dateFormat     = flag.String("dateFormat", "2006/20060102/", "Date directory path format")
	fileNamePrefix = flag.String("fileNamePrefix", "FY4B_china_", "File name prefix")
	fileNameFormat = flag.String("fileNameFormat", "20060102150400.JPG", "File name format")
	outputDir      = flag.String("outputDir", "./archived", "Output directory path")
	updateDelay    = 15 * time.Minute

	// Pushdeer notify configuration
	notifyBaseUrl = flag.String("notifyBaseUrl", "", "Pushdeer Notify Base URL")
	notifyKey     = flag.String("notifyKey", "", "Pushdeer Notify Key")
)

func main() {

	fmt.Println("Parsing flags...")

	flag.Parse()

	if notifyBaseUrl != nil && *notifyBaseUrl != "" && notifyKey != nil && *notifyKey != "" {
		fmt.Println("Pushdeer notify configuration is enabled. Test notification sent.")
		notify(*notifyBaseUrl, *notifyKey, "[CloudMonitor] Startup successful.", "SCIM container is running.")
	} else {
		fmt.Println("Pushdeer notify configuration is disabled.")
	}

	fmt.Println("Running...")

	for {
		// 获取当前最近的15分钟时间
		currentTimeUTC := time.Now().UTC()
		roundedTime := time.Date(
			currentTimeUTC.Year(),
			currentTimeUTC.Month(),
			currentTimeUTC.Day(),
			currentTimeUTC.Hour(),
			currentTimeUTC.Minute()/15*15, // 取最近的15分钟时间
			0,
			0,
			time.UTC,
		)
		oneHour30MinutesAgo := roundedTime.Add(-1*time.Hour - 30*time.Minute)

		// 获取一个半小时前至一个小时45分钟前范围内的最新云图（不获取当前最新的，以防更新不及时导致获取失败）
		imageURL := buildImageURL(oneHour30MinutesAgo)

		fmt.Println("Request " + imageURL)

		err := downloadImage(imageURL)
		if err != nil {
			fmt.Println("Error downloading image:", err)
			notify(*notifyBaseUrl, *notifyKey, "[CloudMonitor] Error downloading image", err.Error())
		}

		// Wait for the next update
		fmt.Println("Sleeping...")
		time.Sleep(updateDelay)
	}
}

func buildImageURL(t time.Time) string {
	datePath := t.Format(*dateFormat)
	fileName := *fileNamePrefix + t.Format(*fileNameFormat)
	return *baseURL + datePath + fileName
}

func downloadImage(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
	}

	// Create the "archived" directory if it doesn't exist
	if err := os.MkdirAll(*outputDir, os.ModePerm); err != nil {
		return err
	}

	// Create the file with the same name as in the URL
	filePath := path.Join(*outputDir, path.Base(url))
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
