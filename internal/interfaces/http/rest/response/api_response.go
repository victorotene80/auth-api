package response

import (
	"encoding/json"
	"net/http"
)

type APIResponse[T any] struct {
	Status  string `json:"status"`            
	Code    string `json:"code,omitempty"`    
	Message string `json:"message,omitempty"` 
	Data    *T     `json:"data,omitempty"`    
	Errors  any    `json:"errors,omitempty"`  
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func WriteJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}

func Success[T any](w http.ResponseWriter, statusCode int, code, message string, data *T) {
	resp := APIResponse[T]{
		Status:  "success",
		Code:    code,
		Message: message,
		Data:    data,
	}
	WriteJSON(w, statusCode, resp)
}

func Error(w http.ResponseWriter, statusCode int, code, message string, errors any) {
	resp := APIResponse[struct{}]{ // empty data
		Status:  "error",
		Code:    code,
		Message: message,
		Errors:  errors,
	}
	WriteJSON(w, statusCode, resp)
}