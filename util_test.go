package twitterapi

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type timeoutErr struct{}

func (timeoutErr) Error() string   { return "timeout" }
func (timeoutErr) Timeout() bool   { return true }
func (timeoutErr) Temporary() bool { return true }

func TestGetDataWithHeader_RetriesTimeoutOnce(t *testing.T) {
	attempts := 0
	api := New("test",
		WithHttpClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			attempts++
			if attempts == 1 {
				return nil, timeoutErr{}
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader("ok")),
				Header:     make(http.Header),
			}, nil
		})}),
		WithRequestTimeout(time.Second),
		WithRetryConfig(RetryConfig{MaxAttempts: 2}),
	)

	body, resp, err := api.getDataWithHeader(context.Background(), "http://example.com", map[string]string{"X-Test": "1"})
	if err != nil {
		t.Fatalf("TwitterApi.getDataWithHeader returned error: %v", err)
	}
	if attempts != 2 {
		t.Fatalf("expected 2 attempts, got %d", attempts)
	}
	if resp == nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 response, got %#v", resp)
	}
	if string(body) != "ok" {
		t.Fatalf("expected body ok, got %q", string(body))
	}
}

func TestGetDataWithHeader_ReturnsDeadlineExceededAfterRetry(t *testing.T) {
	attempts := 0
	api := New("test",
		WithHttpClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			attempts++
			return nil, timeoutErr{}
		})}),
		WithRequestTimeout(time.Second),
		WithRetryConfig(RetryConfig{MaxAttempts: 2}),
	)

	_, _, err := api.getDataWithHeader(context.Background(), "http://example.com", nil)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected deadline exceeded, got %v", err)
	}
	if attempts != 2 {
		t.Fatalf("expected 2 attempts, got %d", attempts)
	}
}

func TestGetDataWithHeader_DoesNotRetryUnauthorized(t *testing.T) {
	attempts := 0
	api := New("test",
		WithHttpClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			attempts++
			return &http.Response{
				StatusCode: http.StatusUnauthorized,
				Body:       io.NopCloser(strings.NewReader("unauthorized")),
				Header:     make(http.Header),
			}, nil
		})}),
		WithRequestTimeout(time.Second),
		WithRetryConfig(RetryConfig{MaxAttempts: 3}),
	)

	_, resp, err := api.getDataWithHeader(context.Background(), "http://example.com", nil)
	if err != nil {
		t.Fatalf("expected status response without helper error, got %v", err)
	}
	if resp == nil || resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 response, got %#v", resp)
	}
	if attempts != 1 {
		t.Fatalf("expected 1 attempt, got %d", attempts)
	}
}

func TestGetDataWithHeader_RetriesServerError(t *testing.T) {
	attempts := 0
	api := New("test",
		WithHttpClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			attempts++
			if attempts == 1 {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(strings.NewReader("server error")),
					Header:     make(http.Header),
				}, nil
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader("ok")),
				Header:     make(http.Header),
			}, nil
		})}),
		WithRequestTimeout(time.Second),
		WithRetryConfig(RetryConfig{MaxAttempts: 2}),
	)

	body, _, err := api.getDataWithHeader(context.Background(), "http://example.com", nil)
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if attempts != 2 {
		t.Fatalf("expected 2 attempts, got %d", attempts)
	}
	if string(body) != "ok" {
		t.Fatalf("expected body ok, got %q", string(body))
	}
}
