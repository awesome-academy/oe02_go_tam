package utils

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

func BuildVnpUrl(params url.Values, secret, baseUrl, hashType string) string {
	params.Del("vnp_SecureHash")
	params.Set("vnp_SecureHashType", hashType)

	var keys []string
	for k := range params {
		if k != "vnp_SecureHash" && k != "vnp_SecureHashType" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	var rawData strings.Builder
	for i, k := range keys {
		if i > 0 {
			rawData.WriteString("&")
		}
		rawData.WriteString(k + "=" + vnpayEscape(params.Get(k)))
	}

	h := hmac.New(sha512.New, []byte(secret))
	h.Write([]byte(rawData.String()))
	signature := strings.ToUpper(hex.EncodeToString(h.Sum(nil)))

	fullUrl := fmt.Sprintf("%s?%s&vnp_SecureHash=%s", baseUrl, rawData.String(), signature)

	return fullUrl
}

func VerifyVnpSignature(query url.Values, secret string) bool {
	var keys []string
	for k := range query {
		if k != "vnp_SecureHash" && k != "vnp_SecureHashType" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	var data []string
	for _, k := range keys {
		data = append(data, k+"="+vnpayEscape(query.Get(k)))
	}
	rawData := strings.Join(data, "&")

	h := hmac.New(sha512.New, []byte(secret))
	h.Write([]byte(rawData))
	expected := strings.ToUpper(hex.EncodeToString(h.Sum(nil)))

	received := strings.ToUpper(query.Get("vnp_SecureHash"))

	return expected == received
}

func vnpayEscape(s string) string {
	return url.QueryEscape(s)
}
