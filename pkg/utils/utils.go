package utils

import (
	"crypto/sha1"
	"fmt"
	"regexp"
	"stock_scraper/types"
	"strconv"
	"strings"
)

func Hash(str string) string {
	h := sha1.New()
	h.Write([]byte(str))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func ParsePrice(str string) *types.Price {
	s := strings.Replace(str, ",", "", -1)
	r := regexp.MustCompile("([$]?)([0-9]+(\\.[0-9]{2})?)")
	match := r.FindStringSubmatch(s)

	if len(match) > 2 {
		symbol := "$"
		if match[1] != "" {
			symbol = match[1]
		}

		if match[2] != "" {
			if value, err := strconv.ParseFloat(match[2], 64); err == nil {
				return &types.Price{Symbol: symbol, Value: value}
			}
		}
	}

	return nil
}
