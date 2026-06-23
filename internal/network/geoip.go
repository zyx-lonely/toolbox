package network

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// GeoIPResult IP 地理位置查询结果
type GeoIPResult struct {
	IP         string  `json:"ip"`
	Country    string  `json:"country"`
	Region     string  `json:"region"`
	City       string  `json:"city"`
	ISP        string  `json:"isp"`
	Org        string  `json:"org"`
	Latitude   float64 `json:"lat"`
	Longitude  float64 `json:"lon"`
	Success    bool    `json:"success"`
	Error      string  `json:"error,omitempty"`
}

type ipAPIData struct {
	Status   string  `json:"status"`
	Country  string  `json:"country"`
	Region   string  `json:"region"`
	City     string  `json:"city"`
	ISP      string  `json:"isp"`
	Org      string  `json:"org"`
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
	Query    string  `json:"query"`
	Message  string  `json:"message"`
}

// QueryGeoIP 查询 IP 地理位置（使用 ip-api.com 免费 API）
func QueryGeoIP(ip string) GeoIPResult {
	if ip == "" {
		return GeoIPResult{Success: false, Error: "IP 地址为空"}
	}

	client := &http.Client{Timeout: 5 * time.Second}
	url := fmt.Sprintf("https://ip-api.com/json/%s?fields=status,message,country,region,city,isp,org,lat,lon,query", ip)

	resp, err := client.Get(url)
	if err != nil {
		return GeoIPResult{IP: ip, Success: false, Error: fmt.Sprintf("请求失败: %v", err)}
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var data ipAPIData
	if err := json.Unmarshal(body, &data); err != nil {
		return GeoIPResult{IP: ip, Success: false, Error: fmt.Sprintf("解析失败: %v", err)}
	}

	if data.Status == "fail" {
		return GeoIPResult{IP: ip, Success: false, Error: data.Message}
	}

	return GeoIPResult{
		IP:        data.Query,
		Country:   data.Country,
		Region:    data.Region,
		City:      data.City,
		ISP:       data.ISP,
		Org:       data.Org,
		Latitude:  data.Lat,
		Longitude: data.Lon,
		Success:   true,
	}
}
