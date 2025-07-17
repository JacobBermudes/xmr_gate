package main

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime"
	"xmr_mixer/connectors"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Разрешить любые Origin (НЕ использовать в продакшене!)
	},
}

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
	http.HandleFunc("/api/test/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"answer": true})
	})
	http.HandleFunc("/ws", wsHandler)

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

func wsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Origin: %s", r.Header.Get("Origin"))
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Ошибка апгрейда:", err)
		return
	}
	go handleWS(conn)
}

func handleWS(conn *websocket.Conn) {
	defer conn.Close()
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Ошибка чтения:", err)
			break
		}
		log.Printf("Получено сообщение: %s", msg)

		conn.WriteMessage(websocket.TextMessage, []byte("pong"))
		log.Printf("Текущее количество горутин: %d", runtime.NumGoroutine())
	}
}
