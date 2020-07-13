package api

import (
	"encoding/json"
	"net/http"
)

func jsonError(w http.ResponseWriter, err error, status int) {
	w.WriteHeader(status)

	resp := Response{Success: false}
	resp.Errors = []string{err.Error()}

	js, _ := json.Marshal(resp)

	_, _ = w.Write(js)

	return
}