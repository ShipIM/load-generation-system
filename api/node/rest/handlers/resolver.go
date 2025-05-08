package handlers

import "load-generation-system/pkg/rest"

type Resolver struct {
	server rest.Server
}

func NewResolver(
	server rest.Server,
) *Resolver {
	resolver := &Resolver{
		server: server,
	}

	return resolver
}
