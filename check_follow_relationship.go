// Doc https://docs.twitterapi.io/api-reference/endpoint/check_follow_relationship

package twitterapi

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type CheckFollowRelationshipData struct {
	Following  bool `json:"following"`
	FollowedBy bool `json:"followed_by"`
}

type CheckFollowRelationshipResponse struct {
	Status  string                       `json:"status"`
	Message string                       `json:"message"`
	Data    *CheckFollowRelationshipData `json:"data"`
}

func (t *twitterApi) CheckFollowRelationship(sourceUserName, targetUserName string) (*CheckFollowRelationshipResponse, error) {
	if strings.TrimSpace(sourceUserName) == "" {
		return nil, errors.New("sourceUserName is required")
	}
	if strings.TrimSpace(targetUserName) == "" {
		return nil, errors.New("targetUserName is required")
	}

	url := userTwitterDomainURI + "/check_follow_relationship?source_user_name=" + sourceUserName + "&target_user_name=" + targetUserName

	ctx1, cancel1 := context.WithTimeout(t.ctx, time.Second*10)
	defer cancel1()

	jsonData, resp, err := getDataWithHeader(ctx1, t.httpClient, url, t.headers)
	if err != nil {
		slog.Error("CheckFollowRelationship failed", "err", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		slog.Error("CheckFollowRelationship failed", "statusCode", resp.StatusCode, "body", string(jsonData))
		return nil, errors.New("CheckFollowRelationship failed")
	}

	response := &CheckFollowRelationshipResponse{}
	if err = jsoniter.Unmarshal(jsonData, &response); err != nil {
		slog.Error("CheckFollowRelationship failed", "err", err)
		return nil, err
	}
	if response.Status != "success" {
		slog.Error("CheckFollowRelationship failed", "status", response.Status, "message", response.Message)
		return nil, errors.New("CheckFollowRelationship failed")
	}

	return response, nil
}
