package main

import (
	"log"
	"net/http"
	"xmr_mixer/connectors"
	library_reporter "xmr_mixer/connectors/bd"
)

// CurrencyProvider определяет интерфейс для получения списка валют
type CurrencyProvider interface {
	GetCurrencies() ([]byte, error)
}

type AddressBookProvider interface {
	GetAddressBook() ([]byte, error)
}

func avoidCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		h(w, r)
	}
}

func main() {

	var crypto_getter CurrencyProvider = &connectors.DexConnector{}
	var addressBookReader AddressBookProvider = &library_reporter.Libreporter{}

	http.HandleFunc("/api/v1/currencies", avoidCORS(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp, err := crypto_getter.GetCurrencies()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`[]`))
			return
		}
		w.Write(resp)
	}))

	http.HandleFunc("/api/v1/addressbook", avoidCORS(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp, err := addressBookReader.GetAddressBook()
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
