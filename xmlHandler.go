package main

import (
	"fmt"
	"time"
)

func parseXMLTimeToTimePath(xmlTime string) string {
	t, err := parseXMLTime(xmlTime)

	if err != nil {
		return ""
	}

	return parseTimeToTimePath(t)
}

func parseXMLTime(timeStr string) (time.Time, error) {
	// 定义时间格式
	timeFormat := "2006-01-02 15:04 (MST)"

	// 将字符串解析为时间对象
	parsedTime, err := time.Parse(timeFormat, timeStr)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return parsedTime, err
	}

	return parsedTime, nil
}

func parseTimeToTimePath(t time.Time) string {
	return t.Format(*dateFormat)
}
