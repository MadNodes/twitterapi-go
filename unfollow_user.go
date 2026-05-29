// Doc https://docs.twitterapi.io/api-reference/endpoint/unfollow_user_v2

package twitterapi

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"maps"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type unfollowUserRequest struct {
	Cookies Cookies `json:"login_cookies"`
	Proxy   string  `json:"proxy"`
	UserID  string  `json:"user_id"`
}

type unfollowUserResponse struct {
	Status  string `json:"status"`
	Message string `json:"msg"`
}

// UnfollowUser
func (t *TwitterApi) UnfollowUser(userID string) error {
	if userID == "" {
		return errors.New("userID is empty")
	}

	if t.proxy == "" {
		return errors.New("proxy is empty, please set WithProxy")
	}

	if t.cookies == "" {
		return errors.New("cookies is empty, please login first")
	}

	request := &unfollowUserRequest{
		Cookies: t.cookies,
		Proxy:   t.proxy,
		UserID:  userID,
	}

	jsonData, _ := jsoniter.Marshal(request)

	ctx1, cancel1 := context.WithTimeout(t.ctx, time.Second*5)
	defer cancel1()
	headers := maps.Clone(t.headers)
	headers["Content-Type"] = "application/json"
	jsonData, resp, err := postDataWithHeader(ctx1, t.httpClient, twitterDomainURI+"/unfollow_user_v2", bytes.NewReader(jsonData), headers)
	if err != nil {
		slog.Error("UnfollowUser failed", "err", err)
		return err
	}

	//slog.Info("UnfollowUser response", "jsonData", string(jsonData))

	if resp.StatusCode != http.StatusOK {
		slog.Error("UnfollowUser failed", "statusCode", resp.StatusCode, "body", string(jsonData))
		return errors.New("UnfollowUser failed")
	}

	response := &unfollowUserResponse{}

	if err = jsoniter.Unmarshal(jsonData, &response); err != nil {
		slog.Error("UnfollowUser failed", "err", err)
		return err
	}

	if response.Status != "success" {
		slog.Error("UnfollowUser failed", "status", response.Status, "message", response.Message)
		return errors.New("UnfollowUser failed")
	}

	return nil
}
