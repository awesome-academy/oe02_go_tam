package constant

import (
	"encoding/json"
	"os"
)

var translations map[string]string

func LoadI18n(locale string) error {
	file, err := os.Open("locales/" + locale + ".json")
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(&translations)
}

func T(key string) string {
	if val, ok := translations[key]; ok {
		return val
	}
	return key
}
