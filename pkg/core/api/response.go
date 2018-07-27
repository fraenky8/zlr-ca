// jsend, https://labs.omniti.com/labs/jsend
package api

import (
	"fmt"

	"github.com/fraenky8/zlr-ca/pkg/core/domain"
)

const (
	// All went well, and (usually) some data was returned.
	StatusOk = "success"
	// There was a problem with the data submitted, or some pre-condition of the API call wasn't satisfied
	StatusFail = "fail"
	// An error occurred in processing the request, i.e. an exception was thrown
	StatusError = "error"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type IcecreamResponse struct {
	Icecream *domain.Icecream `json:"icecream"`
}

type IcecreamsResponse struct {
	Icecreams []*domain.Icecream `json:"icecreams"`
}

type ErrorsResponse struct {
	Error []string `json:"errors"`
}

func Success(data interface{}) Response {
	return Response{
		Status: StatusOk,
		Data:   data,
	}
}

func Fail(errors ...error) Response {

	var errs []string
	for _, e := range errors {
		errs = append(errs, e.Error())
	}

	return Response{
		Status: StatusFail,
		Data:   ErrorsResponse{errs},
	}
}

func FailString(strings ...string) Response {
	var errors []error
	for _, s := range strings {
		errors = append(errors, fmt.Errorf(s))
	}
	return Fail(errors...)
}

func Error(message string) Response {
	return Response{
		Status:  StatusError,
		Message: message,
	}
}
