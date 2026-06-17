// Doc https://docs.twitterapi.io/api-reference/endpoint/check_follow_relationship

package twitterapi

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

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

func (t *TwitterApi) CheckFollowRelationship(sourceUserName, targetUserName string) (*CheckFollowRelationshipResponse, error) {
	if strings.TrimSpace(sourceUserName) == "" {
		return nil, errors.New("sourceUserName is required")
	}
	if strings.TrimSpace(targetUserName) == "" {
		return nil, errors.New("targetUserName is required")
	}

	url := userTwitterDomainURI + "/check_follow_relationship?source_user_name=" + sourceUserName + "&target_user_name=" + targetUserName

	jsonData, resp, err := t.getDataWithHeader(t.ctx, url, t.headers)
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
