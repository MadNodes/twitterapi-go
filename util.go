package twitterapi

import (
	"context"
	"errors"
	"fmt"
	"io"
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

func getDataWithHeader(ctx context.Context, client *http.Client, url string, headers map[string]string) ([]byte, *http.Response, error) {
	var lastErr error
	for attempt := 0; attempt < 2; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, nil, err
		}
		for k, v := range headers {
			req.Header.Add(k, v)
		}
		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			if attempt == 0 && isTimeoutErr(err) {
				continue
			}
			return nil, nil, normalizeTimeoutErr(err)
		}
		defer resp.Body.Close()
		buff, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, nil, err
		}
		return buff, resp, nil
	}
	return nil, nil, normalizeTimeoutErr(lastErr)
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
	return fmt.Errorf("%w: %v", context.DeadlineExceeded, err)
}
