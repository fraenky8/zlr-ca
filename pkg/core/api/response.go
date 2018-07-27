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

type IngredientResponse struct {
	Ingredient domain.Ingredients `json:"ingredients"`
}

type IngredientsResponse struct {
	Ingredients []domain.Ingredients `json:"ingredients"`
}

type SourcingValueResponse struct {
	SourcingValue domain.SourcingValues `json:"sourcing_values"`
}

type SourcingValuesResponse struct {
	SourcingValues []domain.SourcingValues `json:"sourcing_values"`
}

type ErrorsResponse struct {
	Error []string `json:"errors"`
}

func SuccessResponse(data interface{}) Response {
	return Response{
		Status: StatusOk,
		Data:   data,
	}
}

func FailResponse(errors ...error) Response {

	var errs []string
	for _, e := range errors {
		errs = append(errs, e.Error())
	}

	return Response{
		Status: StatusFail,
		Data:   ErrorsResponse{errs},
	}
}

func FailStringResponse(strings ...string) Response {
	var errors []error
	for _, s := range strings {
		errors = append(errors, fmt.Errorf(s))
	}
	return FailResponse(errors...)
}

func ErrorResponse(message string) Response {
	return Response{
		Status:  StatusError,
		Message: message,
	}
}
