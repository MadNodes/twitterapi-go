package twitterapi

import (
	"testing"
)

func TestLogin(t *testing.T) {
	if xApiKey == "" || username == "" || email == "" || password == "" {
		t.Skip("login requires xApiKey, username, email, password, and proxy to be set in twitterapi_test.go")
	}
	x := New(xApiKey, WithProxy(proxy))

	if err := x.Login(username, email, password, nil); err != nil {
		t.Fatal(err)
	}
	t.Log(x.cookies)

}
