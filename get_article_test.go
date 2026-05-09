package twitterapi

import (
	"os"
	"testing"
)

func TestGetArticle(t *testing.T) {
	client := newTestClient(t)
	url := twitterDomainURI + "/article?tweet_id=" + testTweetID

	raw, statusCode, err := doGet(t, client, url)
	if err != nil {
		t.Fatalf("GetArticle request failed: %v", err)
	}
	t.Logf("RAW JSON: %s", string(raw))
	if statusCode/100 != 2 {
		t.Fatalf("GetArticle returned non-2xx status: %d", statusCode)
	}

	var response GetArticleResponse
	if decodeErr := decodeJSONDisallowUnknowns(t, raw, &response); decodeErr != nil {
		t.Logf("DisallowUnknownFields decode error: %v", decodeErr)
	}
	logJSONFieldDiff(t, raw, &response)
}

func TestGetArticle_Behavior(t *testing.T) {
	apiKey := os.Getenv("TWITTERAPI_IO_KEY")
	if apiKey == "" {
		t.Skip("TWITTERAPI_IO_KEY not set")
	}

	t.Run("empty tweetID returns error", func(t *testing.T) {
		x := New(apiKey)
		_, err := x.GetArticle("")
		if err == nil {
			t.Fatal("expected error for empty tweetID, got nil")
		}
	})

	t.Run("valid tweetID returns data", func(t *testing.T) {
		x := New(apiKey)
		resp, err := x.GetArticle("1905545699552375179")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected non-nil response")
		}
		if resp.Article == nil {
			t.Fatal("expected non-nil article")
		}
	})
}
