package http

import (
	"github.com/go-resty/resty/v2"
	"log"
	"time"
)

// New returns a Resty client with default settings
func New() *resty.Client {
	return configure(resty.New())
}

// NewWithBaseURL returns a Resty client with a given base URL
func NewWithBaseURL(baseURL string) *resty.Client {
	client := resty.New().SetBaseURL(baseURL)
	return configure(client)
}

// Internal: common config for both
func configure(client *resty.Client) *resty.Client {
	client.
		// Timeouts
		SetTimeout(10*time.Second).
		// Retry logic
		SetRetryCount(3).
		SetRetryWaitTime(1*time.Second).
		SetRetryMaxWaitTime(5*time.Second).
		// Logging requests/responses
		OnBeforeRequest(func(c *resty.Client, r *resty.Request) error {
			log.Printf("[HTTP] --> %s %s", r.Method, r.URL)
			return nil
		}).
		OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
			log.Printf("[HTTP] <-- %d %s", resp.StatusCode(), resp.Request.URL)
			return nil
		}).
		// JSON content-type by default
		SetHeader("Content-Type", "application/json")

	return client
}
