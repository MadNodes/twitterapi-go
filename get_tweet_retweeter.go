// Doc https://docs.twitterapi.io/api-reference/endpoint/get_tweet_retweeter

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

type GetTweetRetweeterUserAffiliatesHighlightedLabel map[string]any

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
}

func (t *twitterApi) GetTweetRetweeter(tweetID string, cursor *string) (*GetTweetRetweeterResponse, error) {
	if strings.TrimSpace(tweetID) == "" {
		return nil, errors.New("tweetID is required")
	}

	vals := neturl.Values{}
	vals.Set("tweetId", tweetID)
	if cursor != nil && *cursor != "" {
		vals.Set("cursor", *cursor)
	}
	url := twitterDomainURI + "/tweet/retweeters?" + vals.Encode()

	ctx1, cancel1 := context.WithTimeout(t.ctx, time.Second*10)
	defer cancel1()

	jsonData, resp, err := getDataWithHeader(ctx1, t.httpClient, url, t.headers)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			slog.Error("GetTweetRetweeter request timed out", "url", url)
			return nil, errors.New("GetTweetRetweeter request timed out")
		}
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

	return response, nil
}
