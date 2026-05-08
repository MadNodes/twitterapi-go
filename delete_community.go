// Doc https://docs.twitterapi.io/api-reference/endpoint/delete_community_v2

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

type deleteCommunityRequest struct {
	Cookies       Cookies `json:"login_cookies"`
	CommunityID   string  `json:"community_id"`
	CommunityName string  `json:"community_name"`
	Proxy         string  `json:"proxy"`
}

type deleteCommunityResponse struct {
	Status  string `json:"status"`
	Message string `json:"msg"`
}

// DeleteCommunity
func (t *twitterApi) DeleteCommunity(communityID, communityName string) error {
	if communityID == "" {
		return errors.New("communityID is empty")
	}
	if communityName == "" {
		return errors.New("communityName is empty")
	}
	if t.proxy == "" {
		return errors.New("proxy is empty, please set WithProxy")
	}
	if t.cookies == "" {
		return errors.New("cookies is empty, please login first")
	}

	request := &deleteCommunityRequest{
		Cookies:       t.cookies,
		CommunityID:   communityID,
		CommunityName: communityName,
		Proxy:         t.proxy,
	}
	jsonData, _ := jsoniter.Marshal(request)

	ctx1, cancel1 := context.WithTimeout(t.ctx, time.Second*10)
	defer cancel1()
	headers := maps.Clone(t.headers)
	headers["Content-Type"] = "application/json"
	jsonData, resp, err := postDataWithHeader(ctx1, t.httpClient, twitterDomainURI+"/delete_community_v2", bytes.NewReader(jsonData), headers)
	if err != nil {
		slog.Error("DeleteCommunity failed", "err", err)
		return err
	}

	slog.Info("DeleteCommunity response", "jsonData", string(jsonData))
	if resp.StatusCode != http.StatusOK {
		slog.Error("DeleteCommunity failed", "statusCode", resp.StatusCode, "body", string(jsonData))
		return errors.New("DeleteCommunity failed")
	}

	response := &deleteCommunityResponse{}
	if err = jsoniter.Unmarshal(jsonData, &response); err != nil {
		slog.Error("DeleteCommunity failed", "err", err)
		return err
	}
	if response.Status != "success" {
		slog.Error("DeleteCommunity failed", "status", response.Status, "message", response.Message)
		return errors.New("DeleteCommunity failed")
	}

	return nil
}
