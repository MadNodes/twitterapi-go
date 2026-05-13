// Doc https://docs.twitterapi.io/api-reference/endpoint/remove_user_to_monitor_tweet

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

type removeUserToMonitorTweetRequest struct {
	IDForUser string `json:"id_for_user"`
}

type removeUserToMonitorTweetResponse struct {
	Status  string `json:"status"`
	Message string `json:"msg"`
}

// RemoveUserToMonitorTweet
func (t *TwitterApi) RemoveUserToMonitorTweet(idForUser string) error {
	if idForUser == "" {
		return errors.New("idForUser is empty")
	}

	request := &removeUserToMonitorTweetRequest{IDForUser: idForUser}
	jsonData, _ := jsoniter.Marshal(request)

	ctx1, cancel1 := context.WithTimeout(t.ctx, time.Second*10)
	defer cancel1()
	headers := maps.Clone(t.headers)
	headers["Content-Type"] = "application/json"
	jsonData, resp, err := postDataWithHeader(ctx1, t.httpClient, streamDomainURI+"/remove_user_to_monitor_tweet", bytes.NewReader(jsonData), headers)
	if err != nil {
		slog.Error("RemoveUserToMonitorTweet failed", "err", err)
		return err
	}

	slog.Info("RemoveUserToMonitorTweet response", "jsonData", string(jsonData))
	if resp.StatusCode != http.StatusOK {
		slog.Error("RemoveUserToMonitorTweet failed", "statusCode", resp.StatusCode, "body", string(jsonData))
		return errors.New("RemoveUserToMonitorTweet failed")
	}

	response := &removeUserToMonitorTweetResponse{}
	if err = jsoniter.Unmarshal(jsonData, &response); err != nil {
		slog.Error("RemoveUserToMonitorTweet failed", "err", err)
		return err
	}
	if response.Status != "success" {
		slog.Error("RemoveUserToMonitorTweet failed", "status", response.Status, "message", response.Message)
		return errors.New("RemoveUserToMonitorTweet failed")
	}

	return nil
}
