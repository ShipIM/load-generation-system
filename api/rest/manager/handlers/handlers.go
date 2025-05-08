package handlers

import (
	"load-generation-system/api/rest/manager/handlers/model"
	"load-generation-system/api/rest/manager/presenters/attack"
	"load-generation-system/pkg/web"

	"github.com/gofiber/fiber/v2"
)

// @Title  Start new attack
// @Param  config  body  model.StartAttackRequestBody  true  "Attack configuration"
// @Success  201  object  model.StartAttackResponse  "Successful attack start"
// @Failure  400  object  model.BadRequestError  "Bad request error"
// @Failure  422  object  model.ValidationResponse  "Validation error"
// @Failure  500  object  model.InternalServerError  "Internal server error"
// @Resource  Attack
// @Router  /manager/api/v1/attacks [post]
func (r *Resolver) startAttack(ctx *fiber.Ctx) error {
	var presenter attack.StartAttackPresenter
	status, errResp := r.bodyChecker(ctx, &presenter)
	if errResp != nil {
		return ctx.Status(status).JSON(errResp)
	}

	start, err := presenter.ToCore()
	if err != nil {
		response, status := model.MapError(err)
		return ctx.Status(status).JSON(response)
	}

	attackDetails, err := r.attackService.StartAttack(start)
	if err != nil {
		response, status := model.MapError(err)
		return ctx.Status(status).JSON(response)
	}

	pres := attack.PresentAttack(attackDetails)

	resp := web.OKResponse(pres)
	return ctx.Status(fiber.StatusCreated).JSON(resp)
}

// @Title  Stop attack
// @Param  attack_id  path  int64  true  "Attack id"  "1"
// @Success  200  object  model.StopAttackResponse  "Successful attack end"
// @Failure  400  object  model.BadRequestError  "Bad request error"
// @Failure  404  object  model.NotFoundError  "Not found error"
// @Failure  500  object  model.InternalServerError  "Internal server error"
// @Resource  Attack
// @Router  /manager/api/v1/attacks/{attack_id} [delete]
func (r *Resolver) stopAttack(ctx *fiber.Ctx) error {
	id, err := parseInt64Param(ctx, "attack_id")
	if err != nil {
		response, status := model.MapError(err)
		return ctx.Status(status).JSON(response)
	}

	err = r.attackService.StopAttack(id)
	if err != nil {
		response, status := model.MapError(err)
		return ctx.Status(status).JSON(response)
	}

	resp := web.OKResponse(nil)
	return ctx.Status(fiber.StatusOK).JSON(resp)
}

// @Title  Get attack scenarios
// @Success  200  object  model.GetScenariosResponse  "Successful get scenarios"
// @Failure  500  object  model.InternalServerError  "Internal server error"
// @Resource  Attack
// @Router  /manager/api/v1/scenarios [get]
func (r *Resolver) getScenarios(ctx *fiber.Ctx) error {
	scenarios := r.attackService.GetScenarios()

	pres := attack.PresentScenarioList(scenarios)

	resp := web.OKResponse(pres)
	return ctx.Status(fiber.StatusOK).JSON(resp)
}

// @Title  Get active attacks
// @Success  200  object  model.GetAttacksResponse  "Successful get attacks"
// @Failure  500  object  model.InternalServerError "Internal server error"
// @Resource  Attack
// @Router  /manager/api/v1/attacks [get]
func (r *Resolver) getAttacks(ctx *fiber.Ctx) error {
	attacks := r.attackService.GetAttacks()

	pres := attack.PresentAttackList(attacks)

	resp := web.OKResponse(pres)
	return ctx.Status(fiber.StatusOK).JSON(resp)
}

// @Title  Get nodes information
// @Success  200  object  model.GetNodesResponse  "Successful get nodes"
// @Failure  500  object  model.InternalServerError  "Internal server error"
// @Resource  Attack
// @Router  /manager/api/v1/nodes [get]
func (r *Resolver) getNodes(ctx *fiber.Ctx) error {
	nodes := r.attackService.ListNodes()

	pres := attack.PresentNodeList(nodes)

	resp := web.OKResponse(pres)
	return ctx.Status(fiber.StatusOK).JSON(resp)
}

// @Title  Start new increment
// @Param  attack_id  path  int64  true  "Attack id"  "1"
// @Param  config  body  model.StartIncrementRequestBody  true  "Increment configuration"
// @Success 201  object  model.StartIncrementResponse  "Successful increment start"
// @Failure 400  object  model.BadRequestError  "Bad request error"
// @Failure 422  object  model.ValidationResponse  "Validation error"
// @Failure 500  object  model.InternalServerError  "Internal server error"
// @Resource  Attack
// @Router  /manager/api/v1/attacks/{attack_id}/increments [post]
func (r *Resolver) startIncrement(ctx *fiber.Ctx) error {
	id, err := parseInt64Param(ctx, "attack_id")
	if err != nil {
		response, status := model.MapError(err)
		return ctx.Status(status).JSON(response)
	}

	var presenter attack.StartIncrementPresenter
	status, errResp := r.bodyChecker(ctx, &presenter)
	if errResp != nil {
		return ctx.Status(status).JSON(errResp)
	}

	start := presenter.ToCore(id)

	incrementDetails, err := r.attackService.StartIncrement(start)
	if err != nil {
		response, status := model.MapError(err)
		return ctx.Status(status).JSON(response)
	}

	pres := attack.PresentIncrement(incrementDetails)

	resp := web.OKResponse(pres)
	return ctx.Status(fiber.StatusCreated).JSON(resp)
}

// @Title  Stop increment
// @Param  attack_id  path  int64  true  "Attack id"  "1"
// @Param  increment_id  path  int64  true  "Increment id"  "1"
// @Success  200  object  model.StopIncrementResponse  "Successful increment stop"
// @Failure  400  object  model.BadRequestError  "Bad request error"
// @Failure  404  object  model.NotFoundError  "Not found error"
// @Failure  500  object  model.InternalServerError  "Internal server error"
// @Resource  Attack
// @Router  /manager/api/v1/attacks/{attack_id}/increments/{increment_id} [delete]
func (r *Resolver) stopIncrement(ctx *fiber.Ctx) error {
	attackID, err := parseInt64Param(ctx, "attack_id")
	if err != nil {
		response, status := model.MapError(err)
		return ctx.Status(status).JSON(response)
	}

	incrementID, err := parseInt64Param(ctx, "increment_id")
	if err != nil {
		response, status := model.MapError(err)
		return ctx.Status(status).JSON(response)
	}

	err = r.attackService.StopIncrement(attackID, incrementID)
	if err != nil {
		response, status := model.MapError(err)
		return ctx.Status(status).JSON(response)
	}

	resp := web.OKResponse(nil)
	return ctx.Status(fiber.StatusOK).JSON(resp)
}
