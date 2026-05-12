package twitterapi

import (
	neturl "net/url"
	"os"
	"testing"
)

func TestGetCommunityModerators(t *testing.T) {
	client := newTestClient(t)
	// real community ID that exposes moderators
	communityID := "1512879283559141381"
	vals := neturl.Values{}
	vals.Set("community_id", communityID)
	url := twitterDomainURI + "/community/moderators?" + vals.Encode()

	raw, statusCode, err := doGet(t, client, url)
	if err != nil {
		t.Fatalf("GetCommunityModerators request failed: %v", err)
	}
	t.Logf("RAW JSON: %s", string(raw))
	if statusCode/100 != 2 {
		t.Fatalf("GetCommunityModerators returned non-2xx status: %d", statusCode)
	}

	var response GetCommunityModeratorsResponse
	if decodeErr := decodeJSONDisallowUnknowns(t, raw, &response); decodeErr != nil {
		t.Logf("DisallowUnknownFields decode error: %v", decodeErr)
	}
	logJSONFieldDiff(t, raw, &response)
}

func TestGetCommunityModerators_Behavior(t *testing.T) {
	apiKey := os.Getenv("TWITTERAPI_IO_KEY")
	if apiKey == "" {
		t.Skip("TWITTERAPI_IO_KEY not set")
	}

	t.Run("empty communityID returns error", func(t *testing.T) {
		x := New(apiKey)
		_, err := x.GetCommunityModerators("", nil)
		if err == nil {
			t.Fatal("expected error for empty communityID, got nil")
		}
	})

	t.Run("valid communityID returns data", func(t *testing.T) {
		x := New(apiKey)
		resp, err := x.GetCommunityModerators("1512879283559141381", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected non-nil response")
		}
		if len(resp.Moderators) == 0 {
			t.Fatal("expected at least one moderator")
		}
	})
}
