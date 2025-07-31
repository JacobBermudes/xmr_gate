package signin

import (
	"encoding/json"
	"net/http"

	"github.com/dgraph-io/badger/v4"
)

var db *badger.DB

func init() {
	var err error
	db, err = badger.Open(badger.DefaultOptions("./badgerdb"))
	if err != nil {
		panic(err)
	}
}

type GenericController struct{}
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (a *GenericController) Login(w http.ResponseWriter, r *http.Request) {

	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Unvalid body scheme!", http.StatusBadRequest)
		return
	}
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(u.Username))
		if err != nil {
			return err
		}
		var storedPassword []byte
		if err := item.Value(func(val []byte) error {
			storedPassword = val
			return nil
		}); err != nil {
			return err
		}
		if string(storedPassword) != u.Password {
			return badger.ErrKeyNotFound
		}
		return nil
	})
	if err != nil {
		if err == badger.ErrKeyNotFound {
			http.Error(w, "Uncorrect login or email!", http.StatusUnauthorized)
			return
		}
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful"))
}

func (a *GenericController) Register(w http.ResponseWriter, r *http.Request) {

	var u User

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	err := db.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(u.Username))
		return err
	})
	if err == nil {
		http.Error(w, "user exists", http.StatusConflict)
		return
	}
	err = db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(u.Username), []byte(u.Password))
	})
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"status":"registered"}`))
}
