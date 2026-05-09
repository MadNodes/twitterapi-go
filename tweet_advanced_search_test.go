package twitterapi

import (
	"os"
	"testing"
)

func TestTweetAdvancedSearch(t *testing.T) {
	client := newTestClient(t)
	query := "from:" + testUserName
	url := twitterDomainURI + "/tweet/advanced_search?query=" + query

	raw, statusCode, err := doGet(t, client, url)
	if err != nil {
		t.Fatalf("TweetAdvancedSearch request failed: %v", err)
	}
	t.Logf("RAW JSON: %s", string(raw))
	if statusCode/100 != 2 {
		t.Fatalf("TweetAdvancedSearch returned non-2xx status: %d", statusCode)
	}

	var response TweetAdvancedSearchResponse
	if decodeErr := decodeJSONDisallowUnknowns(t, raw, &response); decodeErr != nil {
		t.Logf("DisallowUnknownFields decode error: %v", decodeErr)
	}
	logJSONFieldDiff(t, raw, &response)
}

func TestTweetAdvancedSearch_Behavior(t *testing.T) {
	apiKey := os.Getenv("TWITTERAPI_IO_KEY")
	if apiKey == "" {
		t.Skip("TWITTERAPI_IO_KEY not set")
	}

	t.Run("empty query returns error", func(t *testing.T) {
		x := New(apiKey)
		_, err := x.TweetAdvancedSearch("", nil)
		if err == nil {
			t.Fatal("expected error for empty query, got nil")
		}
	})

	t.Run("valid query returns data", func(t *testing.T) {
		x := New(apiKey)
		resp, err := x.TweetAdvancedSearch("from:"+testUserName, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected non-nil response")
		}
	})
}
