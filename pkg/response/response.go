package response

import (
	"encoding/json"
	"net/http"
)

type errorBody struct {
	Error string `json:"error"`
}

func Error(w http.ResponseWriter, statusCode int, msg string) {
	w.WriteHeader(statusCode)
	body, _ := json.Marshal(errorBody{msg})
	w.Write(body)
}
