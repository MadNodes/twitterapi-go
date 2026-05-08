// Doc https://docs.twitterapi.io/api-reference/endpoint/get_user_by_username

package twitterapi

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type GetUserInfoURL struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
	URL         string `json:"url"`
}

type GetUserInfoURLs struct {
	URLs []*GetUserInfoURL `json:"urls"`
}

type GetUserInfoEntities struct {
	Description *GetUserInfoURLs `json:"description"`
	URL         *GetUserInfoURLs `json:"url"`
}

type GetUserInfoAffiliatesHighlightedLabel map[any]any

type GetUserInfoUser struct {
	ID                         string                                 `json:"id"`
	Name                       string                                 `json:"name"`
	UserName                   string                                 `json:"userName"`
	Location                   string                                 `json:"location"`
	URL                        string                                 `json:"url"`
	Description                string                                 `json:"description"`
	Entities                   *GetUserInfoEntities                   `json:"entities"`
	Protected                  bool                                   `json:"protected"`
	IsVerified                 bool                                   `json:"isVerified"`
	IsBlueVerified             bool                                   `json:"isBlueVerified"`
	VerifiedType               *string                                `json:"verifiedType"`
	Followers                  int                                    `json:"followers"`
	Following                  int                                    `json:"following"`
	FavouritesCount            int                                    `json:"favouritesCount"`
	StatusesCount              int                                    `json:"statusesCount"`
	MediaCount                 int                                    `json:"mediaCount"`
	CreatedAt                  string                                 `json:"createdAt"`
	CoverPicture               string                                 `json:"coverPicture"`
	ProfilePicture             string                                 `json:"profilePicture"`
	CanDM                      bool                                   `json:"canDm"`
	AffiliatesHighlightedLabel *GetUserInfoAffiliatesHighlightedLabel `json:"affiliatesHighlightedLabel"`
	IsAutomated                bool                                   `json:"isAutomated"`
	AutomatedBy                *string                                `json:"automatedBy"`
	PinnedTweetIDs             []string                               `json:"pinnedTweetIds"`
}

type GetUserInfoResponse struct {
	Data    *GetUserInfoUser `json:"data"`
	Status  string           `json:"status"`
	Message string           `json:"msg"`
}

func (t *twitterApi) GetUserInfo(userName string) (*GetUserInfoResponse, error) {
	if userName == "" {
		return nil, errors.New("userName is empty")
	}

	url := userTwitterDomainURI + "/info?userName=" + userName

	ctx1, cancel1 := context.WithTimeout(t.ctx, time.Second*10)
	defer cancel1()

	jsonData, resp, err := getDataWithHeader(ctx1, t.httpClient, url, t.headers)
	if err != nil {
		slog.Error("GetUserInfo failed", "err", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		slog.Error("GetUserInfo failed", "statusCode", resp.StatusCode, "body", string(jsonData))
		return nil, errors.New("GetUserInfo failed")
	}

	response := &GetUserInfoResponse{}
	if err = jsoniter.Unmarshal(jsonData, &response); err != nil {
		slog.Error("GetUserInfo failed", "err", err)
		return nil, err
	}
	if response.Status != "success" {
		slog.Error("GetUserInfo failed", "status", response.Status, "message", response.Message)
		return nil, errors.New("GetUserInfo failed")
	}

	return response, nil
}
