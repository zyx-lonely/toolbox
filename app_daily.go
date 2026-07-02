package main

import (
	"pc-toolbox/internal/daily"
)

// ============================================================
//  翻译工具
// ============================================================

// Translate 翻译文本
func (a *App) Translate(text string, fromLang string, toLang string) (*daily.TranslateResult, error) {
	return daily.Translate(text, fromLang, toLang)
}

// GetClipboardAndTranslate 获取剪贴板内容并翻译
func (a *App) GetClipboardAndTranslate(toLang string) (*daily.TranslateResult, error) {
	return daily.GetClipboardAndTranslate(toLang)
}

// ============================================================
//  汇率查询
// ============================================================

// GetSupportedCurrencies 获取支持的货币列表
func (a *App) GetSupportedCurrencies() []daily.CurrencyInfo {
	return daily.GetSupportedCurrencies()
}

// ConvertCurrency 货币转换
func (a *App) ConvertCurrency(amount float64, from string, to string) (*daily.ExchangeRate, error) {
	return daily.ConvertCurrency(amount, from, to)
}

// GetExchangeRate 获取汇率
func (a *App) GetExchangeRate(from string, to string) (float64, error) {
	return daily.GetExchangeRate(from, to)
}

// ============================================================
//  科学计算器
// ============================================================

// Calculate 计算表达式
func (a *App) Calculate(expr string) (*daily.CalcResult, error) {
	return daily.Calculate(expr)
}
