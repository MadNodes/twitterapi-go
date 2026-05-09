// Doc https://docs.twitterapi.io/api-reference/endpoint/get_community_by_id

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

type GetCommunityByIDCommunityInfoPrimaryTopic struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GetCommunityByIDCommunityInfoRule struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type GetCommunityByIDCommunityInfoCreatorAffiliatesHighlightedLabel map[string]any

type GetCommunityByIDCommunityInfoCreatorProfileBioEntitiesDescriptionURL struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
	URL         string `json:"url"`
}

type GetCommunityByIDCommunityInfoCreatorProfileBioEntitiesDescription struct {
	URLs []*GetCommunityByIDCommunityInfoCreatorProfileBioEntitiesDescriptionURL `json:"urls"`
}

type GetCommunityByIDCommunityInfoCreatorProfileBioEntitiesURLURL struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
	URL         string `json:"url"`
}

type GetCommunityByIDCommunityInfoCreatorProfileBioEntitiesURL struct {
	URLs []*GetCommunityByIDCommunityInfoCreatorProfileBioEntitiesURLURL `json:"urls"`
}

type GetCommunityByIDCommunityInfoCreatorProfileBioEntities struct {
	Description *GetCommunityByIDCommunityInfoCreatorProfileBioEntitiesDescription `json:"description"`
	URL         *GetCommunityByIDCommunityInfoCreatorProfileBioEntitiesURL         `json:"url"`
}

type GetCommunityByIDCommunityInfoCreatorProfileBio struct {
	Description string                                                  `json:"description"`
	Entities    *GetCommunityByIDCommunityInfoCreatorProfileBioEntities `json:"entities"`
}

type GetCommunityByIDCommunityInfoCreator struct {
	Type                       string                                                          `json:"type"`
	UserName                   string                                                          `json:"userName"`
	URL                        string                                                          `json:"url"`
	ID                         string                                                          `json:"id"`
	Name                       string                                                          `json:"name"`
	IsBlueVerified             bool                                                            `json:"isBlueVerified"`
	VerifiedType               string                                                          `json:"verifiedType"`
	ProfilePicture             string                                                          `json:"profilePicture"`
	CoverPicture               string                                                          `json:"coverPicture"`
	Description                string                                                          `json:"description"`
	Location                   string                                                          `json:"location"`
	Followers                  int                                                             `json:"followers"`
	Following                  int                                                             `json:"following"`
	CanDM                      bool                                                            `json:"canDm"`
	CreatedAt                  string                                                          `json:"createdAt"`
	FavouritesCount            int                                                             `json:"favouritesCount"`
	HasCustomTimelines         bool                                                            `json:"hasCustomTimelines"`
	IsTranslator               bool                                                            `json:"isTranslator"`
	MediaCount                 int                                                             `json:"mediaCount"`
	StatusesCount              int                                                             `json:"statusesCount"`
	WithheldInCountries        []string                                                        `json:"withheldInCountries"`
	AffiliatesHighlightedLabel *GetCommunityByIDCommunityInfoCreatorAffiliatesHighlightedLabel `json:"affiliatesHighlightedLabel"`
	PossiblySensitive          bool                                                            `json:"possiblySensitive"`
	PinnedTweetIDs             []string                                                        `json:"pinnedTweetIds"`
	IsAutomated                bool                                                            `json:"isAutomated"`
	AutomatedBy                string                                                          `json:"automatedBy"`
	Unavailable                bool                                                            `json:"unavailable"`
	Message                    string                                                          `json:"message"`
	UnavailableReason          string                                                          `json:"unavailableReason"`
	ProfileBio                 *GetCommunityByIDCommunityInfoCreatorProfileBio                 `json:"profile_bio"`
}

type GetCommunityByIDCommunityInfoAdminAffiliatesHighlightedLabel map[string]any

type GetCommunityByIDCommunityInfoAdminProfileBioEntitiesDescriptionURL struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
	URL         string `json:"url"`
}

type GetCommunityByIDCommunityInfoAdminProfileBioEntitiesDescription struct {
	URLs []*GetCommunityByIDCommunityInfoAdminProfileBioEntitiesDescriptionURL `json:"urls"`
}

type GetCommunityByIDCommunityInfoAdminProfileBioEntitiesURLURL struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
	URL         string `json:"url"`
}

type GetCommunityByIDCommunityInfoAdminProfileBioEntitiesURL struct {
	URLs []*GetCommunityByIDCommunityInfoAdminProfileBioEntitiesURLURL `json:"urls"`
}

type GetCommunityByIDCommunityInfoAdminProfileBioEntities struct {
	Description *GetCommunityByIDCommunityInfoAdminProfileBioEntitiesDescription `json:"description"`
	URL         *GetCommunityByIDCommunityInfoAdminProfileBioEntitiesURL         `json:"url"`
}

type GetCommunityByIDCommunityInfoAdminProfileBio struct {
	Description string                                                `json:"description"`
	Entities    *GetCommunityByIDCommunityInfoAdminProfileBioEntities `json:"entities"`
}

type GetCommunityByIDCommunityInfoAdmin struct {
	Type                       string                                                        `json:"type"`
	UserName                   string                                                        `json:"userName"`
	URL                        string                                                        `json:"url"`
	ID                         string                                                        `json:"id"`
	Name                       string                                                        `json:"name"`
	IsBlueVerified             bool                                                          `json:"isBlueVerified"`
	VerifiedType               string                                                        `json:"verifiedType"`
	ProfilePicture             string                                                        `json:"profilePicture"`
	CoverPicture               string                                                        `json:"coverPicture"`
	Description                string                                                        `json:"description"`
	Location                   string                                                        `json:"location"`
	Followers                  int                                                           `json:"followers"`
	Following                  int                                                           `json:"following"`
	CanDM                      bool                                                          `json:"canDm"`
	CreatedAt                  string                                                        `json:"createdAt"`
	FavouritesCount            int                                                           `json:"favouritesCount"`
	HasCustomTimelines         bool                                                          `json:"hasCustomTimelines"`
	IsTranslator               bool                                                          `json:"isTranslator"`
	MediaCount                 int                                                           `json:"mediaCount"`
	StatusesCount              int                                                           `json:"statusesCount"`
	WithheldInCountries        []string                                                      `json:"withheldInCountries"`
	AffiliatesHighlightedLabel *GetCommunityByIDCommunityInfoAdminAffiliatesHighlightedLabel `json:"affiliatesHighlightedLabel"`
	PossiblySensitive          bool                                                          `json:"possiblySensitive"`
	PinnedTweetIDs             []string                                                      `json:"pinnedTweetIds"`
	IsAutomated                bool                                                          `json:"isAutomated"`
	AutomatedBy                string                                                        `json:"automatedBy"`
	Unavailable                bool                                                          `json:"unavailable"`
	Message                    string                                                        `json:"message"`
	UnavailableReason          string                                                        `json:"unavailableReason"`
	ProfileBio                 *GetCommunityByIDCommunityInfoAdminProfileBio                 `json:"profile_bio"`
}

type GetCommunityByIDCommunityInfoMembersPreview struct {
}

type GetCommunityByIDCommunityInfo struct {
	ID             string                                         `json:"id"`
	Name           string                                         `json:"name"`
	Description    string                                         `json:"description"`
	Question       string                                         `json:"question"`
	MemberCount    int                                            `json:"member_count"`
	ModeratorCount int                                            `json:"moderator_count"`
	CreatedAt      string                                         `json:"created_at"`
	JoinPolicy     string                                         `json:"join_policy"`
	InvitesPolicy  string                                         `json:"invites_policy"`
	IsNsfw         bool                                           `json:"is_nsfw"`
	IsPinned       bool                                           `json:"is_pinned"`
	Role           string                                         `json:"role"`
	PrimaryTopic   *GetCommunityByIDCommunityInfoPrimaryTopic     `json:"primary_topic"`
	BannerURL      string                                         `json:"banner_url"`
	SearchTags     []string                                       `json:"search_tags"`
	Rules          []*GetCommunityByIDCommunityInfoRule           `json:"rules"`
	Creator        *GetCommunityByIDCommunityInfoCreator          `json:"creator"`
	Admin          *GetCommunityByIDCommunityInfoAdmin            `json:"admin"`
	MembersPreview []*GetCommunityByIDCommunityInfoMembersPreview `json:"members_preview"`
}

type GetCommunityByIDResponse struct {
	CommunityInfo *GetCommunityByIDCommunityInfo `json:"community_info"`
	Status        string                         `json:"status"`
	Message       string                         `json:"msg"`
}

func (t *twitterApi) GetCommunityByID(communityID string) (*GetCommunityByIDResponse, error) {
	if strings.TrimSpace(communityID) == "" {
		return nil, errors.New("communityID is required")
	}

	url := twitterDomainURI + "/community/info?" + strings.Join([]string{"community_id=" + communityID}, "&")

	ctx1, cancel1 := context.WithTimeout(t.ctx, time.Second*10)
	defer cancel1()

	jsonData, resp, err := getDataWithHeader(ctx1, t.httpClient, url, t.headers)
	if err != nil {
		slog.Error("GetCommunityByID failed", "err", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		slog.Error("GetCommunityByID failed", "statusCode", resp.StatusCode, "body", string(jsonData))
		return nil, errors.New("GetCommunityByID failed")
	}

	response := &GetCommunityByIDResponse{}
	if err = jsoniter.Unmarshal(jsonData, &response); err != nil {
		slog.Error("GetCommunityByID failed", "err", err)
		return nil, err
	}
	if response.Status != "success" {
		slog.Error("GetCommunityByID failed", "status", response.Status, "message", response.Message)
		return nil, errors.New("GetCommunityByID failed")
	}

	return response, nil
}
