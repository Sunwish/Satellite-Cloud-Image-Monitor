package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

var (
	xmlURL      = flag.String("xmlURL", "http://img.nsmc.org.cn/CLOUDIMAGE/FY4B/AGRI/GCLR/SEC/xml/FY4B-china-72h.xml", "XML URL")
	checkCount  = flag.Int("checkCount", 5, "XML image check count")
	dateFormat  = flag.String("dateFormat", "2006/20060102/", "Date directory path format")
	outputDir   = flag.String("outputDir", "./archived", "Output directory path")
	filterReg   = flag.String("filterReg", ".*thumb.*", "Image name filter regular expression")
	updateDelay = 15 * time.Minute

	// Pushdeer notify configuration
	notifyBaseUrl = flag.String("notifyBaseUrl", "", "Pushdeer Notify Base URL")
	notifyKey     = flag.String("notifyKey", "", "Pushdeer Notify Key")
	notifyPrefix  = flag.String("notifyPrefix", "[CloudMonitor]", "Pushdeer notify title prefix")
)

func main() {
	fmt.Println("Parsing flags...")

	flag.Parse()

	// 编译正则表达式
	regex, err := regexp.Compile(*filterReg)
	regEnable := true
	if err != nil {
		regEnable = false
		fmt.Println("Error compiling filter regular expression, regex filter disabled. Error: ", err)
		notify(*notifyBaseUrl, *notifyKey, *notifyPrefix+" Error compiling filter regular expression", err.Error())
	}

	// 检查Pushdeer推送配置
	if notifyBaseUrl != nil && *notifyBaseUrl != "" && notifyKey != nil && *notifyKey != "" {
		fmt.Println("Pushdeer notify configuration is enabled. Test notification sent.")
		notify(*notifyBaseUrl, *notifyKey, *notifyPrefix+" Startup successful.", "SCIM container is running.")
	} else {
		fmt.Println("Pushdeer notify configuration is disabled.")
	}

	fmt.Println("Running...")

	// 获取云图列表，并记录checkCount范围内未下载的图像
	for {
		imgList, err := getImageList(*xmlURL)
		if err != nil {
			fmt.Println("Error getting or parsing XML, will retry after 10s. Error: ", err)
			notify(*notifyBaseUrl, *notifyKey, *notifyPrefix+" Error parsing XML", err.Error())
			time.Sleep(10 * time.Second)
			continue
		}

		downloadInfos := make([]DownloadInfo, 0)

		checkedCount := 0
		for i := 0; i < len(imgList.Images) && checkedCount < *checkCount; i++ {
			image := imgList.Images[i]

			// 按正则过滤特定名称图像
			fileName := path.Base(image.URL)
			if regEnable && regex.MatchString(fileName) {
				continue
			}

			// 检查图像是否已下载
			localDir := path.Join(*outputDir, strings.Join(strings.Split(imgList.Name, " "), "-"), parseXMLTimeToTimePath(image.Time))
			localFilePath := path.Join(localDir, fileName)

			// 使用os.Stat检查文件是否存在
			_, err := os.Stat(localFilePath)

			// 若本地图像不存在则下载图像
			if os.IsNotExist(err) {
				downloadInfos = append(downloadInfos, DownloadInfo{image.URL, localDir})
			}

			checkedCount++
		}

		// 下载checkCount范围内未下载的图像
		fmt.Printf("Found %d new cloud image(s).\n", len(downloadInfos))
		for _, downloadInfo := range downloadInfos {
			go func(info DownloadInfo) {
				fmt.Printf("Request %s\n", info.URL)
				err := downloadImage(info.URL, info.LocalDir)
				if err != nil {
					fmt.Println("Error downloading image:", err)
					notify(*notifyBaseUrl, *notifyKey, *notifyPrefix+" Error downloading image", "URL: "+info.URL+"\n\n"+"ERR: "+err.Error())
				}
			}(downloadInfo)
		}

		// 打印信息
		fmt.Print("Monitor is sleeping")
		if len(downloadInfos) > 0 {
			fmt.Print(" and starting download new image(s)")
		}
		fmt.Println(".")

		// 等待下次更新
		time.Sleep(updateDelay)
	}
}
