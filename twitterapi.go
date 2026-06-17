package twitterapi

import (
	"context"
	"net/http"
	"sync"
	"time"
)

var (
	domainURI              = "https://api.twitterapi.io"
	oapiDomainURI          = domainURI + "/oapi"
	twitterDomainURI       = domainURI + "/twitter"
	userTwitterDomainURI   = twitterDomainURI + "/user"
	tweetsTwitterDomainURI = twitterDomainURI + "/tweets"
	listTwitterDomainURI   = twitterDomainURI + "/list"
	streamDomainURI        = oapiDomainURI + "/x_user_stream"
	tweetFilterDomainURI   = oapiDomainURI + "/tweet_filter"
)

type TwitterApi struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	httpClient *http.Client
	headers    map[string]string

	requestTimeout time.Duration
	retryConfig    RetryConfig

	proxy   string
	cookies Cookies
}

func New(xApiKey string, opts ...Option) *TwitterApi {

	ctx, cancel := context.WithCancel(context.Background())

	t := &TwitterApi{
		ctx:    ctx,
		cancel: cancel,

		httpClient: newHTTP(),
		headers:    map[string]string{},

		requestTimeout: 10 * time.Second,
		retryConfig:    DefaultRetryConfig,
	}

	for _, opt := range opts {
		opt(t)
	}

	// t.headers["accept-encoding"] = "gzip"
	t.headers["X-API-Key"] = xApiKey

	return t
}

func (t *TwitterApi) Close() {
	t.cancel()
	t.wg.Wait()
}
