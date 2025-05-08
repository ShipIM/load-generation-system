package handlers

import (
	"errors"
	"load-generation-system/api/rest/manager/handlers/model"
	"load-generation-system/pkg/web"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func newValidate() *validator.Validate {
	return validator.New(validator.WithRequiredStructEnabled())
}

func (r *Resolver) validateStruct(entity any) (int, *web.Response) {
	if err := r.validate.Struct(entity); err != nil {
		var fieldErrors validator.ValidationErrors
		ok := errors.As(err, &fieldErrors)
		if !ok {
			log.Println(model.ErrValidation.Error())
			response, status := model.MapError(model.ErrValidation)
			return status, &response
		}

		response := web.ValidationErrorResponse(web.ValidationErrors(fieldErrors))

		return http.StatusUnprocessableEntity, &response
	}

	return 0, nil
}

func (r *Resolver) bodyChecker(ctx *fiber.Ctx, entity any) (int, *web.Response) {
	if err := ctx.BodyParser(entity); err != nil {
		log.Println(model.ErrParseBody.Error(), "error", err.Error())
		response, status := model.MapError(model.ErrParseBody)
		return status, &response
	}

	return r.validateStruct(entity)
}
