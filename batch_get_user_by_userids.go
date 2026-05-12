// Doc https://docs.twitterapi.io/api-reference/endpoint/batch_get_user_by_userids

package twitterapi

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	neturl "net/url"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type BatchGetUserInfoByUserIdsURL struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
	URL         string `json:"url"`
}

type BatchGetUserInfoByUserIdsURLs struct {
	URLs []*BatchGetUserInfoByUserIdsURL `json:"urls"`
}

type BatchGetUserInfoByUserIdsEntities struct {
	Description *BatchGetUserInfoByUserIdsURLs `json:"description"`
	URL         *BatchGetUserInfoByUserIdsURLs `json:"url"`
}

type BatchGetUserInfoByUserIdsAffiliatesHighlightedLabel map[string]any

type BatchGetUserInfoByUserIdsUser struct {
	ID                         string                                               `json:"id"`
	Name                       string                                               `json:"name"`
	UserName                   string                                               `json:"userName"`
	Location                   string                                               `json:"location"`
	URL                        string                                               `json:"url"`
	Description                string                                               `json:"description"`
	Entities                   *BatchGetUserInfoByUserIdsEntities                   `json:"entities"`
	Protected                  bool                                                 `json:"protected"`
	IsVerified                 bool                                                 `json:"isVerified"`
	IsBlueVerified             bool                                                 `json:"isBlueVerified"`
	VerifiedType               *string                                              `json:"verifiedType"`
	Followers                  int                                                  `json:"followers"`
	Following                  int                                                  `json:"following"`
	FavouritesCount            int                                                  `json:"favouritesCount"`
	StatusesCount              int                                                  `json:"statusesCount"`
	MediaCount                 int                                                  `json:"mediaCount"`
	CreatedAt                  string                                               `json:"createdAt"`
	CoverPicture               string                                               `json:"coverPicture"`
	ProfilePicture             string                                               `json:"profilePicture"`
	CanDM                      bool                                                 `json:"canDm"`
	AffiliatesHighlightedLabel *BatchGetUserInfoByUserIdsAffiliatesHighlightedLabel `json:"affiliatesHighlightedLabel"`
	IsAutomated                bool                                                 `json:"isAutomated"`
	AutomatedBy                *string                                              `json:"automatedBy"`
	PinnedTweetIDs             []string                                             `json:"pinnedTweetIds"`
}

type BatchGetUserInfoByUserIdsResponse struct {
	Users   []*BatchGetUserInfoByUserIdsUser `json:"users"`
	Status  string                           `json:"status"`
	Message string                           `json:"msg"`
}

func (t *twitterApi) BatchGetUserInfoByUserIds(userIds []string) (*BatchGetUserInfoByUserIdsResponse, error) {
	if len(userIds) == 0 {
		return nil, errors.New("userIds is empty")
	}

	vals := neturl.Values{}
	vals.Set("userIds", strings.Join(userIds, ","))
	url := userTwitterDomainURI + "/batch_info_by_ids?" + vals.Encode()

	ctx1, cancel1 := context.WithTimeout(t.ctx, time.Second*10)
	defer cancel1()

	jsonData, resp, err := getDataWithHeader(ctx1, t.httpClient, url, t.headers)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			slog.Error("BatchGetUserInfoByUserIds request timed out", "url", url)
			return nil, errors.New("BatchGetUserInfoByUserIds request timed out")
		}
		slog.Error("BatchGetUserInfoByUserIds failed", "err", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		slog.Error("BatchGetUserInfoByUserIds failed", "statusCode", resp.StatusCode, "body", string(jsonData))
		return nil, errors.New("BatchGetUserInfoByUserIds failed")
	}

	response := &BatchGetUserInfoByUserIdsResponse{}
	if err = jsoniter.Unmarshal(jsonData, &response); err != nil {
		slog.Error("BatchGetUserInfoByUserIds failed", "err", err)
		return nil, err
	}
	if response.Status != "success" {
		slog.Error("BatchGetUserInfoByUserIds failed", "status", response.Status, "message", response.Message)
		return nil, errors.New(response.Message)
	}

	return response, nil
}
