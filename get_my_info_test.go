package twitterapi

import (
	"os"
	"testing"
)

func TestGetMyInfo(t *testing.T) {
	client := newTestClient(t)
	url := oapiDomainURI + "/my/info"

	raw, statusCode, err := doGet(t, client, url)
	if err != nil {
		t.Fatalf("GetMyInfo request failed: %v", err)
	}
	t.Logf("RAW JSON: %s", string(raw))
	if statusCode/100 != 2 {
		t.Fatalf("GetMyInfo returned non-2xx status: %d", statusCode)
	}

	var response GetMyInfoResponse
	if decodeErr := decodeJSONDisallowUnknowns(t, raw, &response); decodeErr != nil {
		t.Logf("DisallowUnknownFields decode error: %v", decodeErr)
	}
	logJSONFieldDiff(t, raw, &response)
}

func TestGetMyInfo_Behavior(t *testing.T) {
	apiKey := os.Getenv("TWITTERAPI_IO_KEY")
	if apiKey == "" {
		t.Skip("TWITTERAPI_IO_KEY not set")
	}

	t.Run("returns data", func(t *testing.T) {
		x := New(apiKey)
		resp, err := x.GetMyInfo()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected non-nil response")
		}
	})
}
