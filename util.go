package twitterapi

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"time"

	"github.com/klauspost/compress/gzhttp"
)

var (
	defaultMaxIdleConnsPerHost = 9
	defaultTimeout             = 5 * time.Minute
	defaultKeepAlive           = 180 * time.Second
)

type RetryConfig struct {
	MaxAttempts  int
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Multiplier   float64
}

var DefaultRetryConfig = RetryConfig{
	MaxAttempts:  3,
	InitialDelay: 500 * time.Millisecond,
	MaxDelay:     5 * time.Second,
	Multiplier:   2.0,
}

func withRetry(
	ctx context.Context,
	config RetryConfig,
	isRetryable func(error) bool,
	fn func() error,
) error {
	config = normalizeRetryConfig(config)
	for attempt := 1; attempt <= config.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		err := fn()
		if err == nil {
			return nil
		}
		if attempt == config.MaxAttempts || !isRetryable(err) {
			return err
		}
		if err = sleepBeforeRetry(ctx, config, attempt); err != nil {
			return err
		}
	}

	return ctx.Err()
}

func normalizeRetryConfig(config RetryConfig) RetryConfig {
	if config.MaxAttempts <= 0 {
		config.MaxAttempts = 1
	}
	if config.InitialDelay < 0 {
		config.InitialDelay = 0
	}
	if config.MaxDelay < 0 {
		config.MaxDelay = 0
	}
	if config.MaxDelay > 0 && config.InitialDelay > config.MaxDelay {
		config.InitialDelay = config.MaxDelay
	}
	if config.Multiplier < 1 {
		config.Multiplier = 1
	}
	return config
}

func retryDelay(config RetryConfig, failedAttempt int) time.Duration {
	delay := config.InitialDelay
	for i := 1; i < failedAttempt; i++ {
		delay = time.Duration(float64(delay) * config.Multiplier)
		if config.MaxDelay > 0 && delay > config.MaxDelay {
			delay = config.MaxDelay
			break
		}
	}
	if config.MaxDelay > 0 && delay > config.MaxDelay {
		delay = config.MaxDelay
	}
	return jitterDelay(delay)
}

func jitterDelay(delay time.Duration) time.Duration {
	if delay <= 0 {
		return 0
	}
	spread := int64(delay / 5)
	if spread <= 0 {
		return delay
	}
	offset := rand.Int63n(spread*2+1) - spread
	return delay + time.Duration(offset)
}

func sleepBeforeRetry(ctx context.Context, config RetryConfig, failedAttempt int) error {
	delay := retryDelay(config, failedAttempt)
	if delay <= 0 {
		return nil
	}

	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func newHTTPTransport() *http.Transport {
	return &http.Transport{
		IdleConnTimeout:     defaultTimeout,
		MaxConnsPerHost:     defaultMaxIdleConnsPerHost,
		MaxIdleConnsPerHost: defaultMaxIdleConnsPerHost,
		Proxy:               http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   defaultTimeout,
			KeepAlive: defaultKeepAlive,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2: true,
		// MaxIdleConns:          100,
		TLSHandshakeTimeout: 10 * time.Second,
		// ExpectContinueTimeout: 1 * time.Second,
	}
}

func newHTTP() *http.Client {
	tr := newHTTPTransport()

	return &http.Client{
		Timeout:   defaultTimeout,
		Transport: gzhttp.Transport(tr),
	}
}

type getResponse struct {
	body []byte
	resp *http.Response
}

type retryableStatusError struct {
	body []byte
	resp *http.Response
}

func (e *retryableStatusError) Error() string {
	if e == nil || e.resp == nil {
		return "retryable http status"
	}
	return fmt.Sprintf("retryable http status %d", e.resp.StatusCode)
}

func (t *TwitterApi) getDataWithHeader(ctx context.Context, url string, headers map[string]string) ([]byte, *http.Response, error) {
	var result getResponse
	err := withRetry(ctx, t.retryConfig, isRetryableGetErr, func() error {
		attemptCtx, cancel := context.WithTimeout(ctx, t.requestTimeout)
		defer cancel()

		body, resp, err := getDataWithHeaderOnce(attemptCtx, t.httpClient, url, headers)
		if err != nil {
			return normalizeTimeoutErr(err)
		}
		result = getResponse{body: body, resp: resp}
		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= http.StatusInternalServerError {
			return &retryableStatusError{body: body, resp: resp}
		}
		return nil
	})
	if err != nil {
		var statusErr *retryableStatusError
		if errors.As(err, &statusErr) {
			return statusErr.body, statusErr.resp, nil
		}
		return nil, nil, err
	}

	return result.body, result.resp, nil
}

func getDataWithHeaderOnce(ctx context.Context, client *http.Client, url string, headers map[string]string) ([]byte, *http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	buff, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	return buff, resp, nil
}

func postDataWithHeader(ctx context.Context, client *http.Client, url string, ioParams io.Reader, headers map[string]string) ([]byte, *http.Response, error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, ioParams)
	if err != nil {
		return nil, nil, err
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	buff, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	return buff, resp, nil
}

func patchDataWithHeader(ctx context.Context, client *http.Client, url string, ioParams io.Reader, headers map[string]string) ([]byte, *http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, ioParams)
	if err != nil {
		return nil, nil, err
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	buff, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	return buff, resp, nil
}

func deleteDataWithHeader(ctx context.Context, client *http.Client, url string, ioParams io.Reader, headers map[string]string) ([]byte, *http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, ioParams)
	if err != nil {
		return nil, nil, err
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	buff, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	return buff, resp, nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func isTimeoutErr(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	var netErr net.Error
	return errors.As(err, &netErr) && netErr.Timeout()
}

func normalizeTimeoutErr(err error) error {
	if !isTimeoutErr(err) || errors.Is(err, context.DeadlineExceeded) {
		return err
	}
	return fmt.Errorf("%w: %w", context.DeadlineExceeded, err)
}

func isRetryableGetErr(err error) bool {
	var statusErr *retryableStatusError
	if errors.As(err, &statusErr) {
		return true
	}
	if isTimeoutErr(err) {
		return true
	}
	var netErr net.Error
	return errors.As(err, &netErr) && netErr.Temporary()
}
