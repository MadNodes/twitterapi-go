package twitterapi

import (
	"net/http"
	"net/http/httptest"
	neturl "net/url"
	"os"
	"testing"
)

func TestGetTweetRepliesV2(t *testing.T) {
	client := newTestClient(t)
	vals := neturl.Values{}
	vals.Set("tweetId", testTweetID)
	url := twitterDomainURI + "/tweet/replies/v2?" + vals.Encode()

	raw, statusCode, err := doGet(t, client, url)
	if err != nil {
		t.Fatalf("GetTweetRepliesV2 request failed: %v", err)
	}
	t.Logf("RAW JSON: %s", string(raw))
	if statusCode/100 != 2 {
		t.Fatalf("GetTweetRepliesV2 returned non-2xx status: %d", statusCode)
	}

	var response GetTweetRepliesV2Response
	if decodeErr := decodeJSONDisallowUnknowns(t, raw, &response); decodeErr != nil {
		t.Logf("DisallowUnknownFields decode error: %v", decodeErr)
	}
	logJSONFieldDiff(t, raw, &response)
}

func TestGetTweetRepliesV2_Behavior(t *testing.T) {
	apiKey := os.Getenv("TWITTERAPI_IO_KEY")
	if apiKey == "" {
		t.Skip("TWITTERAPI_IO_KEY not set")
	}

	t.Run("empty tweetID returns error", func(t *testing.T) {
		x := New(apiKey)
		_, err := x.GetTweetRepliesV2("", nil)
		if err == nil {
			t.Fatal("expected error for empty tweetID, got nil")
		}
	})

	t.Run("valid tweetID returns data", func(t *testing.T) {
		x := New(apiKey)
		resp, err := x.GetTweetRepliesV2(testTweetID, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected non-nil response")
		}
	})
}

func TestGetTweetRepliesV2WithQueryType_BuildsQuery(t *testing.T) {
	var gotQuery neturl.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"tweets":[],"has_next_page":false,"next_cursor":"","status":"success","msg":""}`))
	}))
	t.Cleanup(srv.Close)

	origURL := twitterDomainURI
	twitterDomainURI = srv.URL
	t.Cleanup(func() { twitterDomainURI = origURL })

	cursor := "cursor-1"
	client := New("test-key", WithHttpClient(srv.Client()))
	_, err := client.GetTweetRepliesV2WithQueryType("tweet-1", GetTweetRepliesV2QueryTypeLatest, &cursor)
	if err != nil {
		t.Fatalf("GetTweetRepliesV2WithQueryType returned error: %v", err)
	}

	if gotQuery.Get("tweetId") != "tweet-1" {
		t.Fatalf("expected tweetId tweet-1, got %q", gotQuery.Get("tweetId"))
	}
	if gotQuery.Get("queryType") != string(GetTweetRepliesV2QueryTypeLatest) {
		t.Fatalf("expected queryType Latest, got %q", gotQuery.Get("queryType"))
	}
	if gotQuery.Get("cursor") != cursor {
		t.Fatalf("expected cursor %q, got %q", cursor, gotQuery.Get("cursor"))
	}
}

func TestGetTweetRepliesV2_OmitsDefaultQueryType(t *testing.T) {
	var gotQuery neturl.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"tweets":[],"has_next_page":false,"next_cursor":"","status":"success","msg":""}`))
	}))
	t.Cleanup(srv.Close)

	origURL := twitterDomainURI
	twitterDomainURI = srv.URL
	t.Cleanup(func() { twitterDomainURI = origURL })

	client := New("test-key", WithHttpClient(srv.Client()))
	_, err := client.GetTweetRepliesV2("tweet-1", nil)
	if err != nil {
		t.Fatalf("GetTweetRepliesV2 returned error: %v", err)
	}

	if gotQuery.Get("queryType") != "" {
		t.Fatalf("expected queryType to be omitted, got %q", gotQuery.Get("queryType"))
	}
}
