// Doc https://docs.twitterapi.io/api-reference/endpoint/get_community_members

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

type GetCommunityMembersMemberAffiliatesHighlightedLabel map[string]any

type GetCommunityMembersMemberProfileBioEntitiesDescriptionURL struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
	URL         string `json:"url"`
}

type GetCommunityMembersMemberProfileBioEntitiesDescription struct {
	URLs []*GetCommunityMembersMemberProfileBioEntitiesDescriptionURL `json:"urls"`
}

type GetCommunityMembersMemberProfileBioEntitiesURLURL struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
	URL         string `json:"url"`
}

type GetCommunityMembersMemberProfileBioEntitiesURL struct {
	URLs []*GetCommunityMembersMemberProfileBioEntitiesURLURL `json:"urls"`
}

type GetCommunityMembersMemberProfileBioEntities struct {
	Description *GetCommunityMembersMemberProfileBioEntitiesDescription `json:"description"`
	URL         *GetCommunityMembersMemberProfileBioEntitiesURL         `json:"url"`
}

type GetCommunityMembersMemberProfileBio struct {
	Description string                                       `json:"description"`
	Entities    *GetCommunityMembersMemberProfileBioEntities `json:"entities"`
}

type GetCommunityMembersMember struct {
	Type                       string                                               `json:"type"`
	UserName                   string                                               `json:"userName"`
	URL                        string                                               `json:"url"`
	ID                         string                                               `json:"id"`
	Name                       string                                               `json:"name"`
	IsBlueVerified             bool                                                 `json:"isBlueVerified"`
	VerifiedType               string                                               `json:"verifiedType"`
	ProfilePicture             string                                               `json:"profilePicture"`
	CoverPicture               string                                               `json:"coverPicture"`
	Description                string                                               `json:"description"`
	Location                   string                                               `json:"location"`
	Followers                  int                                                  `json:"followers"`
	Following                  int                                                  `json:"following"`
	CanDM                      bool                                                 `json:"canDm"`
	CreatedAt                  string                                               `json:"createdAt"`
	FavouritesCount            int                                                  `json:"favouritesCount"`
	HasCustomTimelines         bool                                                 `json:"hasCustomTimelines"`
	IsTranslator               bool                                                 `json:"isTranslator"`
	MediaCount                 int                                                  `json:"mediaCount"`
	StatusesCount              int                                                  `json:"statusesCount"`
	WithheldInCountries        []string                                             `json:"withheldInCountries"`
	AffiliatesHighlightedLabel *GetCommunityMembersMemberAffiliatesHighlightedLabel `json:"affiliatesHighlightedLabel"`
	PossiblySensitive          bool                                                 `json:"possiblySensitive"`
	PinnedTweetIDs             []string                                             `json:"pinnedTweetIds"`
	IsAutomated                bool                                                 `json:"isAutomated"`
	AutomatedBy                string                                               `json:"automatedBy"`
	Unavailable                bool                                                 `json:"unavailable"`
	Message                    string                                               `json:"message"`
	UnavailableReason          string                                               `json:"unavailableReason"`
	ProfileBio                 *GetCommunityMembersMemberProfileBio                 `json:"profile_bio"`
}

type GetCommunityMembersResponse struct {
	Members     []*GetCommunityMembersMember `json:"members"`
	HasNext     bool                         `json:"has_next"` // added: from live response
	HasNextPage bool                         `json:"has_next_page"`
	NextCursor  string                       `json:"next_cursor"`
	Status      string                       `json:"status"`
	Message     string                       `json:"msg"`
}

func (t *TwitterApi) GetCommunityMembers(communityID string, cursor *string) (*GetCommunityMembersResponse, error) {
	if strings.TrimSpace(communityID) == "" {
		return nil, errors.New("communityID is required")
	}

	vals := neturl.Values{}
	vals.Set("community_id", communityID)
	if cursor != nil && *cursor != "" {
		vals.Set("cursor", *cursor)
	}
	url := twitterDomainURI + "/community/members?" + vals.Encode()

	ctx1, cancel1 := context.WithTimeout(t.ctx, time.Second*10)
	defer cancel1()

	jsonData, resp, err := getDataWithHeader(ctx1, t.httpClient, url, t.headers)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			slog.Error("GetCommunityMembers request timed out", "url", url)
			return nil, errors.New("GetCommunityMembers request timed out")
		}
		slog.Error("GetCommunityMembers failed", "err", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		slog.Error("GetCommunityMembers failed", "statusCode", resp.StatusCode, "body", string(jsonData))
		return nil, errors.New("GetCommunityMembers failed")
	}

	response := &GetCommunityMembersResponse{}
	if err = jsoniter.Unmarshal(jsonData, &response); err != nil {
		slog.Error("GetCommunityMembers failed", "err", err)
		return nil, err
	}
	if response.Status != "success" {
		slog.Error("GetCommunityMembers failed", "status", response.Status, "message", response.Message)
		return nil, errors.New(response.Message)
	}

	return response, nil
}
