package daily

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// ExchangeRate 汇率信息
type ExchangeRate struct {
	From       string  `json:"from"`
	To         string  `json:"to"`
	Rate       float64 `json:"rate"`
	Amount     float64 `json:"amount"`
	Result     float64 `json:"result"`
	UpdateTime string  `json:"updateTime"`
}

// CurrencyInfo 货币信息
type CurrencyInfo struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// GetSupportedCurrencies 获取支持的货币列表
func GetSupportedCurrencies() []CurrencyInfo {
	return []CurrencyInfo{
		{Code: "CNY", Name: "人民币 (CNY)"},
		{Code: "USD", Name: "美元 (USD)"},
		{Code: "EUR", Name: "欧元 (EUR)"},
		{Code: "GBP", Name: "英镑 (GBP)"},
		{Code: "JPY", Name: "日元 (JPY)"},
		{Code: "KRW", Name: "韩元 (KRW)"},
		{Code: "HKD", Name: "港币 (HKD)"},
		{Code: "TWD", Name: "新台币 (TWD)"},
		{Code: "SGD", Name: "新加坡元 (SGD)"},
		{Code: "AUD", Name: "澳元 (AUD)"},
		{Code: "CAD", Name: "加元 (CAD)"},
		{Code: "CHF", Name: "瑞士法郎 (CHF)"},
		{Code: "THB", Name: "泰铢 (THB)"},
		{Code: "MYR", Name: "马来西亚林吉特 (MYR)"},
		{Code: "RUB", Name: "俄罗斯卢布 (RUB)"},
	}
}

// ConvertCurrency 货币转换
func ConvertCurrency(amount float64, from string, to string) (*ExchangeRate, error) {
	from = strings.ToUpper(from)
	to = strings.ToUpper(to)

	if from == to {
		return &ExchangeRate{
			From:       from,
			To:         to,
			Rate:       1,
			Amount:     amount,
			Result:     amount,
			UpdateTime: time.Now().Format("2006-01-02 15:04:05"),
		}, nil
	}

	rate, err := getExchangeRate(from, to)
	if err != nil {
		return nil, err
	}

	return &ExchangeRate{
		From:       from,
		To:         to,
		Rate:       rate,
		Amount:     amount,
		Result:     amount * rate,
		UpdateTime: time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}

func getExchangeRate(from string, to string) (float64, error) {
	// 使用免费汇率 API
	apiURL := fmt.Sprintf("https://api.exchangerate-api.com/v4/latest/%s", from)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(apiURL)
	if err != nil {
		return 0, fmt.Errorf("获取汇率失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("读取响应失败: %w", err)
	}

	var result struct {
		Rates map[string]float64 `json:"rates"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return 0, fmt.Errorf("解析汇率数据失败")
	}

	rate, ok := result.Rates[to]
	if !ok {
		return 0, fmt.Errorf("不支持的货币: %s", to)
	}

	return rate, nil
}

// GetExchangeRate 获取汇率（不转换）
func GetExchangeRate(from string, to string) (float64, error) {
	return getExchangeRate(strings.ToUpper(from), strings.ToUpper(to))
}
