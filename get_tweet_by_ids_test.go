package twitterapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	neturl "net/url"
	"os"
	"strings"
	"sync"
	"testing"

	"net/http/httptest"
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

func TestGetTweetByIDs_Batching(t *testing.T) {
	// newMock creates a TwitterApi client backed by a test HTTP server.
	newMock := func(t *testing.T, handler func(http.ResponseWriter, *http.Request)) *TwitterApi {
		t.Helper()
		srv := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(srv.Close)

		origURL := tweetsTwitterDomainURI
		tweetsTwitterDomainURI = srv.URL
		t.Cleanup(func() { tweetsTwitterDomainURI = origURL })

		return New("test-key", WithHttpClient(srv.Client()))
	}

	// generateIDs creates n sequential ID strings like "id_1", "id_2", ...
	generateIDs := func(n int) []string {
		ids := make([]string, n)
		for i := range ids {
			ids[i] = fmt.Sprintf("id_%d", i+1)
		}
		return ids
	}

	// buildTweetResponse creates a success response with one tweet per input ID.
	buildTweetResponse := func(ids []string) *GetTweetByIDsResponse {
		tweets := make([]*GetTweetByIDsTweet, len(ids))
		for i, id := range ids {
			tweets[i] = &GetTweetByIDsTweet{ID: id}
		}
		return &GetTweetByIDsResponse{
			Tweets: tweets,
			Status: "success",
			Code:   200,
		}
	}

	t.Run("10 ids single batch", func(t *testing.T) {
		api := newMock(t, func(w http.ResponseWriter, r *http.Request) {
			ids := strings.Split(r.URL.Query().Get("tweet_ids"), ",")
			json.NewEncoder(w).Encode(buildTweetResponse(ids))
		})

		input := generateIDs(10)
		resp, err := api.GetTweetByIDs(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(resp.Tweets) != 10 {
			t.Fatalf("expected 10 tweets, got %d", len(resp.Tweets))
		}
		if resp.Status != "success" {
			t.Errorf("expected status 'success', got %q", resp.Status)
		}
		for i, tw := range resp.Tweets {
			if tw.ID != input[i] {
				t.Errorf("tweet %d: expected ID %s, got %s", i, input[i], tw.ID)
			}
		}
	})

	t.Run("50 ids boundary single batch", func(t *testing.T) {
		api := newMock(t, func(w http.ResponseWriter, r *http.Request) {
			ids := strings.Split(r.URL.Query().Get("tweet_ids"), ",")
			if len(ids) != 50 {
				t.Errorf("expected 50 ids in single batch request, got %d", len(ids))
			}
			json.NewEncoder(w).Encode(buildTweetResponse(ids))
		})

		input := generateIDs(50)
		resp, err := api.GetTweetByIDs(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(resp.Tweets) != 50 {
			t.Fatalf("expected 50 tweets, got %d", len(resp.Tweets))
		}
		for i, tw := range resp.Tweets {
			if tw.ID != input[i] {
				t.Errorf("tweet %d: expected ID %s, got %s", i, input[i], tw.ID)
			}
		}
	})

	t.Run("51 ids split into 2 batches", func(t *testing.T) {
		var (
			mu    sync.Mutex
			sizes []int
		)
		api := newMock(t, func(w http.ResponseWriter, r *http.Request) {
			ids := strings.Split(r.URL.Query().Get("tweet_ids"), ",")
			mu.Lock()
			sizes = append(sizes, len(ids))
			mu.Unlock()

			if len(ids) > maxTweetIDsPerRequest {
				t.Errorf("batch size %d exceeds max %d", len(ids), maxTweetIDsPerRequest)
			}
			json.NewEncoder(w).Encode(buildTweetResponse(ids))
		})

		input := generateIDs(51)
		resp, err := api.GetTweetByIDs(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		mu.Lock()
		if len(sizes) != 2 {
			t.Fatalf("expected 2 requests, got %d: %v", len(sizes), sizes)
		}
		if sizes[0] != 50 {
			t.Errorf("first batch: expected 50 ids, got %d", sizes[0])
		}
		if sizes[1] != 1 {
			t.Errorf("second batch: expected 1 id, got %d", sizes[1])
		}
		mu.Unlock()

		if len(resp.Tweets) != 51 {
			t.Fatalf("expected 51 tweets, got %d", len(resp.Tweets))
		}
		for i, tw := range resp.Tweets {
			if tw.ID != input[i] {
				t.Errorf("tweet %d: expected ID %s, got %s", i, input[i], tw.ID)
			}
		}
	})

	t.Run("120 ids split into 3 batches", func(t *testing.T) {
		var (
			mu    sync.Mutex
			sizes []int
		)
		api := newMock(t, func(w http.ResponseWriter, r *http.Request) {
			ids := strings.Split(r.URL.Query().Get("tweet_ids"), ",")
			mu.Lock()
			sizes = append(sizes, len(ids))
			mu.Unlock()

			if len(ids) > maxTweetIDsPerRequest {
				t.Errorf("batch size %d exceeds max %d", len(ids), maxTweetIDsPerRequest)
			}
			json.NewEncoder(w).Encode(buildTweetResponse(ids))
		})

		input := generateIDs(120)
		resp, err := api.GetTweetByIDs(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		mu.Lock()
		if len(sizes) != 3 {
			t.Fatalf("expected 3 requests, got %d: %v", len(sizes), sizes)
		}
		if sizes[0] != 50 {
			t.Errorf("first batch: expected 50 ids, got %d", sizes[0])
		}
		if sizes[1] != 50 {
			t.Errorf("second batch: expected 50 ids, got %d", sizes[1])
		}
		if sizes[2] != 20 {
			t.Errorf("third batch: expected 20 ids, got %d", sizes[2])
		}
		mu.Unlock()

		if len(resp.Tweets) != 120 {
			t.Fatalf("expected 120 tweets, got %d", len(resp.Tweets))
		}
		for i, tw := range resp.Tweets {
			if tw.ID != input[i] {
				t.Errorf("tweet %d: expected ID %s, got %s", i, input[i], tw.ID)
			}
		}
	})

	t.Run("middle batch failure returns error with range info", func(t *testing.T) {
		var (
			mu         sync.Mutex
			batchIndex int
		)
		api := newMock(t, func(w http.ResponseWriter, r *http.Request) {
			ids := strings.Split(r.URL.Query().Get("tweet_ids"), ",")

			mu.Lock()
			batchIndex++
			idx := batchIndex
			mu.Unlock()

			if idx == 2 {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			json.NewEncoder(w).Encode(buildTweetResponse(ids))
		})

		input := generateIDs(120)
		_, err := api.GetTweetByIDs(input)
		if err == nil {
			t.Fatal("expected error from batch failure, got nil")
		}
		if !strings.Contains(err.Error(), "[50:100)") {
			t.Errorf("error should mention batch range [50:100), got: %v", err)
		}
		if !strings.Contains(err.Error(), "50 ids") {
			t.Errorf("error should mention batch count, got: %v", err)
		}
	})
}
