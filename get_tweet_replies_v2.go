// Doc https://docs.twitterapi.io/api-reference/endpoint/get_tweet_replies_v2

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

type GetTweetRepliesV2ReplyAuthorAffiliatesHighlightedLabel map[string]any

type GetTweetRepliesV2ReplyAuthorProfileBioEntitiesDescriptionURL struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
	URL         string `json:"url"`
}

type GetTweetRepliesV2ReplyAuthorProfileBioEntitiesDescription struct {
	URLs []*GetTweetRepliesV2ReplyAuthorProfileBioEntitiesDescriptionURL `json:"urls"`
}

type GetTweetRepliesV2ReplyAuthorProfileBioEntitiesURLURL struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
	URL         string `json:"url"`
}

type GetTweetRepliesV2ReplyAuthorProfileBioEntitiesURL struct {
	URLs []*GetTweetRepliesV2ReplyAuthorProfileBioEntitiesURLURL `json:"urls"`
}

type GetTweetRepliesV2ReplyAuthorProfileBioEntities struct {
	Description *GetTweetRepliesV2ReplyAuthorProfileBioEntitiesDescription `json:"description"`
	URL         *GetTweetRepliesV2ReplyAuthorProfileBioEntitiesURL         `json:"url"`
}

type GetTweetRepliesV2ReplyAuthorProfileBio struct {
	Description string                                          `json:"description"`
	Entities    *GetTweetRepliesV2ReplyAuthorProfileBioEntities `json:"entities"`
}

type GetTweetRepliesV2ReplyAuthor struct {
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
	AffiliatesHighlightedLabel *GetTweetRepliesV2ReplyAuthorAffiliatesHighlightedLabel `json:"affiliatesHighlightedLabel"`
	PossiblySensitive          bool                                                    `json:"possiblySensitive"`
	PinnedTweetIDs             []string                                                `json:"pinnedTweetIds"`
	IsAutomated                bool                                                    `json:"isAutomated"`
	AutomatedBy                string                                                  `json:"automatedBy"`
	Unavailable                bool                                                    `json:"unavailable"`
	Message                    string                                                  `json:"message"`
	UnavailableReason          string                                                  `json:"unavailableReason"`
	ProfileBio                 *GetTweetRepliesV2ReplyAuthorProfileBio                 `json:"profile_bio"`
}

type GetTweetRepliesV2ReplyEntitiesHashtag struct {
	Indices []int  `json:"indices"`
	Text    string `json:"text"`
}

type GetTweetRepliesV2ReplyEntitiesURL struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
	URL         string `json:"url"`
}

type GetTweetRepliesV2ReplyEntitiesUserMention struct {
	IDStr      string `json:"id_str"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
}

type GetTweetRepliesV2ReplyEntities struct {
	Hashtags     []*GetTweetRepliesV2ReplyEntitiesHashtag     `json:"hashtags"`
	URLs         []*GetTweetRepliesV2ReplyEntitiesURL         `json:"urls"`
	UserMentions []*GetTweetRepliesV2ReplyEntitiesUserMention `json:"user_mentions"`
}

type GetTweetRepliesV2Reply struct {
	Type              string                          `json:"type"`
	ID                string                          `json:"id"`
	URL               string                          `json:"url"`
	TwitterURL        string                          `json:"twitterUrl"` // added: missing from original
	Text              string                          `json:"text"`
	Source            string                          `json:"source"`
	RetweetCount      int                             `json:"retweetCount"`
	ReplyCount        int                             `json:"replyCount"`
	LikeCount         int                             `json:"likeCount"`
	QuoteCount        int                             `json:"quoteCount"`
	ViewCount         int                             `json:"viewCount"`
	CreatedAt         string                          `json:"createdAt"`
	Lang              string                          `json:"lang"`
	BookmarkCount     int                             `json:"bookmarkCount"`
	IsReply           bool                            `json:"isReply"`
	InReplyToID       string                          `json:"inReplyToId"`
	ConversationID    string                          `json:"conversationId"`
	DisplayTextRange  []int                           `json:"displayTextRange"`
	InReplyToUserID   string                          `json:"inReplyToUserId"`
	InReplyToUsername string                          `json:"inReplyToUsername"`
	Author            *GetTweetRepliesV2ReplyAuthor   `json:"author"`
	Entities          *GetTweetRepliesV2ReplyEntities `json:"entities"`
	QuotedTweet       *GetTweetRepliesV2Reply         `json:"quoted_tweet"`
	RetweetedTweet    *GetTweetRepliesV2Reply         `json:"retweeted_tweet"`
	IsLimitedReply    bool                            `json:"isLimitedReply"`
}

type GetTweetRepliesV2Response struct {
	Tweets      []*GetTweetRepliesV2Reply `json:"tweets"` // fixed: was replies
	HasNextPage bool                      `json:"has_next_page"`
	NextCursor  string                    `json:"next_cursor"`
	Status      string                    `json:"status"`
	Message     string                    `json:"msg"` // fixed: was message
}

func (t *TwitterApi) GetTweetRepliesV2(tweetID string, cursor *string) (*GetTweetRepliesV2Response, error) {
	if strings.TrimSpace(tweetID) == "" {
		return nil, errors.New("tweetID is required")
	}

	vals := neturl.Values{}
	vals.Set("tweetId", tweetID)
	if cursor != nil && *cursor != "" {
		vals.Set("cursor", *cursor)
	}
	url := twitterDomainURI + "/tweet/replies/v2?" + vals.Encode()

	ctx1, cancel1 := context.WithTimeout(t.ctx, time.Second*10)
	defer cancel1()

	jsonData, resp, err := getDataWithHeader(ctx1, t.httpClient, url, t.headers)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			slog.Error("GetTweetRepliesV2 request timed out", "url", url)
			return nil, errors.New("GetTweetRepliesV2 request timed out")
		}
		slog.Error("GetTweetRepliesV2 failed", "err", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		slog.Error("GetTweetRepliesV2 failed", "statusCode", resp.StatusCode, "body", string(jsonData))
		return nil, errors.New("GetTweetRepliesV2 failed")
	}

	response := &GetTweetRepliesV2Response{}
	if err = jsoniter.Unmarshal(jsonData, &response); err != nil {
		slog.Error("GetTweetRepliesV2 failed", "err", err)
		return nil, err
	}
	if response.Status != "success" {
		slog.Error("GetTweetRepliesV2 failed", "status", response.Status, "message", response.Message)
		return nil, errors.New(response.Message)
	}

	return response, nil
}
