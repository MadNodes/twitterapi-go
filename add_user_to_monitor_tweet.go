// Doc https://docs.twitterapi.io/api-reference/endpoint/add_user_to_monitor_tweet

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

type addUserToMonitorTweetRequest struct {
	XUserName string `json:"x_user_name"`
}

type addUserToMonitorTweetResponse struct {
	Status  string `json:"status"`
	Message string `json:"msg"`
}

// AddUserToMonitorTweet
func (t *TwitterApi) AddUserToMonitorTweet(xUserName string) error {
	if xUserName == "" {
		return errors.New("xUserName is empty")
	}

	request := &addUserToMonitorTweetRequest{XUserName: xUserName}
	jsonData, _ := jsoniter.Marshal(request)

	ctx1, cancel1 := context.WithTimeout(t.ctx, time.Second*10)
	defer cancel1()
	headers := maps.Clone(t.headers)
	headers["Content-Type"] = "application/json"
	jsonData, resp, err := postDataWithHeader(ctx1, t.httpClient, streamDomainURI+"/add_user_to_monitor_tweet", bytes.NewReader(jsonData), headers)
	if err != nil {
		slog.Error("AddUserToMonitorTweet failed", "err", err)
		return err
	}

	//slog.Info("AddUserToMonitorTweet response", "jsonData", string(jsonData))
	if resp.StatusCode != http.StatusOK {
		slog.Error("AddUserToMonitorTweet failed", "statusCode", resp.StatusCode, "body", string(jsonData))
		return errors.New("AddUserToMonitorTweet failed")
	}

	response := &addUserToMonitorTweetResponse{}
	if err = jsoniter.Unmarshal(jsonData, &response); err != nil {
		slog.Error("AddUserToMonitorTweet failed", "err", err)
		return err
	}
	if response.Status != "success" {
		slog.Error("AddUserToMonitorTweet failed", "status", response.Status, "message", response.Message)
		return errors.New("AddUserToMonitorTweet failed")
	}

	return nil
}
