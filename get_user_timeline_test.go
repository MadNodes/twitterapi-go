package twitterapi

import (
	neturl "net/url"
	"os"
	"testing"
)

func TestGetUserTimeline(t *testing.T) {
	client := newTestClient(t)
	vals := neturl.Values{}
	vals.Set("userId", testUserID)
	url := userTwitterDomainURI + "/tweet_timeline?" + vals.Encode()

	raw, statusCode, err := doGet(t, client, url)
	if err != nil {
		t.Fatalf("GetUserTimeline request failed: %v", err)
	}
	t.Logf("RAW JSON: %s", string(raw))
	if statusCode/100 != 2 {
		t.Fatalf("GetUserTimeline returned non-2xx status: %d", statusCode)
	}

	var response GetUserTimelineResponse
	if decodeErr := decodeJSONDisallowUnknowns(t, raw, &response); decodeErr != nil {
		t.Logf("DisallowUnknownFields decode error: %v", decodeErr)
	}
	logJSONFieldDiff(t, raw, &response)
}

func TestGetUserTimeline_Behavior(t *testing.T) {
	apiKey := os.Getenv("TWITTERAPI_IO_KEY")
	if apiKey == "" {
		t.Skip("TWITTERAPI_IO_KEY not set")
	}

	t.Run("empty userId returns error", func(t *testing.T) {
		x := New(apiKey)
		_, err := x.GetUserTimeline("", nil, nil, nil)
		if err == nil {
			t.Fatal("expected error for empty userId, got nil")
		}
	})

	t.Run("valid userId returns data", func(t *testing.T) {
		x := New(apiKey)
		resp, err := x.GetUserTimeline(testUserID, nil, nil, nil)
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
