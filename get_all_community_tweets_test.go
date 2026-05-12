package twitterapi

import (
	neturl "net/url"
	"os"
	"testing"
)

func TestGetAllCommunityTweets(t *testing.T) {
	client := newTestClient(t)
	query := "from:" + testUserName
	vals := neturl.Values{}
	vals.Set("query", query)
	url := twitterDomainURI + "/community/get_tweets_from_all_community?" + vals.Encode()

	raw, statusCode, err := doGet(t, client, url)
	if err != nil {
		t.Fatalf("GetAllCommunityTweets request failed: %v", err)
	}
	t.Logf("RAW JSON: %s", string(raw))
	if statusCode/100 != 2 {
		t.Fatalf("GetAllCommunityTweets returned non-2xx status: %d", statusCode)
	}

	var response GetAllCommunityTweetsResponse
	if decodeErr := decodeJSONDisallowUnknowns(t, raw, &response); decodeErr != nil {
		t.Logf("DisallowUnknownFields decode error: %v", decodeErr)
	}
	logJSONFieldDiff(t, raw, &response)
}

func TestGetAllCommunityTweets_Behavior(t *testing.T) {
	apiKey := os.Getenv("TWITTERAPI_IO_KEY")
	if apiKey == "" {
		t.Skip("TWITTERAPI_IO_KEY not set")
	}

	t.Run("empty query returns error", func(t *testing.T) {
		x := New(apiKey)
		_, err := x.GetAllCommunityTweets("", nil)
		if err == nil {
			t.Fatal("expected error for empty query, got nil")
		}
	})

	t.Run("valid query returns data", func(t *testing.T) {
		x := New(apiKey)
		resp, err := x.GetAllCommunityTweets("from:"+testUserName, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected non-nil response")
		}
	})
}
