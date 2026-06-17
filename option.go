package twitterapi

import (
	"net/http"
	"time"
)

type Option func(*TwitterApi)

func WithHttpClient(cli *http.Client) Option {
	return func(t *TwitterApi) {
		t.httpClient = cli
	}
}

func WithHeader(headers map[string]string) Option {
	return func(t *TwitterApi) {
		t.headers = headers
	}
}

func WithProxy(proxy string) Option {
	return func(t *TwitterApi) {
		t.proxy = proxy
	}
}

func WithCookies(cookies Cookies) Option {
	return func(t *TwitterApi) {
		t.cookies = cookies
	}
}

func WithRequestTimeout(timeout time.Duration) Option {
	return func(t *TwitterApi) {
		if timeout > 0 {
			t.requestTimeout = timeout
		}
	}
}

func WithRetryConfig(config RetryConfig) Option {
	return func(t *TwitterApi) {
		t.retryConfig = normalizeRetryConfig(config)
	}
}
