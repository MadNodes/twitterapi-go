// Doc https://docs.twitterapi.io/api-reference/endpoint/get_user_to_monitor_tweet

package twitterapi

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type MonitoredUser struct {
	IDForUser                  string `json:"id_for_user"`
	XUserID                    int    `json:"x_user_id"`
	XUserName                  string `json:"x_user_name"`
	XUserScreenName            string `json:"x_user_screen_name"`
	IsMonitorTweet             int    `json:"is_monitor_tweet"`
	IsMonitorProfile           int    `json:"is_monitor_profile"`
	MonitorTweetConfigStatus   int    `json:"monitor_tweet_config_status"`
	MonitorProfileConfigStatus int    `json:"monitor_profile_config_status"`
	CreatedAt                  string `json:"created_at"`
}

type GetUsersToMonitorTweetResponse struct {
	Status  string           `json:"status"`
	Message string           `json:"msg"`
	Data    []*MonitoredUser `json:"data"`
}

// GetUsersToMonitorTweet
func (t *TwitterApi) GetUsersToMonitorTweet() (*GetUsersToMonitorTweetResponse, error) {
	url := streamDomainURI + "/get_user_to_monitor_tweet"

	ctx1, cancel1 := context.WithTimeout(t.ctx, time.Second*10)
	defer cancel1()

	jsonData, resp, err := getDataWithHeader(ctx1, t.httpClient, url, t.headers)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			slog.Error("GetUsersToMonitorTweet request timed out", "url", url)
			return nil, errors.New("GetUsersToMonitorTweet request timed out")
		}
		slog.Error("GetUsersToMonitorTweet failed", "err", err)
		return nil, err
	}

	slog.Info("GetUsersToMonitorTweet response", "jsonData", string(jsonData))
	if resp.StatusCode != http.StatusOK {
		slog.Error("GetUsersToMonitorTweet failed", "statusCode", resp.StatusCode, "body", string(jsonData))
		return nil, errors.New("GetUsersToMonitorTweet failed")
	}

	response := &GetUsersToMonitorTweetResponse{}
	if err = jsoniter.Unmarshal(jsonData, &response); err != nil {
		slog.Error("GetUsersToMonitorTweet failed", "err", err)
		return nil, err
	}
	if response.Status != "success" {
		slog.Error("GetUsersToMonitorTweet failed", "status", response.Status, "message", response.Message)
		return nil, errors.New(response.Message)
	}

	return response, nil
}
