package test

import (
	"context"
)

type TestCaller interface {
	Test(ctx context.Context) error
}

const (
	host    = "localhost:8090"
	path    = "/test/api"
	version = "/v1"

	protocol = "http://"
)
