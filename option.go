package twitterapi

import "net/http"

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
