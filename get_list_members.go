// Doc https://docs.twitterapi.io/api-reference/endpoint/get_list_members

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

type GetListMembersMemberAffiliatesHighlightedLabel map[any]any

type GetListMembersMemberProfileBioEntitiesDescriptionURL struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
	URL         string `json:"url"`
}

type GetListMembersMemberProfileBioEntitiesDescription struct {
	URLs []*GetListMembersMemberProfileBioEntitiesDescriptionURL `json:"urls"`
}

type GetListMembersMemberProfileBioEntitiesURLURL struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
	URL         string `json:"url"`
}

type GetListMembersMemberProfileBioEntitiesURL struct {
	URLs []*GetListMembersMemberProfileBioEntitiesURLURL `json:"urls"`
}

type GetListMembersMemberProfileBioEntities struct {
	Description *GetListMembersMemberProfileBioEntitiesDescription `json:"description"`
	URL         *GetListMembersMemberProfileBioEntitiesURL         `json:"url"`
}

type GetListMembersMemberProfileBio struct {
	Description string                                  `json:"description"`
	Entities    *GetListMembersMemberProfileBioEntities `json:"entities"`
}

type GetListMembersMember struct {
	Type                       string                                          `json:"type"`
	UserName                   string                                          `json:"userName"`
	URL                        string                                          `json:"url"`
	ID                         string                                          `json:"id"`
	Name                       string                                          `json:"name"`
	IsBlueVerified             bool                                            `json:"isBlueVerified"`
	VerifiedType               string                                          `json:"verifiedType"`
	ProfilePicture             string                                          `json:"profilePicture"`
	CoverPicture               string                                          `json:"coverPicture"`
	Description                string                                          `json:"description"`
	Location                   string                                          `json:"location"`
	Followers                  int                                             `json:"followers"`
	Following                  int                                             `json:"following"`
	CanDM                      bool                                            `json:"canDm"`
	CreatedAt                  string                                          `json:"createdAt"`
	FavouritesCount            int                                             `json:"favouritesCount"`
	HasCustomTimelines         bool                                            `json:"hasCustomTimelines"`
	IsTranslator               bool                                            `json:"isTranslator"`
	MediaCount                 int                                             `json:"mediaCount"`
	StatusesCount              int                                             `json:"statusesCount"`
	WithheldInCountries        []string                                        `json:"withheldInCountries"`
	AffiliatesHighlightedLabel *GetListMembersMemberAffiliatesHighlightedLabel `json:"affiliatesHighlightedLabel"`
	PossiblySensitive          bool                                            `json:"possiblySensitive"`
	PinnedTweetIDs             []string                                        `json:"pinnedTweetIds"`
	IsAutomated                bool                                            `json:"isAutomated"`
	AutomatedBy                string                                          `json:"automatedBy"`
	Unavailable                bool                                            `json:"unavailable"`
	Message                    string                                          `json:"message"`
	UnavailableReason          string                                          `json:"unavailableReason"`
	ProfileBio                 *GetListMembersMemberProfileBio                 `json:"profile_bio"`
}

type GetListMembersResponse struct {
	Members     []*GetListMembersMember `json:"members"`
	HasNextPage bool                    `json:"has_next_page"`
	NextCursor  string                  `json:"next_cursor"`
	Status      string                  `json:"status"`
	Message     string                  `json:"msg"`
}

func (t *twitterApi) GetListMembers(listID string, cursor *string) (*GetListMembersResponse, error) {
	if listID == "" {
		return nil, errors.New("list_id is empty")
	}

	queryParts := []string{}
	queryParts = append(queryParts, "list_id="+listID)
	if cursor != nil && *cursor != "" {
		queryParts = append(queryParts, "cursor="+*cursor)
	}
	url := listTwitterDomainURI + "/members"
	if len(queryParts) > 0 {
		url += "?" + strings.Join(queryParts, "&")
	}

	ctx1, cancel1 := context.WithTimeout(t.ctx, time.Second*10)
	defer cancel1()

	jsonData, resp, err := getDataWithHeader(ctx1, t.httpClient, url, t.headers)
	if err != nil {
		slog.Error("GetListMembers failed", "err", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		slog.Error("GetListMembers failed", "statusCode", resp.StatusCode, "body", string(jsonData))
		return nil, errors.New("GetListMembers failed")
	}

	response := &GetListMembersResponse{}
	if err = jsoniter.Unmarshal(jsonData, &response); err != nil {
		slog.Error("GetListMembers failed", "err", err)
		return nil, err
	}
	if response.Status != "success" {
		slog.Error("GetListMembers failed", "status", response.Status, "message", response.Message)
		return nil, errors.New("GetListMembers failed")
	}

	return response, nil
}
