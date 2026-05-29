// Doc https://docs.twitterapi.io/api-reference/endpoint/get_webhook_rules

package twitterapi

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type WebhookRule struct {
	RuleID          string  `json:"rule_id"`
	Tag             string  `json:"tag"`
	Value           string  `json:"value"`
	IntervalSeconds float64 `json:"interval_seconds"`
}

type GetWebhookRulesResponse struct {
	Rules   []WebhookRule `json:"rules"`
	Status  string        `json:"status"`
	Message string        `json:"msg"`
}

// GetWebhookRules
func (t *TwitterApi) GetWebhookRules() (*GetWebhookRulesResponse, error) {
	url := tweetFilterDomainURI + "/get_rules"

	ctx1, cancel1 := context.WithTimeout(t.ctx, time.Second*10)
	defer cancel1()

	jsonData, resp, err := getDataWithHeader(ctx1, t.httpClient, url, t.headers)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			slog.Error("GetWebhookRules request timed out", "url", url)
			return nil, errors.New("GetWebhookRules request timed out")
		}
		slog.Error("GetWebhookRules failed", "err", err)
		return nil, err
	}

	//slog.Info("GetWebhookRules response", "jsonData", string(jsonData))
	if resp.StatusCode != http.StatusOK {
		slog.Error("GetWebhookRules failed", "statusCode", resp.StatusCode, "body", string(jsonData))
		return nil, errors.New("GetWebhookRules failed")
	}

	response := &GetWebhookRulesResponse{}
	if err = jsoniter.Unmarshal(jsonData, &response); err != nil {
		slog.Error("GetWebhookRules failed", "err", err)
		return nil, err
	}
	if response.Status != "success" {
		slog.Error("GetWebhookRules failed", "status", response.Status, "message", response.Message)
		return nil, errors.New(response.Message)
	}

	return response, nil
}
