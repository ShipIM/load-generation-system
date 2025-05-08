package scenarios

import (
	"context"
	"load-generation-system/internal/service/callers"
)

const (
	testHTTP = "test_http"
)

var (
	testHTTPScen Scenario = New(
		testHTTP,
		"test http",
		func(ctx context.Context, caller *callers.Caller) error {
			return caller.TestCaller.Test(ctx)
		},
	)
)
