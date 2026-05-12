package twitterapi

import (
	neturl "net/url"
	"os"
	"testing"
)

func TestGetListTimeline(t *testing.T) {
	client := newTestClient(t)
	vals := neturl.Values{}
	vals.Set("listId", testListID)
	url := listTwitterDomainURI + "/tweets_timeline?" + vals.Encode()

	raw, statusCode, err := doGet(t, client, url)
	if err != nil {
		t.Fatalf("GetListTimeline request failed: %v", err)
	}
	t.Logf("RAW JSON: %s", string(raw))
	if statusCode/100 != 2 {
		t.Fatalf("GetListTimeline returned non-2xx status: %d", statusCode)
	}

	var response ListTimelineResponse
	if decodeErr := decodeJSONDisallowUnknowns(t, raw, &response); decodeErr != nil {
		t.Logf("DisallowUnknownFields decode error: %v", decodeErr)
	}
	logJSONFieldDiff(t, raw, &response)
}

func TestGetListTimeline_Behavior(t *testing.T) {
	apiKey := os.Getenv("TWITTERAPI_IO_KEY")
	if apiKey == "" {
		t.Skip("TWITTERAPI_IO_KEY not set")
	}

	t.Run("empty listID returns error", func(t *testing.T) {
		x := New(apiKey)
		_, err := x.GetListTimeline("", nil)
		if err == nil {
			t.Fatal("expected error for empty listID, got nil")
		}
	})

	t.Run("valid listID returns data", func(t *testing.T) {
		x := New(apiKey)
		resp, err := x.GetListTimeline(testListID, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected non-nil response")
		}
	})
}
