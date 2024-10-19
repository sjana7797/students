package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

const (
	StatusOk    = "OK"
	StatusError = "Error"
)

func WriteJSON(w http.ResponseWriter, status int, data interface{}) error {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

func GeneralError(err error) Response {
	return Response{
		Status: StatusError,
		Error:  err.Error(),
	}
}

func ValidatorError(errs validator.ValidationErrors) Response {
	var errMgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMgs = append(errMgs, fmt.Sprintf("field %s is required field", err.Field()))

		default:
			errMgs = append(errMgs, fmt.Sprintf("field %s is invalid feild", err.Field()))
		}
	}

	return Response{
		Status: StatusError,
		Error:  strings.Join(errMgs, ","),
	}
}
