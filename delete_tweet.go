// Doc https://docs.twitterapi.io/api-reference/endpoint/delete_tweet_v2

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

type deleteTweetRequest struct {
	Cookies Cookies `json:"login_cookies"`
	Proxy   string  `json:"proxy"`
	TweetID string  `json:"tweet_id"`
}

type deleteTweetResponse struct {
	Status  string `json:"status"`
	Message string `json:"msg"`
}

// DeleteTweet
func (t *twitterApi) DeleteTweet(tweetID string) error {
	if tweetID == "" {
		return errors.New("tweetID is empty")
	}

	if t.proxy == "" {
		return errors.New("proxy is empty, please set WithProxy")
	}

	if t.cookies == "" {
		return errors.New("cookies is empty, please login first")
	}

	request := &deleteTweetRequest{
		Cookies: t.cookies,
		Proxy:   t.proxy,
		TweetID: tweetID,
	}

	jsonData, _ := jsoniter.Marshal(request)

	ctx1, cancel1 := context.WithTimeout(t.ctx, time.Second*5)
	defer cancel1()
	headers := maps.Clone(t.headers)
	headers["Content-Type"] = "application/json"
	jsonData, resp, err := postDataWithHeader(ctx1, t.httpClient, twitterDomainURI+"/delete_tweet_v2", bytes.NewReader(jsonData), headers)
	if err != nil {
		slog.Error("DeleteTweet failed", "err", err)
		return err
	}

	slog.Info("DeleteTweet response", "jsonData", string(jsonData))

	if resp.StatusCode != http.StatusOK {
		slog.Error("DeleteTweet failed", "statusCode", resp.StatusCode, "body", string(jsonData))
		return errors.New("DeleteTweet failed")
	}

	response := &deleteTweetResponse{}

	if err = jsoniter.Unmarshal(jsonData, &response); err != nil {
		slog.Error("DeleteTweet failed", "err", err)
		return err
	}

	if response.Status != "success" {
		slog.Error("DeleteTweet failed", "status", response.Status, "message", response.Message)
		return errors.New("DeleteTweet failed")
	}

	return nil
}
