package twitterapi

import (
	neturl "net/url"
	"os"
	"testing"
)

func TestGetTweetByIDs(t *testing.T) {
	client := newTestClient(t)
	vals := neturl.Values{}
	vals.Set("tweet_ids", testTweetID)
	url := tweetsTwitterDomainURI + "?" + vals.Encode()

	raw, statusCode, err := doGet(t, client, url)
	if err != nil {
		t.Fatalf("GetTweetByIDs request failed: %v", err)
	}
	t.Logf("RAW JSON: %s", string(raw))
	if statusCode/100 != 2 {
		t.Fatalf("GetTweetByIDs returned non-2xx status: %d", statusCode)
	}

	var response GetTweetByIDsResponse
	if decodeErr := decodeJSONDisallowUnknowns(t, raw, &response); decodeErr != nil {
		t.Logf("DisallowUnknownFields decode error: %v", decodeErr)
	}
	logJSONFieldDiff(t, raw, &response)
}

func TestGetTweetByIDs_Behavior(t *testing.T) {
	apiKey := os.Getenv("TWITTERAPI_IO_KEY")
	if apiKey == "" {
		t.Skip("TWITTERAPI_IO_KEY not set")
	}

	t.Run("empty slice returns error", func(t *testing.T) {
		x := New(apiKey)
		_, err := x.GetTweetByIDs([]string{})
		if err == nil {
			t.Fatal("expected error for empty tweetIDs, got nil")
		}
	})

	t.Run("valid tweetIDs returns data", func(t *testing.T) {
		x := New(apiKey)
		resp, err := x.GetTweetByIDs([]string{testTweetID})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected non-nil response")
		}
	})
}
