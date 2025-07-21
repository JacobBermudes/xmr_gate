package main

import (
	"log"
	"net/http"
	"xmr_mixer/connectors"
)

// CurrencyProvider определяет интерфейс для получения списка валют
type CurrencyProvider interface {
	GetCurrencies() ([]byte, error)
}

func avoidCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		h(w, r)
	}
}

func main() {

	// DexConnector как CurrencyProvider
	var provider CurrencyProvider = &connectors.DexConnector{}

	http.HandleFunc("/api/v1/currencies", avoidCORS(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp, err := provider.GetCurrencies()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`[]`))
			return
		}
		w.Write(resp)
	}))

	log.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
