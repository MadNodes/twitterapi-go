package twitterapi

import (
	neturl "net/url"
	"os"
	"testing"
)

func TestCheckFollowRelationship(t *testing.T) {
	client := newTestClient(t)
	// source and target must be different and both valid, otherwise the API
	// returns status=failed with no data
	sourceUserName := "xiaohu"
	targetUserName := "elonmusk"
	vals := neturl.Values{}
	vals.Set("source_user_name", sourceUserName)
	vals.Set("target_user_name", targetUserName)
	url := userTwitterDomainURI + "/check_follow_relationship?" + vals.Encode()

	raw, statusCode, err := doGet(t, client, url)
	if err != nil {
		t.Fatalf("CheckFollowRelationship request failed: %v", err)
	}
	t.Logf("RAW JSON: %s", string(raw))
	if statusCode/100 != 2 {
		t.Fatalf("CheckFollowRelationship returned non-2xx status: %d", statusCode)
	}

	var response CheckFollowRelationshipResponse
	if decodeErr := decodeJSONDisallowUnknowns(t, raw, &response); decodeErr != nil {
		t.Logf("DisallowUnknownFields decode error: %v", decodeErr)
	}
	logJSONFieldDiff(t, raw, &response)
}

func TestCheckFollowRelationship_Behavior(t *testing.T) {
	apiKey := os.Getenv("TWITTERAPI_IO_KEY")
	if apiKey == "" {
		t.Skip("TWITTERAPI_IO_KEY not set")
	}

	t.Run("empty sourceUserName returns error", func(t *testing.T) {
		x := New(apiKey)
		_, err := x.CheckFollowRelationship("", "elonmusk")
		if err == nil {
			t.Fatal("expected error for empty sourceUserName, got nil")
		}
	})

	t.Run("empty targetUserName returns error", func(t *testing.T) {
		x := New(apiKey)
		_, err := x.CheckFollowRelationship("xiaohu", "")
		if err == nil {
			t.Fatal("expected error for empty targetUserName, got nil")
		}
	})

	t.Run("valid params returns data", func(t *testing.T) {
		x := New(apiKey)
		resp, err := x.CheckFollowRelationship("xiaohu", "elonmusk")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected non-nil response")
		}
		if resp.Data == nil {
			t.Fatal("expected non-nil data")
		}
	})
}
