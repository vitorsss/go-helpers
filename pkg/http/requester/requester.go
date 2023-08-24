package requester

import (
	"net/http"

	"golang.org/x/time/rate"
)

type Requester interface {
	Do(req *http.Request) (*http.Response, error)
}

type rateLimitedRequester struct {
	requester   Requester
	rateLimiter *rate.Limiter
}

func (c *rateLimitedRequester) Do(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	err := c.rateLimiter.Wait(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := c.requester.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func NewRateLimitedRequester(
	requester Requester,
	rateLimiter *rate.Limiter,
) Requester {
	c := &rateLimitedRequester{
		requester:   requester,
		rateLimiter: rateLimiter,
	}
	return c
}
