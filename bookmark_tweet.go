// Doc https://docs.twitterapi.io/api-reference/endpoint/bookmark_tweet_v2

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

type bookmarkTweetRequest struct {
	Cookies Cookies `json:"login_cookies"`
	Proxy   string  `json:"proxy"`
	TweetID string  `json:"tweet_id"`
}

type bookmarkTweetResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// BookmarkTweet
func (t *TwitterApi) BookmarkTweet(tweetID string) error {
	if tweetID == "" {
		return errors.New("tweetID is empty")
	}

	if t.proxy == "" {
		return errors.New("proxy is empty, please set WithProxy")
	}

	if t.cookies == "" {
		return errors.New("cookies is empty, please login first")
	}

	request := &bookmarkTweetRequest{
		Cookies: t.cookies,
		Proxy:   t.proxy,
		TweetID: tweetID,
	}

	jsonData, _ := jsoniter.Marshal(request)

	ctx1, cancel1 := context.WithTimeout(t.ctx, time.Second*5)
	defer cancel1()
	headers := maps.Clone(t.headers)
	headers["Content-Type"] = "application/json"
	jsonData, resp, err := postDataWithHeader(ctx1, t.httpClient, twitterDomainURI+"/bookmark_tweet_v2", bytes.NewReader(jsonData), headers)
	if err != nil {
		slog.Error("BookmarkTweet failed", "err", err)
		return err
	}

	slog.Info("BookmarkTweet response", "jsonData", string(jsonData))

	if resp.StatusCode != http.StatusOK {
		slog.Error("BookmarkTweet failed", "statusCode", resp.StatusCode, "body", string(jsonData))
		return errors.New("BookmarkTweet failed")
	}

	response := &bookmarkTweetResponse{}

	if err = jsoniter.Unmarshal(jsonData, &response); err != nil {
		slog.Error("BookmarkTweet failed", "err", err)
		return err
	}

	if response.Status != "success" {
		slog.Error("BookmarkTweet failed", "status", response.Status, "message", response.Message)
		return errors.New("BookmarkTweet failed")
	}

	return nil
}
