// Doc https://docs.twitterapi.io/api-reference/endpoint/update_webhook_rule

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

type updateWebhookRuleRequest struct {
	RuleID          string  `json:"rule_id"`
	Tag             string  `json:"tag"`
	Value           string  `json:"value"`
	IntervalSeconds float64 `json:"interval_seconds"`
	IsEffect        int     `json:"is_effect"`
}

type updateWebhookRuleResponse struct {
	Status  string `json:"status"`
	Message string `json:"msg"`
}

// UpdateWebhookRule
func (t *twitterApi) UpdateWebhookRule(ruleID, tag, value string, intervalSeconds float64, isEffect int) error {
	if ruleID == "" {
		return errors.New("ruleID is empty")
	}
	if tag == "" {
		return errors.New("tag is empty")
	}
	if value == "" {
		return errors.New("value is empty")
	}
	if intervalSeconds < 0.1 || intervalSeconds > 86400 {
		return errors.New("intervalSeconds is invalid")
	}
	if isEffect != 0 && isEffect != 1 {
		return errors.New("isEffect is invalid")
	}

	request := &updateWebhookRuleRequest{RuleID: ruleID, Tag: tag, Value: value, IntervalSeconds: intervalSeconds, IsEffect: isEffect}
	jsonData, _ := jsoniter.Marshal(request)

	ctx1, cancel1 := context.WithTimeout(t.ctx, time.Second*10)
	defer cancel1()
	headers := maps.Clone(t.headers)
	headers["Content-Type"] = "application/json"
	jsonData, resp, err := postDataWithHeader(ctx1, t.httpClient, tweetFilterDomainURI+"/update_rule", bytes.NewReader(jsonData), headers)
	if err != nil {
		slog.Error("UpdateWebhookRule failed", "err", err)
		return err
	}

	slog.Info("UpdateWebhookRule response", "jsonData", string(jsonData))
	if resp.StatusCode != http.StatusOK {
		slog.Error("UpdateWebhookRule failed", "statusCode", resp.StatusCode, "body", string(jsonData))
		return errors.New("UpdateWebhookRule failed")
	}

	response := &updateWebhookRuleResponse{}
	if err = jsoniter.Unmarshal(jsonData, &response); err != nil {
		slog.Error("UpdateWebhookRule failed", "err", err)
		return err
	}
	if response.Status != "success" {
		slog.Error("UpdateWebhookRule failed", "status", response.Status, "message", response.Message)
		return errors.New("UpdateWebhookRule failed")
	}

	return nil
}
