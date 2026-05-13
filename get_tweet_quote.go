// Doc https://docs.twitterapi.io/api-reference/endpoint/get_tweet_quote

package twitterapi

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	neturl "net/url"
	"strconv"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type GetTweetQuoteTweetAuthorAffiliatesHighlightedLabel map[string]any

type GetTweetQuoteTweetAuthorProfileBioEntitiesDescriptionURL struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
	URL         string `json:"url"`
}

type GetTweetQuoteTweetAuthorProfileBioEntitiesDescription struct {
	URLs []*GetTweetQuoteTweetAuthorProfileBioEntitiesDescriptionURL `json:"urls"`
}

type GetTweetQuoteTweetAuthorProfileBioEntitiesURLURL struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
	URL         string `json:"url"`
}

type GetTweetQuoteTweetAuthorProfileBioEntitiesURL struct {
	URLs []*GetTweetQuoteTweetAuthorProfileBioEntitiesURLURL `json:"urls"`
}

type GetTweetQuoteTweetAuthorProfileBioEntities struct {
	Description *GetTweetQuoteTweetAuthorProfileBioEntitiesDescription `json:"description"`
	URL         *GetTweetQuoteTweetAuthorProfileBioEntitiesURL         `json:"url"`
}

type GetTweetQuoteTweetAuthorProfileBio struct {
	Description string                                      `json:"description"`
	Entities    *GetTweetQuoteTweetAuthorProfileBioEntities `json:"entities"`
}

type GetTweetQuoteTweetAuthor struct {
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
	AffiliatesHighlightedLabel *GetTweetQuoteTweetAuthorAffiliatesHighlightedLabel `json:"affiliatesHighlightedLabel"`
	PossiblySensitive          bool                                                `json:"possiblySensitive"`
	PinnedTweetIDs             []string                                            `json:"pinnedTweetIds"`
	IsAutomated                bool                                                `json:"isAutomated"`
	AutomatedBy                string                                              `json:"automatedBy"`
	Unavailable                bool                                                `json:"unavailable"`
	Message                    string                                              `json:"message"`
	UnavailableReason          string                                              `json:"unavailableReason"`
	ProfileBio                 *GetTweetQuoteTweetAuthorProfileBio                 `json:"profile_bio"`
}

type GetTweetQuoteTweetEntitiesHashtag struct {
	Indices []int  `json:"indices"`
	Text    string `json:"text"`
}

type GetTweetQuoteTweetEntitiesURL struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
	URL         string `json:"url"`
}

type GetTweetQuoteTweetEntitiesUserMention struct {
	IDStr      string `json:"id_str"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
}

type GetTweetQuoteTweetEntities struct {
	Hashtags     []*GetTweetQuoteTweetEntitiesHashtag     `json:"hashtags"`
	URLs         []*GetTweetQuoteTweetEntitiesURL         `json:"urls"`
	UserMentions []*GetTweetQuoteTweetEntitiesUserMention `json:"user_mentions"`
}

type GetTweetQuoteTweet struct {
	Type              string                      `json:"type"`
	ID                string                      `json:"id"`
	URL               string                      `json:"url"`
	TwitterURL        string                      `json:"twitterUrl"` // added: missing from original
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
	Author            *GetTweetQuoteTweetAuthor   `json:"author"`
	Entities          *GetTweetQuoteTweetEntities `json:"entities"`
	QuotedTweet       *GetTweetQuoteTweet         `json:"quoted_tweet"`
	RetweetedTweet    *GetTweetQuoteTweet         `json:"retweeted_tweet"`
	IsLimitedReply    bool                        `json:"isLimitedReply"`
}

type GetTweetQuoteResponse struct {
	Tweets      []*GetTweetQuoteTweet `json:"tweets"`
	HasNextPage bool                  `json:"has_next_page"`
	NextCursor  string                `json:"next_cursor"`
	Status      string                `json:"status"`
	Message     string                `json:"msg"` // fixed: was message
}

func (t *TwitterApi) GetTweetQuote(tweetID string, sinceTime *int, untilTime *int, includeReplies *bool, cursor *string) (*GetTweetQuoteResponse, error) {
	if strings.TrimSpace(tweetID) == "" {
		return nil, errors.New("tweetID is required")
	}

	vals := neturl.Values{}
	vals.Set("tweetId", tweetID)
	if sinceTime != nil {
		vals.Set("sinceTime", strconv.Itoa(*sinceTime))
	}
	if untilTime != nil {
		vals.Set("untilTime", strconv.Itoa(*untilTime))
	}
	if includeReplies != nil {
		vals.Set("includeReplies", strconv.FormatBool(*includeReplies))
	}
	if cursor != nil && *cursor != "" {
		vals.Set("cursor", *cursor)
	}
	url := twitterDomainURI + "/tweet/quotes?" + vals.Encode()

	ctx1, cancel1 := context.WithTimeout(t.ctx, time.Second*10)
	defer cancel1()

	jsonData, resp, err := getDataWithHeader(ctx1, t.httpClient, url, t.headers)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			slog.Error("GetTweetQuote request timed out", "url", url)
			return nil, errors.New("GetTweetQuote request timed out")
		}
		slog.Error("GetTweetQuote failed", "err", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		slog.Error("GetTweetQuote failed", "statusCode", resp.StatusCode, "body", string(jsonData))
		return nil, errors.New("GetTweetQuote failed")
	}

	response := &GetTweetQuoteResponse{}
	if err = jsoniter.Unmarshal(jsonData, &response); err != nil {
		slog.Error("GetTweetQuote failed", "err", err)
		return nil, err
	}
	if response.Status != "success" {
		slog.Error("GetTweetQuote failed", "status", response.Status, "message", response.Message)
		return nil, errors.New(response.Message)
	}

	return response, nil
}
