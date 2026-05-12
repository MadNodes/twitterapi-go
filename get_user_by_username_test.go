package twitterapi

import (
	"net/url"
	"os"
	"testing"
)

func TestGetUserInfo(t *testing.T) {
	client := newTestClient(t)
	vals := url.Values{}
	vals.Set("userName", testUserName)
	url := userTwitterDomainURI + "/info?" + vals.Encode()

	raw, statusCode, err := doGet(t, client, url)
	if err != nil {
		t.Fatalf("GetUserInfo request failed: %v", err)
	}
	t.Logf("RAW JSON: %s", string(raw))
	if statusCode/100 != 2 {
		t.Fatalf("GetUserInfo returned non-2xx status: %d", statusCode)
	}

	var response GetUserInfoResponse
	if decodeErr := decodeJSONDisallowUnknowns(t, raw, &response); decodeErr != nil {
		t.Logf("DisallowUnknownFields decode error: %v", decodeErr)
	}
	logJSONFieldDiff(t, raw, &response)
}

func TestGetUserInfo_Behavior(t *testing.T) {
	apiKey := os.Getenv("TWITTERAPI_IO_KEY")
	if apiKey == "" {
		t.Skip("TWITTERAPI_IO_KEY not set")
	}

	t.Run("empty userName returns error", func(t *testing.T) {
		x := New(apiKey)
		_, err := x.GetUserInfo("")
		if err == nil {
			t.Fatal("expected error for empty userName, got nil")
		}
	})

	t.Run("valid userName returns data", func(t *testing.T) {
		x := New(apiKey)
		resp, err := x.GetUserInfo(testUserName)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected non-nil response")
		}
		if resp.Data == nil || resp.Data.ID == "" {
			t.Fatal("expected non-empty user ID in response")
		}
	})
}
