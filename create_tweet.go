// Doc https://docs.twitterapi.io/api-reference/endpoint/create_tweet_v2

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

type createTweetRequest struct {
	Cookies        Cookies  `json:"login_cookies"`
	Text           string   `json:"tweet_text"`
	Proxy          string   `json:"proxy"`
	ReplyToTweetID string   `json:"reply_to_tweet_id"`
	AttachmentURL  string   `json:"attachment_url"`
	CommunityID    string   `json:"community_id"`
	IsNoteTweet    bool     `json:"is_note_tweet"`
	MediaIDs       []string `json:"media_ids"`
	QuoteTweetID   string   `json:"quote_tweet_id"`
	ScheduleFor    string   `json:"schedule_for"`
}

type createTweetResponse struct {
	TweetID string `json:"tweet_id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (t *twitterApi) ReplyTweet(text string, mediaIDs []string, replyToTweetID, attachmentURL, communityID, quoteTweetID, scheduleFor *string, isNoteTweet *bool) (string, error) {
	return t.CreateTweet(text, mediaIDs, replyToTweetID, attachmentURL, communityID, quoteTweetID, scheduleFor, isNoteTweet)
}

func (t *twitterApi) CreateTweet(text string, mediaIDs []string, replyToTweetID, attachmentURL, communityID, quoteTweetID, scheduleFor *string, isNoteTweet *bool) (string, error) {
	if text == "" {
		return "", errors.New("text is empty")
	}

	if t.proxy == "" {
		return "", errors.New("proxy is empty, please set WithProxy")
	}

	if t.cookies == "" {
		return "", errors.New("cookies is empty, please login first")
	}

	request := &createTweetRequest{
		Cookies:  t.cookies,
		Proxy:    t.proxy,
		Text:     text,
		MediaIDs: mediaIDs,
	}

	if replyToTweetID != nil {
		request.ReplyToTweetID = *replyToTweetID
	}
	if attachmentURL != nil {
		request.AttachmentURL = *attachmentURL
	}
	if communityID != nil {
		request.CommunityID = *communityID
	}
	if quoteTweetID != nil {
		request.QuoteTweetID = *quoteTweetID
	}
	if scheduleFor != nil {
		request.ScheduleFor = *scheduleFor
	}
	if isNoteTweet != nil {
		request.IsNoteTweet = *isNoteTweet
	}

	jsonData, _ := jsoniter.Marshal(request)

	slog.Info("CreateTweet request", "jsonData", string(jsonData))

	ctx1, cancel1 := context.WithTimeout(t.ctx, time.Second*5)
	defer cancel1()
	headers := maps.Clone(t.headers)
	headers["Content-Type"] = "application/json"
	slog.Info("CreateTweet request headers", "headers", headers)

	jsonData, resp, err := postDataWithHeader(ctx1, t.httpClient, twitterDomainURI+"/create_tweet_v2", bytes.NewReader(jsonData), headers)
	if err != nil {
		slog.Error("CreateTweet failed", "err", err)
		return "", err
	}

	slog.Info("CreateTweet response", "jsonData", string(jsonData))

	if resp.StatusCode != http.StatusOK {
		slog.Error("CreateTweet failed", "statusCode", resp.StatusCode, "body", string(jsonData))
		return "", errors.New("CreateTweet failed")
	}

	response := &createTweetResponse{}

	if err = jsoniter.Unmarshal(jsonData, &response); err != nil {
		slog.Error("CreateTweet failed", "err", err)
		return "", err
	}

	if response.Status != "success" {
		slog.Error("CreateTweet failed", "status", response.Status, "message", response.Message)
		return "", errors.New("CreateTweet failed")
	}

	if response.TweetID == "" {
		slog.Error("CreateTweet failed", "tweetID", response.TweetID, "message", response.Message)
		return "", errors.New("TweetID is empty")
	}

	return response.TweetID, nil
}
