package utils

import "strings"

func ParseCurrencies(currencyStr string) []string{
	currencies := strings.Split(currencyStr, ",")
	for i := range currencies{
		curClean := strings.TrimSpace(currencies[i])
		currencies[i] = curClean
	}
	return currencies
}