// Doc https://docs.twitterapi.io/api-reference/endpoint/get_tweet_retweeter

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

type GetTweetRetweeterUserAffiliatesHighlightedLabel map[any]any

type GetTweetRetweeterUserProfileBioEntitiesDescriptionURL struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
	URL         string `json:"url"`
}

type GetTweetRetweeterUserProfileBioEntitiesDescription struct {
	URLs []*GetTweetRetweeterUserProfileBioEntitiesDescriptionURL `json:"urls"`
}

type GetTweetRetweeterUserProfileBioEntitiesURLURL struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
	URL         string `json:"url"`
}

type GetTweetRetweeterUserProfileBioEntitiesURL struct {
	URLs []*GetTweetRetweeterUserProfileBioEntitiesURLURL `json:"urls"`
}

type GetTweetRetweeterUserProfileBioEntities struct {
	Description *GetTweetRetweeterUserProfileBioEntitiesDescription `json:"description"`
	URL         *GetTweetRetweeterUserProfileBioEntitiesURL         `json:"url"`
}

type GetTweetRetweeterUserProfileBio struct {
	Description string                                   `json:"description"`
	Entities    *GetTweetRetweeterUserProfileBioEntities `json:"entities"`
}

type GetTweetRetweeterUser struct {
	Type                       string                                           `json:"type"`
	UserName                   string                                           `json:"userName"`
	URL                        string                                           `json:"url"`
	ID                         string                                           `json:"id"`
	Name                       string                                           `json:"name"`
	IsBlueVerified             bool                                             `json:"isBlueVerified"`
	VerifiedType               string                                           `json:"verifiedType"`
	ProfilePicture             string                                           `json:"profilePicture"`
	CoverPicture               string                                           `json:"coverPicture"`
	Description                string                                           `json:"description"`
	Location                   string                                           `json:"location"`
	Followers                  int                                              `json:"followers"`
	Following                  int                                              `json:"following"`
	CanDM                      bool                                             `json:"canDm"`
	CreatedAt                  string                                           `json:"createdAt"`
	FavouritesCount            int                                              `json:"favouritesCount"`
	HasCustomTimelines         bool                                             `json:"hasCustomTimelines"`
	IsTranslator               bool                                             `json:"isTranslator"`
	MediaCount                 int                                              `json:"mediaCount"`
	StatusesCount              int                                              `json:"statusesCount"`
	WithheldInCountries        []string                                         `json:"withheldInCountries"`
	AffiliatesHighlightedLabel *GetTweetRetweeterUserAffiliatesHighlightedLabel `json:"affiliatesHighlightedLabel"`
	PossiblySensitive          bool                                             `json:"possiblySensitive"`
	PinnedTweetIDs             []string                                         `json:"pinnedTweetIds"`
	IsAutomated                bool                                             `json:"isAutomated"`
	AutomatedBy                string                                           `json:"automatedBy"`
	Unavailable                bool                                             `json:"unavailable"`
	Message                    string                                           `json:"message"`
	UnavailableReason          string                                           `json:"unavailableReason"`
	ProfileBio                 *GetTweetRetweeterUserProfileBio                 `json:"profile_bio"`
}

type GetTweetRetweeterResponse struct {
	Users       []*GetTweetRetweeterUser `json:"users"`
	HasNextPage bool                     `json:"has_next_page"`
	NextCursor  string                   `json:"next_cursor"`
	Status      string                   `json:"status"`
	Message     string                   `json:"message"`
}

func (t *twitterApi) GetTweetRetweeter(tweetID string, cursor *string) (*GetTweetRetweeterResponse, error) {
	if tweetID == "" {
		return nil, errors.New("tweetId is empty")
	}

	queryParts := []string{}
	queryParts = append(queryParts, "tweetId="+tweetID)
	if cursor != nil && *cursor != "" {
		queryParts = append(queryParts, "cursor="+*cursor)
	}
	url := twitterDomainURI + "/tweet/retweeters"
	if len(queryParts) > 0 {
		url += "?" + strings.Join(queryParts, "&")
	}

	ctx1, cancel1 := context.WithTimeout(t.ctx, time.Second*10)
	defer cancel1()

	jsonData, resp, err := getDataWithHeader(ctx1, t.httpClient, url, t.headers)
	if err != nil {
		slog.Error("GetTweetRetweeter failed", "err", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		slog.Error("GetTweetRetweeter failed", "statusCode", resp.StatusCode, "body", string(jsonData))
		return nil, errors.New("GetTweetRetweeter failed")
	}

	response := &GetTweetRetweeterResponse{}
	if err = jsoniter.Unmarshal(jsonData, &response); err != nil {
		slog.Error("GetTweetRetweeter failed", "err", err)
		return nil, err
	}
	if response.Status != "success" {
		slog.Error("GetTweetRetweeter failed", "status", response.Status, "message", response.Message)
		return nil, errors.New("GetTweetRetweeter failed")
	}

	return response, nil
}
