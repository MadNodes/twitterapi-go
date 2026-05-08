// Doc https://docs.twitterapi.io/api-reference/endpoint/delete_webhook_rule

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

type deleteWebhookRuleRequest struct {
	RuleID string `json:"rule_id"`
}

type deleteWebhookRuleResponse struct {
	Status  string `json:"status"`
	Message string `json:"msg"`
}

// DeleteWebhookRule
func (t *twitterApi) DeleteWebhookRule(ruleID string) error {
	if ruleID == "" {
		return errors.New("ruleID is empty")
	}

	request := &deleteWebhookRuleRequest{RuleID: ruleID}
	jsonData, _ := jsoniter.Marshal(request)

	ctx1, cancel1 := context.WithTimeout(t.ctx, time.Second*10)
	defer cancel1()

	headers := maps.Clone(t.headers)
	headers["Content-Type"] = "application/json"

	jsonData, resp, err := deleteDataWithHeader(ctx1, t.httpClient, tweetFilterDomainURI+"/delete_rule", bytes.NewReader(jsonData), headers)
	if err != nil {
		slog.Error("DeleteWebhookRule failed", "err", err)
		return err
	}

	slog.Info("DeleteWebhookRule response", "jsonData", string(jsonData))
	if resp.StatusCode != http.StatusOK {
		slog.Error("DeleteWebhookRule failed", "statusCode", resp.StatusCode, "body", string(jsonData))
		return errors.New("DeleteWebhookRule failed")
	}

	response := &deleteWebhookRuleResponse{}
	if err = jsoniter.Unmarshal(jsonData, &response); err != nil {
		slog.Error("DeleteWebhookRule failed", "err", err)
		return err
	}
	if response.Status != "success" {
		slog.Error("DeleteWebhookRule failed", "status", response.Status, "message", response.Message)
		return errors.New("DeleteWebhookRule failed")
	}

	return nil
}
