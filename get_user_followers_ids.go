// Doc https://docs.twitterapi.io/api-reference/endpoint/get_user_followers_ids

package twitterapi

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	neturl "net/url"
	"strconv"
	"strings"
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

func (t *TwitterApi) GetUserFollowersIDs(userName, userId *string, count *int, cursor *string) (*GetUserFollowersIDsResponse, error) {
	if (userName == nil || strings.TrimSpace(*userName) == "") && (userId == nil || strings.TrimSpace(*userId) == "") {
		return nil, errors.New("userName or userId is required")
	}

	vals := neturl.Values{}
	if userName != nil && *userName != "" {
		vals.Set("userName", *userName)
	} else {
		vals.Set("userId", *userId)
	}
	if count != nil {
		vals.Set("count", strconv.Itoa(*count))
	}
	if cursor != nil && *cursor != "" {
		vals.Set("cursor", *cursor)
	}
	url := userTwitterDomainURI + "/followers_ids?" + vals.Encode()

	ctx1, cancel1 := context.WithTimeout(t.ctx, time.Second*10)
	defer cancel1()

	jsonData, resp, err := getDataWithHeader(ctx1, t.httpClient, url, t.headers)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			slog.Error("GetUserFollowersIDs request timed out", "url", url)
			return nil, errors.New("GetUserFollowersIDs request timed out")
		}
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
		return nil, errors.New(response.Message)
	}

	return response, nil
}
