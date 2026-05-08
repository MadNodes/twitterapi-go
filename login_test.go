package twitterapi

import (
	"testing"
)

func TestLogin(t *testing.T) {
	x := New(xApiKey, WithProxy(proxy))

	if err := x.Login(username, email, password, nil); err != nil {
		t.Fatal(err)
	}
	t.Log(x.cookies)

}
