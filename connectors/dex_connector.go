package connectors

import (
	"encoding/json"
	"log"
)

// DexConnector реализует работу с внешними DEX API

type DexConnector struct{}

func (d *DexConnector) GetFixedFloatCurrencies() ([]byte, error) {
	resp, err := PostCurrencies()
	if err != nil {
		log.Println("Ошибка при запросе FixedFloat:", err)
		return nil, err
	}
	return resp, nil
}

func (d *DexConnector) GetCurrencies() ([]byte, error) {
	cryptoResp, err := d.GetFixedFloatCurrencies()
	if err != nil {
		return nil, err
	}

	var cryptoObj map[string]any
	if err := json.Unmarshal(cryptoResp, &cryptoObj); err != nil {
		return nil, err
	}

	return json.Marshal(cryptoObj)
}
