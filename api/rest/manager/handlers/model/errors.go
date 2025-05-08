package model

import (
	"errors"
	"load-generation-system/pkg/web"
)

var (
	ErrValidation            = errors.New("error when validation")
	ErrRequestBodyIsRequired = errors.New("request body is required")
	ErrParseBody             = errors.New("error when body parsing")
	ErrParseQuery            = errors.New("error when query parsing")
	ErrInvalidPathParam      = errors.New("error when path param parsing")
)

type ValidationErrors struct {
	Errors []web.ValidationError `json:"errors"`
}

type ValidationResponse struct {
	Status string           `json:"status" example:"ERROR"`
	Data   ValidationErrors `json:"data" example:"{\"errors\": [{\"tag\": \"<tag>\", \"field\": \"<field>\", \"param\": \"<param>\"}]}"`
}

type InternalServerError struct {
	Status string `json:"status" example:"ERROR"`
}

type BadRequestError struct {
	Status string `json:"status" example:"ERROR"`
}

type NotFoundError struct {
	Status string `json:"status" example:"ERROR"`
}
