package twitterapi

import (
	neturl "net/url"
	"os"
	"strconv"
	"testing"
)

func TestGetTweetQuote(t *testing.T) {
	client := newTestClient(t)
	sinceTime, untilTime := last7DaysRangeUnixInt()
	vals := neturl.Values{}
	vals.Set("tweetId", testTweetID)
	vals.Set("sinceTime", strconv.Itoa(sinceTime))
	vals.Set("untilTime", strconv.Itoa(untilTime))
	url := twitterDomainURI + "/tweet/quotes?" + vals.Encode()

	raw, statusCode, err := doGet(t, client, url)
	if err != nil {
		t.Fatalf("GetTweetQuote request failed: %v", err)
	}
	t.Logf("RAW JSON: %s", string(raw))
	if statusCode/100 != 2 {
		t.Fatalf("GetTweetQuote returned non-2xx status: %d", statusCode)
	}

	var response GetTweetQuoteResponse
	if decodeErr := decodeJSONDisallowUnknowns(t, raw, &response); decodeErr != nil {
		t.Logf("DisallowUnknownFields decode error: %v", decodeErr)
	}
	logJSONFieldDiff(t, raw, &response)
}

func TestGetTweetQuote_Behavior(t *testing.T) {
	apiKey := os.Getenv("TWITTERAPI_IO_KEY")
	if apiKey == "" {
		t.Skip("TWITTERAPI_IO_KEY not set")
	}

	t.Run("empty tweetID returns error", func(t *testing.T) {
		x := New(apiKey)
		_, err := x.GetTweetQuote("", nil, nil, nil, nil)
		if err == nil {
			t.Fatal("expected error for empty tweetID, got nil")
		}
	})

	t.Run("valid tweetID returns data", func(t *testing.T) {
		x := New(apiKey)
		resp, err := x.GetTweetQuote(testTweetID, nil, nil, nil, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected non-nil response")
		}
	})
}
