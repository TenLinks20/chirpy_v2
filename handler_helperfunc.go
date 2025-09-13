package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func respondWithErr(w http.ResponseWriter, code int, msg string)  {
	type apiError struct {
		Error string `json:"error"`
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	apiErr := apiError{
		Error: msg,
	}
	if err := json.NewEncoder(w).Encode(&apiErr); err != nil {
		fmt.Println("Failed to encode JSON", err)
	}
}

func respondWithJSON(w http.ResponseWriter, code int, payload any)  {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		fmt.Println("Failed to encode JSON", err)
	}
}