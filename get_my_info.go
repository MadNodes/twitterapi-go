// Doc https://docs.twitterapi.io/api-reference/endpoint/get_my_info

package twitterapi

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type GetMyInfoResponse struct {
	RechargeCredits int    `json:"recharge_credits"`
	Status          string `json:"status"`
	Message         string `json:"msg"`
}

func (t *twitterApi) GetMyInfo() (*GetMyInfoResponse, error) {
	url := oapiDomainURI + "/my/info"

	ctx1, cancel1 := context.WithTimeout(t.ctx, time.Second*10)
	defer cancel1()

	jsonData, resp, err := getDataWithHeader(ctx1, t.httpClient, url, t.headers)
	if err != nil {
		slog.Error("GetMyInfo failed", "err", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		slog.Error("GetMyInfo failed", "statusCode", resp.StatusCode, "body", string(jsonData))
		return nil, errors.New("GetMyInfo failed")
	}

	response := &GetMyInfoResponse{}
	if err = jsoniter.Unmarshal(jsonData, &response); err != nil {
		slog.Error("GetMyInfo failed", "err", err)
		return nil, err
	}
	if response.Status != "success" {
		slog.Error("GetMyInfo failed", "status", response.Status, "message", response.Message)
		return nil, errors.New("GetMyInfo failed")
	}

	return response, nil
}
