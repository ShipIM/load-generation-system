package handlers

import (
	"load-generation-system/internal/core"
	"load-generation-system/pkg/rest"

	"github.com/go-playground/validator/v10"
)

type Resolver struct {
	server        rest.Server
	attackService core.AttackService
	validate      *validator.Validate
}

const (
	serviceName = "manager"
	APIPrefix   = "api"
	APIVersion  = "v1"
	pathPrefix  = serviceName + "/" + APIPrefix + "/" + APIVersion
)

func NewResolver(
	server rest.Server,
	attackService core.AttackService,
) *Resolver {
	resolver := &Resolver{
		server:        server,
		attackService: attackService,
		validate:      newValidate(),
	}

	resolver.initRoutes()

	return resolver
}

func (r *Resolver) initRoutes() {
	r.server.Router().Post(pathPrefix+"/attacks", r.startAttack)
	r.server.Router().Post(pathPrefix+"/attacks/:attack_id/increments", r.startIncrement)
	r.server.Router().Delete(pathPrefix+"/attacks/:attack_id", r.stopAttack)
	r.server.Router().Delete(pathPrefix+"/attacks/:attack_id/increments/:increment_id", r.stopIncrement)
	r.server.Router().Get(pathPrefix+"/scenarios", r.getScenarios)
	r.server.Router().Get(pathPrefix+"/attacks", r.getAttacks)
	r.server.Router().Get(pathPrefix+"/nodes", r.getNodes)
}
