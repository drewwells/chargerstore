package api

import (
	"encoding/json"
	"net/http"
)

func marshal(w http.ResponseWriter, v interface{}) error {
	bs, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err

	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bs)
	return nil
}
