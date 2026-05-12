package twitterapi

import (
	neturl "net/url"
	"os"
	"testing"
)

func TestGetUserAbout(t *testing.T) {
	client := newTestClient(t)
	vals := neturl.Values{}
	vals.Set("userName", testUserName)
	url := userTwitterDomainURI + "_about?" + vals.Encode()

	raw, statusCode, err := doGet(t, client, url)
	if err != nil {
		t.Fatalf("GetUserAbout request failed: %v", err)
	}
	t.Logf("RAW JSON: %s", string(raw))
	if statusCode/100 != 2 {
		t.Fatalf("GetUserAbout returned non-2xx status: %d", statusCode)
	}

	var response GetUserAboutResponse
	if decodeErr := decodeJSONDisallowUnknowns(t, raw, &response); decodeErr != nil {
		t.Logf("DisallowUnknownFields decode error: %v", decodeErr)
	}
	logJSONFieldDiff(t, raw, &response)
}

func TestGetUserAbout_Behavior(t *testing.T) {
	apiKey := os.Getenv("TWITTERAPI_IO_KEY")
	if apiKey == "" {
		t.Skip("TWITTERAPI_IO_KEY not set")
	}

	t.Run("nil userName returns error", func(t *testing.T) {
		x := New(apiKey)
		_, err := x.GetUserAbout(nil)
		if err == nil {
			t.Fatal("expected error for nil userName, got nil")
		}
	})

	t.Run("valid userName returns data", func(t *testing.T) {
		x := New(apiKey)
		un := testUserName
		resp, err := x.GetUserAbout(&un)
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
