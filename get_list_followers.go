// Doc https://docs.twitterapi.io/api-reference/endpoint/get_list_followers

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

type GetListFollowersFollowerAffiliatesHighlightedLabel map[string]any

type GetListFollowersFollowerProfileBioEntitiesDescriptionURL struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
	URL         string `json:"url"`
}

type GetListFollowersFollowerProfileBioEntitiesDescription struct {
	URLs []*GetListFollowersFollowerProfileBioEntitiesDescriptionURL `json:"urls"`
}

type GetListFollowersFollowerProfileBioEntitiesURLURL struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
	URL         string `json:"url"`
}

type GetListFollowersFollowerProfileBioEntitiesURL struct {
	URLs []*GetListFollowersFollowerProfileBioEntitiesURLURL `json:"urls"`
}

type GetListFollowersFollowerProfileBioEntities struct {
	Description *GetListFollowersFollowerProfileBioEntitiesDescription `json:"description"`
	URL         *GetListFollowersFollowerProfileBioEntitiesURL         `json:"url"`
}

type GetListFollowersFollowerProfileBio struct {
	Description string                                      `json:"description"`
	Entities    *GetListFollowersFollowerProfileBioEntities `json:"entities"`
}

type GetListFollowersFollower struct {
	Type                       string                                              `json:"type"`
	UserName                   string                                              `json:"userName"`
	URL                        string                                              `json:"url"`
	ID                         string                                              `json:"id"`
	Name                       string                                              `json:"name"`
	IsBlueVerified             bool                                                `json:"isBlueVerified"`
	VerifiedType               string                                              `json:"verifiedType"`
	ProfilePicture             string                                              `json:"profilePicture"`
	CoverPicture               string                                              `json:"coverPicture"`
	Description                string                                              `json:"description"`
	Location                   string                                              `json:"location"`
	Followers                  int                                                 `json:"followers"`
	Following                  int                                                 `json:"following"`
	CanDM                      bool                                                `json:"canDm"`
	CreatedAt                  string                                              `json:"createdAt"`
	FavouritesCount            int                                                 `json:"favouritesCount"`
	HasCustomTimelines         bool                                                `json:"hasCustomTimelines"`
	IsTranslator               bool                                                `json:"isTranslator"`
	MediaCount                 int                                                 `json:"mediaCount"`
	StatusesCount              int                                                 `json:"statusesCount"`
	WithheldInCountries        []string                                            `json:"withheldInCountries"`
	AffiliatesHighlightedLabel *GetListFollowersFollowerAffiliatesHighlightedLabel `json:"affiliatesHighlightedLabel"`
	PossiblySensitive          bool                                                `json:"possiblySensitive"`
	PinnedTweetIDs             []string                                            `json:"pinnedTweetIds"`
	IsAutomated                bool                                                `json:"isAutomated"`
	AutomatedBy                string                                              `json:"automatedBy"`
	Unavailable                bool                                                `json:"unavailable"`
	Message                    string                                              `json:"message"`
	UnavailableReason          string                                              `json:"unavailableReason"`
	ProfileBio                 *GetListFollowersFollowerProfileBio                 `json:"profile_bio"`
}

type GetListFollowersResponse struct {
	Followers   []*GetListFollowersFollower `json:"followers"`
	HasNextPage bool                        `json:"has_next_page"`
	NextCursor  string                      `json:"next_cursor"`
	Status      string                      `json:"status"`
	Message     string                      `json:"msg"`
}

func (t *TwitterApi) GetListFollowers(listID string, cursor *string) (*GetListFollowersResponse, error) {
	if strings.TrimSpace(listID) == "" {
		return nil, errors.New("listID is required")
	}

	vals := neturl.Values{}
	vals.Set("list_id", listID)
	if cursor != nil && *cursor != "" {
		vals.Set("cursor", *cursor)
	}
	url := listTwitterDomainURI + "/followers?" + vals.Encode()

	ctx1, cancel1 := context.WithTimeout(t.ctx, time.Second*10)
	defer cancel1()

	jsonData, resp, err := getDataWithHeader(ctx1, t.httpClient, url, t.headers)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			slog.Error("GetListFollowers request timed out", "url", url)
			return nil, errors.New("GetListFollowers request timed out")
		}
		slog.Error("GetListFollowers failed", "err", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		slog.Error("GetListFollowers failed", "statusCode", resp.StatusCode, "body", string(jsonData))
		return nil, errors.New("GetListFollowers failed")
	}

	response := &GetListFollowersResponse{}
	if err = jsoniter.Unmarshal(jsonData, &response); err != nil {
		slog.Error("GetListFollowers failed", "err", err)
		return nil, err
	}
	if response.Status != "success" {
		slog.Error("GetListFollowers failed", "status", response.Status, "message", response.Message)
		return nil, errors.New(response.Message)
	}

	return response, nil
}
