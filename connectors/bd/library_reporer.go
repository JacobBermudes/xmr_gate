package library_reporter

import (
	"encoding/json"
)

type Libreporter struct{}

func (d *Libreporter) GetAddressBook() ([]byte, error) {
	return json.Marshal("cryptoObj")
}
