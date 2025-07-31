package main

import (
	"log"
	"net/http"
	"xmr_mixer/connectors"
	"xmr_mixer/signin"
)

// CurrencyProvider определяет интерфейс для получения списка валют
type CurrencyProvider interface {
	GetCurrencies() ([]byte, error)
}

type AddressBookProvider interface {
	GetAddressBook() ([]byte, error)
}

type GenericLoginHandler interface {
	Login(w http.ResponseWriter, r *http.Request)
}

type GenericRegisterHandler interface {
	Register(w http.ResponseWriter, r *http.Request)
}

func avoidCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		h(w, r)
	}
}

func main() {

	var crypto_getter CurrencyProvider = &connectors.DexConnector{}

	var genericLogin GenericLoginHandler = &signin.GenericController{}
	var genericRegister GenericRegisterHandler = &signin.GenericController{}

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
		w.Write([]byte(`{"addressbook": "This is a placeholder for address book data"}`))
	}))

	http.HandleFunc("/api/login", avoidCORS(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		genericLogin.Login(w, r)
	}))

	http.HandleFunc("/api/register", avoidCORS(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		genericRegister.Register(w, r)
	}))

	log.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
