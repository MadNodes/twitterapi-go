// Doc https://docs.twitterapi.io/api-reference/endpoint/get_article

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

type GetArticleArticleAuthorAffiliatesHighlightedLabel map[any]any

type GetArticleArticleAuthorProfileBioEntitiesDescriptionURL struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
	URL         string `json:"url"`
}

type GetArticleArticleAuthorProfileBioEntitiesDescription struct {
	URLs []*GetArticleArticleAuthorProfileBioEntitiesDescriptionURL `json:"urls"`
}

type GetArticleArticleAuthorProfileBioEntitiesURLURL struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
	URL         string `json:"url"`
}

type GetArticleArticleAuthorProfileBioEntitiesURL struct {
	URLs []*GetArticleArticleAuthorProfileBioEntitiesURLURL `json:"urls"`
}

type GetArticleArticleAuthorProfileBioEntities struct {
	Description *GetArticleArticleAuthorProfileBioEntitiesDescription `json:"description"`
	URL         *GetArticleArticleAuthorProfileBioEntitiesURL         `json:"url"`
}

type GetArticleArticleAuthorProfileBio struct {
	Description string                                     `json:"description"`
	Entities    *GetArticleArticleAuthorProfileBioEntities `json:"entities"`
}

type GetArticleArticleAuthor struct {
	Type                       string                                             `json:"type"`
	UserName                   string                                             `json:"userName"`
	URL                        string                                             `json:"url"`
	ID                         string                                             `json:"id"`
	Name                       string                                             `json:"name"`
	IsBlueVerified             bool                                               `json:"isBlueVerified"`
	VerifiedType               string                                             `json:"verifiedType"`
	ProfilePicture             string                                             `json:"profilePicture"`
	CoverPicture               string                                             `json:"coverPicture"`
	Description                string                                             `json:"description"`
	Location                   string                                             `json:"location"`
	Followers                  int                                                `json:"followers"`
	Following                  int                                                `json:"following"`
	CanDM                      bool                                               `json:"canDm"`
	CreatedAt                  string                                             `json:"createdAt"`
	FavouritesCount            int                                                `json:"favouritesCount"`
	HasCustomTimelines         bool                                               `json:"hasCustomTimelines"`
	IsTranslator               bool                                               `json:"isTranslator"`
	MediaCount                 int                                                `json:"mediaCount"`
	StatusesCount              int                                                `json:"statusesCount"`
	WithheldInCountries        []string                                           `json:"withheldInCountries"`
	AffiliatesHighlightedLabel *GetArticleArticleAuthorAffiliatesHighlightedLabel `json:"affiliatesHighlightedLabel"`
	PossiblySensitive          bool                                               `json:"possiblySensitive"`
	PinnedTweetIDs             []string                                           `json:"pinnedTweetIds"`
	IsAutomated                bool                                               `json:"isAutomated"`
	AutomatedBy                string                                             `json:"automatedBy"`
	Unavailable                bool                                               `json:"unavailable"`
	Message                    string                                             `json:"message"`
	UnavailableReason          string                                             `json:"unavailableReason"`
	ProfileBio                 *GetArticleArticleAuthorProfileBio                 `json:"profile_bio"`
}

type GetArticleArticleContentInlineStyleRange struct {
	Offset int    `json:"offset"`
	Length int    `json:"length"`
	Style  string `json:"style"`
}

type GetArticleArticleContent struct {
	Type              string                                      `json:"type"`
	Text              string                                      `json:"text"`
	URL               string                                      `json:"url"`
	PreviewURL        string                                      `json:"previewUrl"`
	Width             int                                         `json:"width"`
	Height            int                                         `json:"height"`
	InlineStyleRanges []*GetArticleArticleContentInlineStyleRange `json:"inlineStyleRanges"`
}

type GetArticleArticle struct {
	Author           *GetArticleArticleAuthor    `json:"author"`
	ReplyCount       int                         `json:"replyCount"`
	LikeCount        int                         `json:"likeCount"`
	QuoteCount       int                         `json:"quoteCount"`
	ViewCount        int                         `json:"viewCount"`
	CreatedAt        string                      `json:"createdAt"`
	Title            string                      `json:"title"`
	PreviewText      string                      `json:"preview_text"`
	CoverMediaImgURL string                      `json:"cover_media_img_url"`
	Contents         []*GetArticleArticleContent `json:"contents"`
}

type GetArticleResponse struct {
	Article *GetArticleArticle `json:"article"`
	Status  string             `json:"status"`
	Message string             `json:"message"`
}

func (t *twitterApi) GetArticle(tweetID string) (*GetArticleResponse, error) {
	if tweetID == "" {
		return nil, errors.New("tweet_id is empty")
	}

	url := twitterDomainURI + "/article?" + strings.Join([]string{"tweet_id=" + tweetID}, "&")

	ctx1, cancel1 := context.WithTimeout(t.ctx, time.Second*10)
	defer cancel1()

	jsonData, resp, err := getDataWithHeader(ctx1, t.httpClient, url, t.headers)
	if err != nil {
		slog.Error("GetArticle failed", "err", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		slog.Error("GetArticle failed", "statusCode", resp.StatusCode, "body", string(jsonData))
		return nil, errors.New("GetArticle failed")
	}

	response := &GetArticleResponse{}
	if err = jsoniter.Unmarshal(jsonData, &response); err != nil {
		slog.Error("GetArticle failed", "err", err)
		return nil, err
	}
	if response.Status != "success" {
		slog.Error("GetArticle failed", "status", response.Status, "message", response.Message)
		return nil, errors.New("GetArticle failed")
	}

	return response, nil
}
