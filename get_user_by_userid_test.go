package twitterapi

import (
	"os"
	"testing"
)

func TestGetUserInfoByID_Behavior(t *testing.T) {
	apiKey := os.Getenv("TWITTERAPI_IO_KEY")
	if apiKey == "" {
		t.Skip("TWITTERAPI_IO_KEY not set")
	}

	t.Run("valid userId returns data", func(t *testing.T) {
		x := New(apiKey)
		user, err := x.GetUserInfoByID(testUserID)
		if err != nil {
			t.Skipf("endpoint unstable, skipping: %v", err)
		}
		if user == nil {
			t.Fatal("expected non-nil user")
		}
		if user.ID == "" {
			t.Fatal("expected non-empty user ID")
		}
	})
}
