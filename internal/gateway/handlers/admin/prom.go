package admin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// promClient is a minimal Prometheus HTTP API query client.
// No gRPC, no Prometheus Go client library — just net/http + encoding/json.
type promClient struct {
	addr   string       // e.g. "http://localhost:9090"
	client *http.Client // timeout baked in
}

func newPromClient(addr string, timeout time.Duration) *promClient {
	return &promClient{
		addr:   addr,
		client: &http.Client{Timeout: timeout},
	}
}

// promResult represents a single vector result from Prometheus.
type promResult struct {
	Metric map[string]string  `json:"metric"`
	Value  [2]json.RawMessage `json:"value"` // [timestamp, "value_string"]
}

// promQueryResponse is the top-level Prometheus /api/v1/query response.
type promQueryResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string       `json:"resultType"`
		Result     []promResult `json:"result"`
	} `json:"data"`
}

// query executes an instant PromQL query and returns the result vector.
func (p *promClient) query(ctx context.Context, expr string) ([]promResult, error) {
	u := fmt.Sprintf("%s/api/v1/query?query=%s", p.addr, url.QueryEscape(expr))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("prometheus returned %d", resp.StatusCode)
	}

	var out promQueryResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	if out.Status != "success" {
		return nil, fmt.Errorf("prometheus query status: %s", out.Status)
	}
	return out.Data.Result, nil
}

// reachable returns true if the Prometheus /-/ready endpoint responds 200.
func (p *promClient) reachable(ctx context.Context) bool {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.addr+"/-/ready", nil)
	if err != nil {
		return false
	}
	resp, err := p.client.Do(req)
	if err != nil {
		return false
	}
	resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}
