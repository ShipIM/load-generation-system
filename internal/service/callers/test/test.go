package test

import (
	"context"
	"load-generation-system/internal/core"
)

type testCaller struct {
	urlBase    string
	httpClient core.Client
}

func NewCaller(
	httpClient core.Client,
) TestCaller {
	return &testCaller{
		urlBase:    protocol + host + path,
		httpClient: httpClient,
	}
}

func (c *testCaller) Test(ctx context.Context) error {
	_, err := c.httpClient.R().
		SetPath(c.urlBase + version + "/test").
		Get(ctx)

	return err
}
