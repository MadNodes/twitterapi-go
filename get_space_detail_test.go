package twitterapi

import (
	"os"
	"testing"
)

func TestGetSpaceDetail(t *testing.T) {
	t.Skip("spaces are ephemeral, skip in CI")
	client := newTestClient(t)
	url := twitterDomainURI + "/spaces/detail?space_id=" + testSpaceID

	raw, statusCode, err := doGet(t, client, url)
	if err != nil {
		t.Fatalf("GetSpaceDetail request failed: %v", err)
	}
	t.Logf("RAW JSON: %s", string(raw))
	if statusCode/100 != 2 {
		t.Fatalf("GetSpaceDetail returned non-2xx status: %d", statusCode)
	}

	var response GetSpaceDetailResponse
	if decodeErr := decodeJSONDisallowUnknowns(t, raw, &response); decodeErr != nil {
		t.Logf("DisallowUnknownFields decode error: %v", decodeErr)
	}
	logJSONFieldDiff(t, raw, &response)
}

func TestGetSpaceDetail_Behavior(t *testing.T) {
	apiKey := os.Getenv("TWITTERAPI_IO_KEY")
	if apiKey == "" {
		t.Skip("TWITTERAPI_IO_KEY not set")
	}

	t.Run("nil spaceID returns error", func(t *testing.T) {
		x := New(apiKey)
		_, err := x.GetSpaceDetail(nil)
		if err == nil {
			t.Fatal("expected error for nil spaceID, got nil")
		}
	})

	t.Run("valid spaceID returns data or skips", func(t *testing.T) {
		x := New(apiKey)
		sid := testSpaceID
		resp, err := x.GetSpaceDetail(&sid)
		if err != nil {
			t.Skipf("endpoint unstable, skipping: %v", err)
		}
		if resp == nil {
			t.Fatal("expected non-nil response")
		}
	})
}
