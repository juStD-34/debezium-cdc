package http

import (
	"fmt"
	"go.uber.org/zap"
	"register/pkg/logger"
	"time"

	"github.com/go-resty/resty/v2"
)

type HTTPClient interface {
	Get(url string, result interface{}) error
	Post(url string, body interface{}, result interface{}) error
	Put(url string, body interface{}, result interface{}) error
	Delete(url string) error
}

type RestyClient struct {
	client *resty.Client
	logger logger.Logger
}

func NewRestyClient(logger logger.Logger) HTTPClient {
	client := resty.New()
	client.SetTimeout(30 * time.Second)
	client.SetRetryCount(3)
	client.SetRetryWaitTime(5 * time.Second)

	// Add logging middleware
	client.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
		logger.Info("HTTP Request",
			zap.String("method", req.Method),
			zap.String("url", req.URL),
		)
		return nil
	})

	client.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
		logger.Info("HTTP Response",
			zap.String("method", resp.Request.Method),
			zap.String("url", resp.Request.URL),
			zap.Int("status", resp.StatusCode()),
		)
		return nil
	})

	return &RestyClient{
		client: client,
		logger: logger,
	}
}

func (r *RestyClient) Get(url string, result interface{}) error {
	resp, err := r.client.R().
		SetResult(result).
		Get(url)

	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("HTTP error: %d - %s", resp.StatusCode(), resp.String())
	}

	return nil
}

func (r *RestyClient) Post(url string, body interface{}, result interface{}) error {
	resp, err := r.client.R().
		SetBody(body).
		SetResult(result).
		Post(url)

	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("HTTP error: %d - %s", resp.StatusCode(), resp.String())
	}

	return nil
}

func (r *RestyClient) Put(url string, body interface{}, result interface{}) error {
	resp, err := r.client.R().
		SetBody(body).
		SetResult(result).
		Put(url)

	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("HTTP error: %d - %s", resp.StatusCode(), resp.String())
	}

	return nil
}

func (r *RestyClient) Delete(url string) error {
	resp, err := r.client.R().Delete(url)
	if err != nil {
		return err
	}

	if resp.IsError() && resp.StatusCode() != 404 {
		return fmt.Errorf("HTTP error: %d - %s", resp.StatusCode(), resp.String())
	}

	return nil
}
