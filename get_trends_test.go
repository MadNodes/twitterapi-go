package twitterapi

import (
	neturl "net/url"
	"os"
	"strconv"
	"testing"
)

func TestGetTrends(t *testing.T) {
	client := newTestClient(t)
	count := 10
	vals := neturl.Values{}
	vals.Set("woeid", "1")
	vals.Set("count", strconv.Itoa(count))
	url := twitterDomainURI + "/trends?" + vals.Encode()

	raw, statusCode, err := doGet(t, client, url)
	if err != nil {
		t.Fatalf("GetTrends request failed: %v", err)
	}
	t.Logf("RAW JSON: %s", string(raw))
	if statusCode/100 != 2 {
		t.Fatalf("GetTrends returned non-2xx status: %d", statusCode)
	}

	var response GetTrendsResponse
	if decodeErr := decodeJSONDisallowUnknowns(t, raw, &response); decodeErr != nil {
		t.Logf("DisallowUnknownFields decode error: %v", decodeErr)
	}
	logJSONFieldDiff(t, raw, &response)
}

func TestGetTrends_Behavior(t *testing.T) {
	apiKey := os.Getenv("TWITTERAPI_IO_KEY")
	if apiKey == "" {
		t.Skip("TWITTERAPI_IO_KEY not set")
	}

	t.Run("invalid woeid returns error", func(t *testing.T) {
		x := New(apiKey)
		_, err := x.GetTrends(0, nil, nil)
		if err == nil {
			t.Fatal("expected error for invalid woeid, got nil")
		}
	})

	t.Run("valid woeid returns data", func(t *testing.T) {
		x := New(apiKey)
		resp, err := x.GetTrends(1, nil, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected non-nil response")
		}
		if len(resp.Trends) == 0 {
			t.Fatal("expected at least one trend")
		}
	})
}
