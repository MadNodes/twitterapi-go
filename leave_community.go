// Doc https://docs.twitterapi.io/api-reference/endpoint/leave_community_v2

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

type leaveCommunityRequest struct {
	Cookies     Cookies `json:"login_cookies"`
	CommunityID string  `json:"community_id"`
	Proxy       string  `json:"proxy"`
}

type leaveCommunityResponse struct {
	CommunityID   string `json:"community_id"`
	CommunityName string `json:"community_name"`
	Status        string `json:"status"`
	Message       string `json:"msg"`
}

// LeaveCommunity
func (t *TwitterApi) LeaveCommunity(communityID string) (string, string, error) {
	if communityID == "" {
		return "", "", errors.New("communityID is empty")
	}
	if t.proxy == "" {
		return "", "", errors.New("proxy is empty, please set WithProxy")
	}
	if t.cookies == "" {
		return "", "", errors.New("cookies is empty, please login first")
	}

	request := &leaveCommunityRequest{Cookies: t.cookies, CommunityID: communityID, Proxy: t.proxy}
	jsonData, _ := jsoniter.Marshal(request)

	ctx1, cancel1 := context.WithTimeout(t.ctx, time.Second*10)
	defer cancel1()
	headers := maps.Clone(t.headers)
	headers["Content-Type"] = "application/json"
	jsonData, resp, err := postDataWithHeader(ctx1, t.httpClient, twitterDomainURI+"/leave_community_v2", bytes.NewReader(jsonData), headers)
	if err != nil {
		slog.Error("LeaveCommunity failed", "err", err)
		return "", "", err
	}

	//slog.Info("LeaveCommunity response", "jsonData", string(jsonData))
	if resp.StatusCode != http.StatusOK {
		slog.Error("LeaveCommunity failed", "statusCode", resp.StatusCode, "body", string(jsonData))
		return "", "", errors.New("LeaveCommunity failed")
	}

	response := &leaveCommunityResponse{}
	if err = jsoniter.Unmarshal(jsonData, &response); err != nil {
		slog.Error("LeaveCommunity failed", "err", err)
		return "", "", err
	}
	if response.Status != "success" {
		slog.Error("LeaveCommunity failed", "status", response.Status, "message", response.Message)
		return "", "", errors.New("LeaveCommunity failed")
	}

	return response.CommunityID, response.CommunityName, nil
}
