// Doc https://docs.twitterapi.io/api-reference/endpoint/follow_user_v2
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

type followUserRequest struct {
	Cookies Cookies `json:"login_cookies"`
	Proxy   string  `json:"proxy"`
	UserID  string  `json:"user_id"`
}

type followUserResponse struct {
	Status  string `json:"status"`
	Message string `json:"msg"`
}

// FollowUser
func (t *twitterApi) FollowUser(userID string) error {
	if userID == "" {
		return errors.New("userID is empty")
	}

	if t.proxy == "" {
		return errors.New("proxy is empty, please set WithProxy")
	}

	if t.cookies == "" {
		return errors.New("cookies is empty, please login first")
	}

	request := &followUserRequest{
		Cookies: t.cookies,
		Proxy:   t.proxy,
		UserID:  userID,
	}

	jsonData, _ := jsoniter.Marshal(request)

	ctx1, cancel1 := context.WithTimeout(t.ctx, time.Second*5)
	defer cancel1()
	headers := maps.Clone(t.headers)
	headers["Content-Type"] = "application/json"
	jsonData, resp, err := postDataWithHeader(ctx1, t.httpClient, twitterDomainURI+"/follow_user_v2", bytes.NewReader(jsonData), headers)
	if err != nil {
		slog.Error("FollowUser failed", "err", err)
		return err
	}

	slog.Info("FollowUser response", "jsonData", string(jsonData))

	if resp.StatusCode != http.StatusOK {
		slog.Error("FollowUser failed", "statusCode", resp.StatusCode, "body", string(jsonData))
		return errors.New("FollowUser failed")
	}

	response := &followUserResponse{}

	if err = jsoniter.Unmarshal(jsonData, &response); err != nil {
		slog.Error("FollowUser failed", "err", err)
		return err
	}

	if response.Status != "success" {
		slog.Error("FollowUser failed", "status", response.Status, "message", response.Message)
		return errors.New("FollowUser failed")
	}

	return nil
}
