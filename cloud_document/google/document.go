package document

import (
	"fmt"
	"net/url"
	"regexp"
)

// 转换 Google Sheets URL 为 CSV 导出 URL
func ConvertToCSVURL(sheetURL string) (string, error) {
	parsedURL, err := url.Parse(sheetURL)
	if err != nil {
		return "", fmt.Errorf("解析 URL 失败: %v", err)
	}

	// 正则匹配 Google Sheets ID
	re := regexp.MustCompile(`/d/([a-zA-Z0-9-_]+)`)
	matches := re.FindStringSubmatch(sheetURL)
	if len(matches) < 2 {
		return "", fmt.Errorf("未找到 Google Sheets ID")
	}
	sheetID := matches[1]

	// 提取 gid
	gid := "0" // 默认工作表 ID
	queryParams := parsedURL.Query()
	if val, exists := queryParams["gid"]; exists && len(val) > 0 {
		gid = val[0]
	}

	// 构造 CSV 导出 URL
	csvURL := fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s/gviz/tq?tqx=out:csv&gid=%s", sheetID, gid)
	return csvURL, nil
}
