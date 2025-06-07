package responseutil

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func JSON(response http.ResponseWriter, code int, status, message string, data interface{}) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(code)

	json.NewEncoder(response).Encode(Response{
		Status:  status,
		Message: message,
		Data:    data,
	})
}

func Success(response http.ResponseWriter, message string, data interface{}) {
	JSON(response, http.StatusOK, "success", message, data)
}

func Error(response http.ResponseWriter, code int, message string) {
	JSON(response, code, "error", message, nil)
}
