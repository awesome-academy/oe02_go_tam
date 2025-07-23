package config

import "os"

type VnpayConfig struct {
	TmnCode    string
	HashSecret string
	PayURL     string
	ReturnURL  string
	HashType   string
}

func GetVnpayConfig() VnpayConfig {
	return VnpayConfig{
		TmnCode:    os.Getenv("VNP_TMNCODE"),
		HashSecret: os.Getenv("VNP_HASHSECRET"),
		PayURL:     os.Getenv("VNP_URL"),
		ReturnURL:  os.Getenv("VNP_RETURN_URL"),
		HashType:   os.Getenv("VNP_HASH_TYPE"),
	}
}
