package connectors

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"os"
)

func PostCurrencies() ([]byte, error) {
	url := "https://ff.io/api/v2/ccies"
	client := &http.Client{}

	apiKey := os.Getenv("FIXEDFLOAT_API_KEY")
	apiSecret := os.Getenv("FIXEDFLOAT_API_SECRET")

	jsonData := []byte("{}")

	// HMAC SHA256 подпись
	h := hmac.New(sha256.New, []byte(apiSecret))
	h.Write(jsonData)
	sign := hex.EncodeToString(h.Sum(nil))

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("X-API-KEY", apiKey)
	request.Header.Set("X-API-SIGN", sign)

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
