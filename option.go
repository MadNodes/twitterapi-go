package twitterapi

import "net/http"

type Option func(*twitterApi)

func WithHttpClient(cli *http.Client) Option {
	return func(t *twitterApi) {
		t.httpClient = cli
	}
}

func WithHeader(headers map[string]string) Option {
	return func(t *twitterApi) {
		t.headers = headers
	}
}

func WithProxy(proxy string) Option {
	return func(t *twitterApi) {
		t.proxy = proxy
	}
}

func WithCookies(cookies Cookies) Option {
	return func(t *twitterApi) {
		t.cookies = cookies
	}
}
