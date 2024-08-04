package response

import (
	"encoding/json"
	"net/http"
)

type errorBody struct {
	Error string `json:"error"`
}

func Error(w http.ResponseWriter, statusCode int, msg string) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(statusCode)
	body, _ := json.Marshal(errorBody{msg})
	w.Write(body)
}

func Json(w http.ResponseWriter, data any) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(200)
	body, _ := json.Marshal(data)
	w.Write(body)
}
