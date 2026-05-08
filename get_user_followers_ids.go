// Doc https://docs.twitterapi.io/api-reference/endpoint/get_user_followers_ids

package twitterapi

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type GetUserFollowersIDsResponse struct {
	IDs         []string `json:"ids"`
	HasNextPage bool     `json:"has_next_page"`
	NextCursor  string   `json:"next_cursor"`
	Status      string   `json:"status"`
	Message     string   `json:"msg"`
	Code        int      `json:"code"`
}

func (t *twitterApi) GetUserFollowersIDs(userName, userId *string, count *int, cursor *string) (*GetUserFollowersIDsResponse, error) {
	if (userName == nil || *userName == "") && (userId == nil || *userId == "") {
		return nil, errors.New("userName or userId is required")
	}

	url := userTwitterDomainURI + "/followers_ids?"
	if userName != nil && *userName != "" {
		url += "userName=" + *userName
	} else {
		url += "userId=" + *userId
	}
	if count != nil {
		url += "&count=" + strconv.Itoa(*count)
	}
	if cursor != nil && *cursor != "" {
		url += "&cursor=" + *cursor
	}

	ctx1, cancel1 := context.WithTimeout(t.ctx, time.Second*10)
	defer cancel1()

	jsonData, resp, err := getDataWithHeader(ctx1, t.httpClient, url, t.headers)
	if err != nil {
		slog.Error("GetUserFollowersIDs failed", "err", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		slog.Error("GetUserFollowersIDs failed", "statusCode", resp.StatusCode, "body", string(jsonData))
		return nil, errors.New("GetUserFollowersIDs failed")
	}

	response := &GetUserFollowersIDsResponse{}
	if err = jsoniter.Unmarshal(jsonData, &response); err != nil {
		slog.Error("GetUserFollowersIDs failed", "err", err)
		return nil, err
	}
	if response.Status != "success" {
		slog.Error("GetUserFollowersIDs failed", "status", response.Status, "message", response.Message)
		return nil, errors.New("GetUserFollowersIDs failed")
	}

	return response, nil
}
