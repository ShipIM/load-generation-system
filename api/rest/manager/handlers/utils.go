package handlers

import (
	"load-generation-system/api/rest/manager/handlers/model"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func parseInt64Param(ctx *fiber.Ctx, param string) (int64, error) {
	idStr := ctx.Params(param)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("%s: %s - %s", model.ErrInvalidPathParam.Error(), param, err.Error())
		return 0, model.ErrInvalidPathParam
	}
	return id, nil
}
