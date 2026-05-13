package twitterapi

import (
	neturl "net/url"
	"os"
	"strings"
	"testing"
	"time"
)

const (
	batchGetUserInfoMaxAttempts = 3
	batchGetUserInfoRetryDelay  = 200 * time.Millisecond
)

func buildBatchGetUserInfoURL(userIds []string) string {
	vals := neturl.Values{}
	vals.Set("userIds", strings.Join(userIds, ","))
	return userTwitterDomainURI + "/batch_info_by_ids?" + vals.Encode()
}

func doGetWithRetry5xx(t *testing.T, client *TwitterApi, url string) ([]byte, int, error) {
	t.Helper()

	var raw []byte
	var statusCode int
	var err error

	for attempt := 1; attempt <= batchGetUserInfoMaxAttempts; attempt++ {
		raw, statusCode, err = doGet(t, client, url)
		if err != nil {
			return raw, statusCode, err
		}
		if statusCode/100 != 5 {
			return raw, statusCode, nil
		}
		if attempt < batchGetUserInfoMaxAttempts {
			t.Logf("BatchGetUserInfoByUserIds attempt %d returned %d, retrying after %s", attempt, statusCode, batchGetUserInfoRetryDelay)
			time.Sleep(batchGetUserInfoRetryDelay)
		}
	}

	return raw, statusCode, nil
}

func batchGetUserInfoByUserIdsWithRetry5xx(t *testing.T, client *TwitterApi, userIds []string) (*BatchGetUserInfoByUserIdsResponse, error) {
	t.Helper()

	if len(userIds) == 0 {
		return client.BatchGetUserInfoByUserIds(userIds)
	}

	url := buildBatchGetUserInfoURL(userIds)
	var lastErr error

	for attempt := 1; attempt <= batchGetUserInfoMaxAttempts; attempt++ {
		resp, err := client.BatchGetUserInfoByUserIds(userIds)
		if err == nil {
			return resp, nil
		}
		lastErr = err

		raw, statusCode, probeErr := doGet(t, client, url)
		if probeErr != nil || statusCode/100 != 5 {
			return nil, lastErr
		}
		if attempt < batchGetUserInfoMaxAttempts {
			t.Logf("BatchGetUserInfoByUserIds behavior attempt %d returned %d, retrying after %s; body: %s", attempt, statusCode, batchGetUserInfoRetryDelay, string(raw))
			time.Sleep(batchGetUserInfoRetryDelay)
		}
	}

	return nil, lastErr
}

func TestBatchGetUserInfoByUserIds(t *testing.T) {
	client := newTestClient(t)
	// single id
	vals := neturl.Values{}
	vals.Set("userIds", testUserID)
	url := userTwitterDomainURI + "/batch_info_by_ids?" + vals.Encode()

	raw, statusCode, err := doGetWithRetry5xx(t, client, url)
	if err != nil {
		t.Fatalf("BatchGetUserInfoByUserIds request failed: %v", err)
	}
	t.Logf("RAW JSON: %s", string(raw))
	if statusCode/100 != 2 {
		t.Fatalf("BatchGetUserInfoByUserIds returned non-2xx status: %d, body: %s", statusCode, string(raw))
	}

	var response BatchGetUserInfoByUserIdsResponse
	if decodeErr := decodeJSONDisallowUnknowns(t, raw, &response); decodeErr != nil {
		t.Logf("DisallowUnknownFields decode error: %v", decodeErr)
	}
	logJSONFieldDiff(t, raw, &response)
}

func TestBatchGetUserInfoByUserIds_Behavior(t *testing.T) {
	apiKey := os.Getenv("TWITTERAPI_IO_KEY")
	if apiKey == "" {
		t.Skip("TWITTERAPI_IO_KEY not set")
	}

	t.Run("empty userIds returns error", func(t *testing.T) {
		x := New(apiKey)
		_, err := x.BatchGetUserInfoByUserIds([]string{})
		if err == nil {
			t.Fatal("expected error for empty userIds, got nil")
		}
	})

	t.Run("valid userIds returns data", func(t *testing.T) {
		x := New(apiKey)
		resp, err := batchGetUserInfoByUserIdsWithRetry5xx(t, x, []string{testUserID})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected non-nil response")
		}
		if len(resp.Users) == 0 {
			t.Fatal("expected at least one user in response")
		}
		for _, u := range resp.Users {
			if u == nil {
				t.Fatal("nil user in response list")
			}
			if strings.TrimSpace(u.ID) == "" {
				t.Fatal("expected non-empty user ID")
			}
		}
	})
}
