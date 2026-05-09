package twitterapi

import (
	"os"
	"testing"
)

func TestGetUserLastTweets(t *testing.T) {
	client := newTestClient(t)
	url := userTwitterDomainURI + "/last_tweets?userId=" + testUserID

	raw, statusCode, err := doGet(t, client, url)
	if err != nil {
		t.Fatalf("GetUserLastTweets request failed: %v", err)
	}
	t.Logf("RAW JSON: %s", string(raw))
	if statusCode/100 != 2 {
		t.Fatalf("GetUserLastTweets returned non-2xx status: %d", statusCode)
	}

	var response GetUserLastTweetsResponse
	if decodeErr := decodeJSONDisallowUnknowns(t, raw, &response); decodeErr != nil {
		t.Logf("DisallowUnknownFields decode error: %v", decodeErr)
	}
	logJSONFieldDiff(t, raw, &response)
}

func TestGetUserLastTweets_Behavior(t *testing.T) {
	apiKey := os.Getenv("TWITTERAPI_IO_KEY")
	if apiKey == "" {
		t.Skip("TWITTERAPI_IO_KEY not set")
	}

	t.Run("empty userId and userName returns error", func(t *testing.T) {
		x := New(apiKey)
		_, err := x.GetUserLastTweets(nil, nil, nil, nil)
		if err == nil {
			t.Fatal("expected error for empty params, got nil")
		}
	})

	t.Run("valid userId returns data", func(t *testing.T) {
		x := New(apiKey)
		uid := testUserID
		resp, err := x.GetUserLastTweets(&uid, nil, nil, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected non-nil response")
		}
		if resp.Data == nil || len(resp.Data.Tweets) == 0 {
			t.Fatal("expected at least one tweet")
		}
	})
}
