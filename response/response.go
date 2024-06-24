package response

import (
	"encoding/json"
	"net/http"
)

type ResponseData struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SuccessResponse(w http.ResponseWriter, code int, data interface{}) {
	response := &ResponseData{
		Code:    code,
		Message: msg[code],
		Data:    data,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func ErrorResponse(w http.ResponseWriter, Code int) {
	response := &ResponseData{
		Code:    Code,
		Message: msg[Code],
		Data:    nil,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
