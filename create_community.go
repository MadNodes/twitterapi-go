// Doc https://docs.twitterapi.io/api-reference/endpoint/create_community_v2

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

type createCommunityRequest struct {
	Cookies     Cookies `json:"login_cookies"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Proxy       string  `json:"proxy"`
}

type createCommunityResponse struct {
	CommunityID string `json:"community_id"`
	Status      string `json:"status"`
	Message     string `json:"msg"`
}

// CreateCommunity
func (t *TwitterApi) CreateCommunity(name, description string) (string, error) {
	if name == "" {
		return "", errors.New("name is empty")
	}
	if description == "" {
		return "", errors.New("description is empty")
	}
	if t.proxy == "" {
		return "", errors.New("proxy is empty, please set WithProxy")
	}
	if t.cookies == "" {
		return "", errors.New("cookies is empty, please login first")
	}

	request := &createCommunityRequest{
		Cookies:     t.cookies,
		Name:        name,
		Description: description,
		Proxy:       t.proxy,
	}
	jsonData, _ := jsoniter.Marshal(request)

	ctx1, cancel1 := context.WithTimeout(t.ctx, time.Second*10)
	defer cancel1()
	headers := maps.Clone(t.headers)
	headers["Content-Type"] = "application/json"
	jsonData, resp, err := postDataWithHeader(ctx1, t.httpClient, twitterDomainURI+"/create_community_v2", bytes.NewReader(jsonData), headers)
	if err != nil {
		slog.Error("CreateCommunity failed", "err", err)
		return "", err
	}

	//slog.Info("CreateCommunity response", "jsonData", string(jsonData))
	if resp.StatusCode != http.StatusOK {
		slog.Error("CreateCommunity failed", "statusCode", resp.StatusCode, "body", string(jsonData))
		return "", errors.New("CreateCommunity failed")
	}

	response := &createCommunityResponse{}
	if err = jsoniter.Unmarshal(jsonData, &response); err != nil {
		slog.Error("CreateCommunity failed", "err", err)
		return "", err
	}
	if response.Status != "success" {
		slog.Error("CreateCommunity failed", "status", response.Status, "message", response.Message)
		return "", errors.New("CreateCommunity failed")
	}
	if response.CommunityID == "" {
		slog.Error("CreateCommunity failed", "communityID", response.CommunityID, "message", response.Message)
		return "", errors.New("CommunityID is empty")
	}

	return response.CommunityID, nil
}
