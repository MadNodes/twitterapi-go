package twitterapi

import (
	neturl "net/url"
	"os"
	"testing"
)

func TestGetCommunityMembers(t *testing.T) {
	client := newTestClient(t)
	vals := neturl.Values{}
	vals.Set("community_id", testCommunityID)
	url := twitterDomainURI + "/community/members?" + vals.Encode()

	raw, statusCode, err := doGet(t, client, url)
	if err != nil {
		t.Fatalf("GetCommunityMembers request failed: %v", err)
	}
	t.Logf("RAW JSON: %s", string(raw))
	if statusCode/100 != 2 {
		t.Fatalf("GetCommunityMembers returned non-2xx status: %d", statusCode)
	}

	var response GetCommunityMembersResponse
	if decodeErr := decodeJSONDisallowUnknowns(t, raw, &response); decodeErr != nil {
		t.Logf("DisallowUnknownFields decode error: %v", decodeErr)
	}
	logJSONFieldDiff(t, raw, &response)
}

func TestGetCommunityMembers_Behavior(t *testing.T) {
	apiKey := os.Getenv("TWITTERAPI_IO_KEY")
	if apiKey == "" {
		t.Skip("TWITTERAPI_IO_KEY not set")
	}

	t.Run("empty communityID returns error", func(t *testing.T) {
		x := New(apiKey)
		_, err := x.GetCommunityMembers("", nil)
		if err == nil {
			t.Fatal("expected error for empty communityID, got nil")
		}
	})

	t.Run("valid communityID returns data", func(t *testing.T) {
		x := New(apiKey)
		resp, err := x.GetCommunityMembers("1512879283559141381", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected non-nil response")
		}
	})
}
