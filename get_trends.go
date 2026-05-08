// Doc https://docs.twitterapi.io/api-reference/endpoint/get_trends

package twitterapi

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type GetTrendsTrendTarget struct {
	Query string `json:"query"`
}

type GetTrendsTrend struct {
	Name            string                `json:"name"`
	Target          *GetTrendsTrendTarget `json:"target"`
	Rank            int                   `json:"rank"`
	MetaDescription string                `json:"meta_description"`
}

type GetTrendsResponse struct {
	Trends  []*GetTrendsTrend `json:"trends"`
	Status  string            `json:"status"`
	Message string            `json:"msg"`
}

func (t *twitterApi) GetTrends(woeid int, count *int, the *string) (*GetTrendsResponse, error) {
	if woeid <= 0 {
		return nil, errors.New("woeid is empty")
	}

	queryParts := []string{}
	queryParts = append(queryParts, "woeid="+strconv.Itoa(woeid))
	if count != nil {
		queryParts = append(queryParts, "count="+strconv.Itoa(*count))
	}
	if the != nil && *the != "" {
		queryParts = append(queryParts, "The="+*the)
	}
	url := twitterDomainURI + "/trends"
	if len(queryParts) > 0 {
		url += "?" + strings.Join(queryParts, "&")
	}

	ctx1, cancel1 := context.WithTimeout(t.ctx, time.Second*10)
	defer cancel1()

	jsonData, resp, err := getDataWithHeader(ctx1, t.httpClient, url, t.headers)
	if err != nil {
		slog.Error("GetTrends failed", "err", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		slog.Error("GetTrends failed", "statusCode", resp.StatusCode, "body", string(jsonData))
		return nil, errors.New("GetTrends failed")
	}

	response := &GetTrendsResponse{}
	if err = jsoniter.Unmarshal(jsonData, &response); err != nil {
		slog.Error("GetTrends failed", "err", err)
		return nil, err
	}
	if response.Status != "success" {
		slog.Error("GetTrends failed", "status", response.Status, "message", response.Message)
		return nil, errors.New("GetTrends failed")
	}

	return response, nil
}
