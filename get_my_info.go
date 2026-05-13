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
	RechargeCredits   int `json:"recharge_credits"`
	TotalBonusCredits int `json:"total_bonus_credits"` // added: from live response
}

func (t *TwitterApi) GetMyInfo() (*GetMyInfoResponse, error) {
	url := oapiDomainURI + "/my/info"

	ctx1, cancel1 := context.WithTimeout(t.ctx, time.Second*10)
	defer cancel1()

	jsonData, resp, err := getDataWithHeader(ctx1, t.httpClient, url, t.headers)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			slog.Error("GetMyInfo request timed out", "url", url)
			return nil, errors.New("GetMyInfo request timed out")
		}
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

	return response, nil
}
