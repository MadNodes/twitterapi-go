// Doc https://docs.twitterapi.io/api-reference/endpoint/get_trends

package twitterapi

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	neturl "net/url"
	"strconv"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type GetTrendsTrendTarget struct {
	Query string `json:"query"`
}

type GetTrendsTrendInner struct {
	Name            string                `json:"name"`
	Target          *GetTrendsTrendTarget `json:"target"`
	Rank            int                   `json:"rank"`
	MetaDescription string                `json:"meta_description"`
}

type GetTrendsTrend struct {
	Trend *GetTrendsTrendInner `json:"trend"` // fixed: each element wraps its trend in this key
}

type GetTrendsMetadataWoeid struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type GetTrendsMetadata struct {
	Timestamp             int64                   `json:"timestamp"`
	RefreshIntervalMillis int                     `json:"refresh_interval_millis"`
	Woeid                 *GetTrendsMetadataWoeid `json:"woeid"`
	ContextMode           string                  `json:"context_mode"`
}

type GetTrendsResponse struct {
	Trends   []*GetTrendsTrend  `json:"trends"`
	Status   string             `json:"status"`
	Message  string             `json:"msg"`
	Metadata *GetTrendsMetadata `json:"metadata"` // added: from live response
}

func (t *twitterApi) GetTrends(woeid int, count *int, the *string) (*GetTrendsResponse, error) {
	if woeid <= 0 {
		return nil, errors.New("woeid is empty")
	}

	vals := neturl.Values{}
	vals.Set("woeid", strconv.Itoa(woeid))
	if count != nil {
		vals.Set("count", strconv.Itoa(*count))
	}
	if the != nil && *the != "" {
		vals.Set("The", *the)
	}
	url := twitterDomainURI + "/trends?" + vals.Encode()

	ctx1, cancel1 := context.WithTimeout(t.ctx, time.Second*10)
	defer cancel1()

	jsonData, resp, err := getDataWithHeader(ctx1, t.httpClient, url, t.headers)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			slog.Error("GetTrends request timed out", "url", url)
			return nil, errors.New("GetTrends request timed out")
		}
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
		return nil, errors.New(response.Message)
	}

	return response, nil
}
