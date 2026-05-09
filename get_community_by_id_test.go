package twitterapi

import (
	"os"
	"testing"
)

func TestGetCommunityByID(t *testing.T) {
	client := newTestClient(t)
	url := twitterDomainURI + "/community/info?community_id=" + testCommunityID

	raw, statusCode, err := doGet(t, client, url)
	if err != nil {
		t.Fatalf("GetCommunityByID request failed: %v", err)
	}
	t.Logf("RAW JSON: %s", string(raw))
	if statusCode/100 != 2 {
		t.Fatalf("GetCommunityByID returned non-2xx status: %d", statusCode)
	}

	var response GetCommunityByIDResponse
	if decodeErr := decodeJSONDisallowUnknowns(t, raw, &response); decodeErr != nil {
		t.Logf("DisallowUnknownFields decode error: %v", decodeErr)
	}
	logJSONFieldDiff(t, raw, &response)
}

func TestGetCommunityByID_Behavior(t *testing.T) {
	apiKey := os.Getenv("TWITTERAPI_IO_KEY")
	if apiKey == "" {
		t.Skip("TWITTERAPI_IO_KEY not set")
	}

	t.Run("empty communityID returns error", func(t *testing.T) {
		x := New(apiKey)
		_, err := x.GetCommunityByID("")
		if err == nil {
			t.Fatal("expected error for empty communityID, got nil")
		}
	})

	t.Run("valid communityID returns data", func(t *testing.T) {
		x := New(apiKey)
		resp, err := x.GetCommunityByID("1512879283559141381")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected non-nil response")
		}
	})
}
