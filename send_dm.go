// Doc https://docs.twitterapi.io/api-reference/endpoint/send_dm_v2

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

type sendDMRequest struct {
	Cookies          Cookies `json:"login_cookies"`
	UserID           string  `json:"user_id"`
	Text             string  `json:"text"`
	Proxy            string  `json:"proxy"`
	MediaID          string  `json:"media_id,omitempty"`
	ReplyToMessageID string  `json:"reply_to_message_id,omitempty"`
}

type sendDMResponse struct {
	MessageID string `json:"message_id"`
	Status    string `json:"status"`
	Message   string `json:"msg"`
}

// SendDM
func (t *TwitterApi) SendDM(userID, text string, mediaID, replyToMessageID *string) (string, error) {
	if userID == "" {
		return "", errors.New("userID is empty")
	}

	if text == "" {
		return "", errors.New("text is empty")
	}

	if t.proxy == "" {
		return "", errors.New("proxy is empty, please set WithProxy")
	}

	if t.cookies == "" {
		return "", errors.New("cookies is empty, please login first")
	}

	request := &sendDMRequest{
		Cookies: t.cookies,
		UserID:  userID,
		Text:    text,
		Proxy:   t.proxy,
	}

	if mediaID != nil {
		request.MediaID = *mediaID
	}
	if replyToMessageID != nil {
		request.ReplyToMessageID = *replyToMessageID
	}

	jsonData, _ := jsoniter.Marshal(request)

	ctx1, cancel1 := context.WithTimeout(t.ctx, time.Second*5)
	defer cancel1()
	headers := maps.Clone(t.headers)
	headers["Content-Type"] = "application/json"
	jsonData, resp, err := postDataWithHeader(ctx1, t.httpClient, twitterDomainURI+"/send_dm_to_user", bytes.NewReader(jsonData), headers)
	if err != nil {
		slog.Error("SendDM failed", "err", err)
		return "", err
	}

	//slog.Info("SendDM response", "jsonData", string(jsonData))

	if resp.StatusCode != http.StatusOK {
		slog.Error("SendDM failed", "statusCode", resp.StatusCode, "body", string(jsonData))
		return "", errors.New("SendDM failed")
	}

	response := &sendDMResponse{}

	if err = jsoniter.Unmarshal(jsonData, &response); err != nil {
		slog.Error("SendDM failed", "err", err)
		return "", err
	}

	if response.Status != "success" {
		slog.Error("SendDM failed", "status", response.Status, "message", response.Message)
		return "", errors.New("SendDM failed")
	}

	if response.MessageID == "" {
		slog.Error("SendDM failed", "messageID", response.MessageID, "message", response.Message)
		return "", errors.New("MessageID is empty")
	}

	return response.MessageID, nil
}
