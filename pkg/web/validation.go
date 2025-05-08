package web

import (
	"github.com/go-playground/validator/v10"
)

type (
	ErrorPayload struct {
		Reason string `json:"reason"`
	}

	ValidationError struct {
		Tag   string `json:"tag"`
		Field string `json:"field"`
		Param string `json:"param"`
	}
)

func ValidationErrors(validatorErrors validator.ValidationErrors) []ValidationError {
	errors := make([]ValidationError, len(validatorErrors))

	for i, err := range validatorErrors {
		errors[i] = ValidationError{
			Tag:   err.ActualTag(),
			Field: err.Field(),
			Param: err.Param(),
		}
	}

	return errors
}

func ValidationErrorResponse(errors []ValidationError) Response {
	return Response{
		Status: ERROR,
		Data: map[string]any{
			"errors": errors,
		},
	}
}
