package common

import (
	"encoding/json"
	"net/http"
	"time"
)

type Meta struct {
	Timestamp string `json:"timestamp"`
}

type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Meta    Meta        `json:"meta"`
}

func getMeta() Meta {
	return Meta{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

func WriteJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func WriteSuccess(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := SuccessResponse{
		Success: true,
		Data:    data,
		Meta:    getMeta(),
	}

	json.NewEncoder(w).Encode(response)
}
