package twitterapi

import (
	neturl "net/url"
	"os"
	"strconv"
	"testing"
)

func TestGetUserMentions(t *testing.T) {
	client := newTestClient(t)
	sinceTime, untilTime := last7DaysRangeUnix()
	vals := neturl.Values{}
	vals.Set("userName", testUserName)
	vals.Set("sinceTime", strconv.FormatInt(sinceTime, 10))
	vals.Set("untilTime", strconv.FormatInt(untilTime, 10))
	url := userTwitterDomainURI + "/mentions?" + vals.Encode()

	raw, statusCode, err := doGet(t, client, url)
	if err != nil {
		t.Fatalf("GetUserMentions request failed: %v", err)
	}
	t.Logf("RAW JSON: %s", string(raw))
	if statusCode/100 != 2 {
		t.Fatalf("GetUserMentions returned non-2xx status: %d", statusCode)
	}

	var response GetUserMentionsResponse
	if decodeErr := decodeJSONDisallowUnknowns(t, raw, &response); decodeErr != nil {
		t.Logf("DisallowUnknownFields decode error: %v", decodeErr)
	}
	logJSONFieldDiff(t, raw, &response)
}

func TestGetUserMentions_Behavior(t *testing.T) {
	apiKey := os.Getenv("TWITTERAPI_IO_KEY")
	if apiKey == "" {
		t.Skip("TWITTERAPI_IO_KEY not set")
	}

	t.Run("empty userName returns error", func(t *testing.T) {
		x := New(apiKey)
		_, err := x.GetUserMentions("", nil, nil, nil)
		if err == nil {
			t.Fatal("expected error for empty userName, got nil")
		}
	})

	t.Run("valid userName returns data", func(t *testing.T) {
		x := New(apiKey)
		resp, err := x.GetUserMentions(testUserName, nil, nil, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected non-nil response")
		}
	})
}
