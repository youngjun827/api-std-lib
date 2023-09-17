package error_response

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *ErrorResponse) JsonErrorResponse(w http.ResponseWriter, err error, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	errorResponse := ErrorResponse{
		Code:    code,
		Message: err.Error(),
	}

	json.NewEncoder(w).Encode(errorResponse)
}
