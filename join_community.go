// Doc https://docs.twitterapi.io/api-reference/endpoint/join_community_v2

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

type joinCommunityRequest struct {
	Cookies     Cookies `json:"login_cookies"`
	CommunityID string  `json:"community_id"`
	Proxy       string  `json:"proxy"`
}

type joinCommunityResponse struct {
	CommunityID   string `json:"community_id"`
	CommunityName string `json:"community_name"`
	Status        string `json:"status"`
	Message       string `json:"msg"`
}

// JoinCommunity
func (t *twitterApi) JoinCommunity(communityID string) (string, string, error) {
	if communityID == "" {
		return "", "", errors.New("communityID is empty")
	}
	if t.proxy == "" {
		return "", "", errors.New("proxy is empty, please set WithProxy")
	}
	if t.cookies == "" {
		return "", "", errors.New("cookies is empty, please login first")
	}

	request := &joinCommunityRequest{Cookies: t.cookies, CommunityID: communityID, Proxy: t.proxy}
	jsonData, _ := jsoniter.Marshal(request)

	ctx1, cancel1 := context.WithTimeout(t.ctx, time.Second*10)
	defer cancel1()
	headers := maps.Clone(t.headers)
	headers["Content-Type"] = "application/json"
	jsonData, resp, err := postDataWithHeader(ctx1, t.httpClient, twitterDomainURI+"/join_community_v2", bytes.NewReader(jsonData), headers)
	if err != nil {
		slog.Error("JoinCommunity failed", "err", err)
		return "", "", err
	}

	slog.Info("JoinCommunity response", "jsonData", string(jsonData))
	if resp.StatusCode != http.StatusOK {
		slog.Error("JoinCommunity failed", "statusCode", resp.StatusCode, "body", string(jsonData))
		return "", "", errors.New("JoinCommunity failed")
	}

	response := &joinCommunityResponse{}
	if err = jsoniter.Unmarshal(jsonData, &response); err != nil {
		slog.Error("JoinCommunity failed", "err", err)
		return "", "", err
	}
	if response.Status != "success" {
		slog.Error("JoinCommunity failed", "status", response.Status, "message", response.Message)
		return "", "", errors.New("JoinCommunity failed")
	}

	return response.CommunityID, response.CommunityName, nil
}
