// Doc https://docs.twitterapi.io/api-reference/endpoint/add_webhook_rule

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

type addWebhookRuleRequest struct {
	Tag             string  `json:"tag"`
	Value           string  `json:"value"`
	IntervalSeconds float64 `json:"interval_seconds"`
}

type addWebhookRuleResponse struct {
	RuleID  string `json:"rule_id"`
	Status  string `json:"status"`
	Message string `json:"msg"`
}

// AddWebhookRule
func (t *TwitterApi) AddWebhookRule(tag, value string, intervalSeconds float64) (string, error) {
	if tag == "" {
		return "", errors.New("tag is empty")
	}
	if value == "" {
		return "", errors.New("value is empty")
	}
	if intervalSeconds <= 0 {
		intervalSeconds = 60
	} else if intervalSeconds < 0.05 || intervalSeconds > 86400 {
		return "", errors.New("intervalSeconds is invalid")
	}

	request := &addWebhookRuleRequest{Tag: tag, Value: value, IntervalSeconds: intervalSeconds}
	jsonData, _ := jsoniter.Marshal(request)

	ctx1, cancel1 := context.WithTimeout(t.ctx, time.Second*10)
	defer cancel1()
	headers := maps.Clone(t.headers)
	headers["Content-Type"] = "application/json"
	jsonData, resp, err := postDataWithHeader(ctx1, t.httpClient, tweetFilterDomainURI+"/add_rule", bytes.NewReader(jsonData), headers)
	if err != nil {
		slog.Error("AddWebhookRule failed", "err", err)
		return "", err
	}

	//slog.Info("AddWebhookRule response", "jsonData", string(jsonData))
	if resp.StatusCode != http.StatusOK {
		slog.Error("AddWebhookRule failed", "statusCode", resp.StatusCode, "body", string(jsonData))
		return "", errors.New("AddWebhookRule failed")
	}

	response := &addWebhookRuleResponse{}
	if err = jsoniter.Unmarshal(jsonData, &response); err != nil {
		slog.Error("AddWebhookRule failed", "err", err)
		return "", err
	}
	if response.Status != "success" {
		slog.Error("AddWebhookRule failed", "status", response.Status, "message", response.Message)
		return "", errors.New("AddWebhookRule failed")
	}
	if response.RuleID == "" {
		slog.Error("AddWebhookRule failed", "ruleID", response.RuleID, "message", response.Message)
		return "", errors.New("RuleID is empty")
	}

	return response.RuleID, nil
}
