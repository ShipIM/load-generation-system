package model

import (
	"errors"
	"load-generation-system/internal/core"
	"load-generation-system/pkg/web"

	"github.com/gofiber/fiber/v2"
)

func MapError(err error) (resp web.Response, status int) {
	switch {
	case errors.Is(err, ErrRequestBodyIsRequired):
		return web.ErrorResponse(
			web.ErrorPayload{
				Reason: err.Error(),
			},
		), fiber.StatusBadRequest
	case errors.Is(err, ErrParseBody):
		return web.ErrorResponse(
			web.ErrorPayload{
				Reason: err.Error(),
			},
		), fiber.StatusBadRequest
	case errors.Is(err, ErrInvalidPathParam):
		return web.ErrorResponse(
			web.ErrorPayload{
				Reason: err.Error(),
			},
		), fiber.StatusBadRequest
	case errors.Is(err, core.ErrAttackNotFound):
		return web.ErrorResponse(
			web.ErrorPayload{
				Reason: err.Error(),
			},
		), fiber.StatusNotFound
	case errors.Is(err, core.ErrIncrementNotFound):
		return web.ErrorResponse(
			web.ErrorPayload{
				Reason: err.Error(),
			},
		), fiber.StatusNotFound
	case errors.Is(err, core.ErrScenarioNotFound):
		return web.ErrorResponse(
			web.ErrorPayload{
				Reason: err.Error(),
			},
		), fiber.StatusNotFound
	case errors.Is(err, core.ErrEmptyAttack):
		return web.ErrorResponse(
			web.ErrorPayload{
				Reason: err.Error(),
			},
		), fiber.StatusBadRequest
	case errors.Is(err, core.ErrBrokenScheduler):
		return web.ErrorResponse(
			web.ErrorPayload{
				Reason: err.Error(),
			},
		), fiber.StatusInternalServerError
	case errors.Is(err, core.ErrBadConfig):
		return web.ErrorResponse(
			web.ErrorPayload{
				Reason: err.Error(),
			},
		), fiber.StatusBadRequest
	}

	return web.ErrorResponse(nil), fiber.StatusInternalServerError
}
