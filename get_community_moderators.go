// Doc https://docs.twitterapi.io/api-reference/endpoint/get_community_moderators

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

type GetCommunityModeratorsMemberAffiliatesHighlightedLabel map[any]any

type GetCommunityModeratorsMemberProfileBioEntitiesDescriptionURL struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
	URL         string `json:"url"`
}

type GetCommunityModeratorsMemberProfileBioEntitiesDescription struct {
	URLs []*GetCommunityModeratorsMemberProfileBioEntitiesDescriptionURL `json:"urls"`
}

type GetCommunityModeratorsMemberProfileBioEntitiesURLURL struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
	URL         string `json:"url"`
}

type GetCommunityModeratorsMemberProfileBioEntitiesURL struct {
	URLs []*GetCommunityModeratorsMemberProfileBioEntitiesURLURL `json:"urls"`
}

type GetCommunityModeratorsMemberProfileBioEntities struct {
	Description *GetCommunityModeratorsMemberProfileBioEntitiesDescription `json:"description"`
	URL         *GetCommunityModeratorsMemberProfileBioEntitiesURL         `json:"url"`
}

type GetCommunityModeratorsMemberProfileBio struct {
	Description string                                          `json:"description"`
	Entities    *GetCommunityModeratorsMemberProfileBioEntities `json:"entities"`
}

type GetCommunityModeratorsMember struct {
	Type                       string                                                  `json:"type"`
	UserName                   string                                                  `json:"userName"`
	URL                        string                                                  `json:"url"`
	ID                         string                                                  `json:"id"`
	Name                       string                                                  `json:"name"`
	IsBlueVerified             bool                                                    `json:"isBlueVerified"`
	VerifiedType               string                                                  `json:"verifiedType"`
	ProfilePicture             string                                                  `json:"profilePicture"`
	CoverPicture               string                                                  `json:"coverPicture"`
	Description                string                                                  `json:"description"`
	Location                   string                                                  `json:"location"`
	Followers                  int                                                     `json:"followers"`
	Following                  int                                                     `json:"following"`
	CanDM                      bool                                                    `json:"canDm"`
	CreatedAt                  string                                                  `json:"createdAt"`
	FavouritesCount            int                                                     `json:"favouritesCount"`
	HasCustomTimelines         bool                                                    `json:"hasCustomTimelines"`
	IsTranslator               bool                                                    `json:"isTranslator"`
	MediaCount                 int                                                     `json:"mediaCount"`
	StatusesCount              int                                                     `json:"statusesCount"`
	WithheldInCountries        []string                                                `json:"withheldInCountries"`
	AffiliatesHighlightedLabel *GetCommunityModeratorsMemberAffiliatesHighlightedLabel `json:"affiliatesHighlightedLabel"`
	PossiblySensitive          bool                                                    `json:"possiblySensitive"`
	PinnedTweetIDs             []string                                                `json:"pinnedTweetIds"`
	IsAutomated                bool                                                    `json:"isAutomated"`
	AutomatedBy                string                                                  `json:"automatedBy"`
	Unavailable                bool                                                    `json:"unavailable"`
	Message                    string                                                  `json:"message"`
	UnavailableReason          string                                                  `json:"unavailableReason"`
	ProfileBio                 *GetCommunityModeratorsMemberProfileBio                 `json:"profile_bio"`
}

type GetCommunityModeratorsResponse struct {
	Members     []*GetCommunityModeratorsMember `json:"members"`
	HasNextPage bool                            `json:"has_next_page"`
	NextCursor  string                          `json:"next_cursor"`
	Status      string                          `json:"status"`
	Message     string                          `json:"msg"`
}

func (t *twitterApi) GetCommunityModerators(communityID string, cursor *string) (*GetCommunityModeratorsResponse, error) {
	if communityID == "" {
		return nil, errors.New("community_id is empty")
	}

	queryParts := []string{}
	queryParts = append(queryParts, "community_id="+communityID)
	if cursor != nil && *cursor != "" {
		queryParts = append(queryParts, "cursor="+*cursor)
	}
	url := twitterDomainURI + "/community/moderators"
	if len(queryParts) > 0 {
		url += "?" + strings.Join(queryParts, "&")
	}

	ctx1, cancel1 := context.WithTimeout(t.ctx, time.Second*10)
	defer cancel1()

	jsonData, resp, err := getDataWithHeader(ctx1, t.httpClient, url, t.headers)
	if err != nil {
		slog.Error("GetCommunityModerators failed", "err", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		slog.Error("GetCommunityModerators failed", "statusCode", resp.StatusCode, "body", string(jsonData))
		return nil, errors.New("GetCommunityModerators failed")
	}

	response := &GetCommunityModeratorsResponse{}
	if err = jsoniter.Unmarshal(jsonData, &response); err != nil {
		slog.Error("GetCommunityModerators failed", "err", err)
		return nil, err
	}
	if response.Status != "success" {
		slog.Error("GetCommunityModerators failed", "status", response.Status, "message", response.Message)
		return nil, errors.New("GetCommunityModerators failed")
	}

	return response, nil
}
