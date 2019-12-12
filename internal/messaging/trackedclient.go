package messaging

import (
	"context"
	phttp "github.com/beatlabs/patron/trace/http"
	"net/http"
	"time"
)

// HTTPClient is an http client which adds prometheus metrics to each call
type HTTPClient struct {
	provider string
	cl       phttp.Client
}

// Do performs an HTTP request and adds prometheus metrics
func (tcl *HTTPClient) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	url := ""
	if req.URL != nil {
		url = req.URL.String()
	}

	start := time.Now()
	rsp, err := tcl.cl.Do(ctx, req)
	if err != nil {
		ObserveCount(tcl.provider, url, req.Method, 0)
		ObserveLatency(tcl.provider, url, req.Method, 0, time.Since(start))
		return rsp, err
	}

	ObserveCount(tcl.provider, url, req.Method, rsp.StatusCode)
	ObserveLatency(tcl.provider, url, req.Method, rsp.StatusCode, time.Since(start))

	return rsp, err
}

// NewHTTPClient returns a new client
func NewHTTPClient(cl phttp.Client, pr string) *HTTPClient {
	return &HTTPClient{cl: cl, provider: pr}
}
