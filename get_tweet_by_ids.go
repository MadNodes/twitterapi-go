// Doc https://docs.twitterapi.io/api-reference/endpoint/get_tweet_by_ids

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

type GetTweetByIDsTweetAuthorAffiliatesHighlightedLabel map[any]any

type GetTweetByIDsTweetAuthorProfileBioEntitiesDescriptionURL struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
	URL         string `json:"url"`
}

type GetTweetByIDsTweetAuthorProfileBioEntitiesDescription struct {
	URLs []*GetTweetByIDsTweetAuthorProfileBioEntitiesDescriptionURL `json:"urls"`
}

type GetTweetByIDsTweetAuthorProfileBioEntitiesURLURL struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
	URL         string `json:"url"`
}

type GetTweetByIDsTweetAuthorProfileBioEntitiesURL struct {
	URLs []*GetTweetByIDsTweetAuthorProfileBioEntitiesURLURL `json:"urls"`
}

type GetTweetByIDsTweetAuthorProfileBioEntities struct {
	Description *GetTweetByIDsTweetAuthorProfileBioEntitiesDescription `json:"description"`
	URL         *GetTweetByIDsTweetAuthorProfileBioEntitiesURL         `json:"url"`
}

type GetTweetByIDsTweetAuthorProfileBio struct {
	Description string                                      `json:"description"`
	Entities    *GetTweetByIDsTweetAuthorProfileBioEntities `json:"entities"`
}

type GetTweetByIDsTweetAuthor struct {
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
	AffiliatesHighlightedLabel *GetTweetByIDsTweetAuthorAffiliatesHighlightedLabel `json:"affiliatesHighlightedLabel"`
	PossiblySensitive          bool                                                `json:"possiblySensitive"`
	PinnedTweetIDs             []string                                            `json:"pinnedTweetIds"`
	IsAutomated                bool                                                `json:"isAutomated"`
	AutomatedBy                string                                              `json:"automatedBy"`
	Unavailable                bool                                                `json:"unavailable"`
	Message                    string                                              `json:"message"`
	UnavailableReason          string                                              `json:"unavailableReason"`
	ProfileBio                 *GetTweetByIDsTweetAuthorProfileBio                 `json:"profile_bio"`
}

type GetTweetByIDsTweetEntitiesHashtag struct {
	Indices []int  `json:"indices"`
	Text    string `json:"text"`
}

type GetTweetByIDsTweetEntitiesURL struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
	URL         string `json:"url"`
}

type GetTweetByIDsTweetEntitiesUserMention struct {
	IDStr      string `json:"id_str"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
}

type GetTweetByIDsTweetEntities struct {
	Hashtags     []*GetTweetByIDsTweetEntitiesHashtag     `json:"hashtags"`
	URLs         []*GetTweetByIDsTweetEntitiesURL         `json:"urls"`
	UserMentions []*GetTweetByIDsTweetEntitiesUserMention `json:"user_mentions"`
}

type GetTweetByIDsTweet struct {
	Type              string                      `json:"type"`
	ID                string                      `json:"id"`
	URL               string                      `json:"url"`
	Text              string                      `json:"text"`
	Source            string                      `json:"source"`
	RetweetCount      int                         `json:"retweetCount"`
	ReplyCount        int                         `json:"replyCount"`
	LikeCount         int                         `json:"likeCount"`
	QuoteCount        int                         `json:"quoteCount"`
	ViewCount         int                         `json:"viewCount"`
	CreatedAt         string                      `json:"createdAt"`
	Lang              string                      `json:"lang"`
	BookmarkCount     int                         `json:"bookmarkCount"`
	IsReply           bool                        `json:"isReply"`
	InReplyToID       string                      `json:"inReplyToId"`
	ConversationID    string                      `json:"conversationId"`
	DisplayTextRange  []int                       `json:"displayTextRange"`
	InReplyToUserID   string                      `json:"inReplyToUserId"`
	InReplyToUsername string                      `json:"inReplyToUsername"`
	Author            *GetTweetByIDsTweetAuthor   `json:"author"`
	Entities          *GetTweetByIDsTweetEntities `json:"entities"`
	QuotedTweet       *GetTweetByIDsTweet         `json:"quoted_tweet"`
	RetweetedTweet    *GetTweetByIDsTweet         `json:"retweeted_tweet"`
	IsLimitedReply    bool                        `json:"isLimitedReply"`
}

type GetTweetByIDsResponse struct {
	Tweets  []*GetTweetByIDsTweet `json:"tweets"`
	Status  string                `json:"status"`
	Message string                `json:"message"`
}

func (t *twitterApi) GetTweetByIDs(tweetIDs []string) (*GetTweetByIDsResponse, error) {
	if len(tweetIDs) == 0 {
		return nil, errors.New("tweet_ids is empty")
	}

	url := tweetsTwitterDomainURI + "?" + strings.Join([]string{"tweet_ids=" + strings.Join(tweetIDs, ",")}, "&")

	ctx1, cancel1 := context.WithTimeout(t.ctx, time.Second*10)
	defer cancel1()

	jsonData, resp, err := getDataWithHeader(ctx1, t.httpClient, url, t.headers)
	if err != nil {
		slog.Error("GetTweetByIDs failed", "err", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		slog.Error("GetTweetByIDs failed", "statusCode", resp.StatusCode, "body", string(jsonData))
		return nil, errors.New("GetTweetByIDs failed")
	}

	response := &GetTweetByIDsResponse{}
	if err = jsoniter.Unmarshal(jsonData, &response); err != nil {
		slog.Error("GetTweetByIDs failed", "err", err)
		return nil, err
	}
	if response.Status != "success" {
		slog.Error("GetTweetByIDs failed", "status", response.Status, "message", response.Message)
		return nil, errors.New("GetTweetByIDs failed")
	}

	return response, nil
}
