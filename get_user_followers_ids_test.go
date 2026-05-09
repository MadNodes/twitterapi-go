package twitterapi

import (
	"os"
	"strconv"
	"testing"
)

func TestGetUserFollowersIDs(t *testing.T) {
	client := newTestClient(t)
	count := 10
	// use xiaohu since Twitter/X doesn't expose this endpoint cleanly
	userName := "xiaohu"
	url := userTwitterDomainURI + "/followers_ids?userName=" + userName + "&count=" + strconv.Itoa(count)

	raw, statusCode, err := doGet(t, client, url)
	if err != nil {
		t.Fatalf("GetUserFollowersIDs request failed: %v", err)
	}
	t.Logf("RAW JSON: %s", string(raw))
	if statusCode/100 != 2 {
		t.Fatalf("GetUserFollowersIDs returned non-2xx status: %d", statusCode)
	}

	var response GetUserFollowersIDsResponse
	if decodeErr := decodeJSONDisallowUnknowns(t, raw, &response); decodeErr != nil {
		t.Logf("DisallowUnknownFields decode error: %v", decodeErr)
	}
	logJSONFieldDiff(t, raw, &response)
}

func TestGetUserFollowersIDs_Behavior(t *testing.T) {
	apiKey := os.Getenv("TWITTERAPI_IO_KEY")
	if apiKey == "" {
		t.Skip("TWITTERAPI_IO_KEY not set")
	}

	t.Run("empty userName and userId returns error", func(t *testing.T) {
		x := New(apiKey)
		_, err := x.GetUserFollowersIDs(nil, nil, nil, nil)
		if err == nil {
			t.Fatal("expected error for empty params, got nil")
		}
	})

	t.Run("valid userName returns data", func(t *testing.T) {
		x := New(apiKey)
		un := "xiaohu"
		resp, err := x.GetUserFollowersIDs(&un, nil, nil, nil)
		if err != nil {
			t.Skipf("endpoint unstable, skipping: %v", err)
		}
		if resp == nil {
			t.Fatal("expected non-nil response")
		}
		if len(resp.IDs) == 0 {
			t.Fatal("expected at least one follower ID")
		}
	})
}
