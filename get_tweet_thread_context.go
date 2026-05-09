// Doc https://docs.twitterapi.io/api-reference/endpoint/get_tweet_thread_context

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

type GetTweetThreadContextReplyAuthorAffiliatesHighlightedLabel map[string]any

type GetTweetThreadContextReplyAuthorProfileBioEntitiesDescriptionURL struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
	URL         string `json:"url"`
}

type GetTweetThreadContextReplyAuthorProfileBioEntitiesDescription struct {
	URLs []*GetTweetThreadContextReplyAuthorProfileBioEntitiesDescriptionURL `json:"urls"`
}

type GetTweetThreadContextReplyAuthorProfileBioEntitiesURLURL struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
	URL         string `json:"url"`
}

type GetTweetThreadContextReplyAuthorProfileBioEntitiesURL struct {
	URLs []*GetTweetThreadContextReplyAuthorProfileBioEntitiesURLURL `json:"urls"`
}

type GetTweetThreadContextReplyAuthorProfileBioEntities struct {
	Description *GetTweetThreadContextReplyAuthorProfileBioEntitiesDescription `json:"description"`
	URL         *GetTweetThreadContextReplyAuthorProfileBioEntitiesURL         `json:"url"`
}

type GetTweetThreadContextReplyAuthorProfileBio struct {
	Description string                                              `json:"description"`
	Entities    *GetTweetThreadContextReplyAuthorProfileBioEntities `json:"entities"`
}

type GetTweetThreadContextReplyAuthor struct {
	Type                       string                                                      `json:"type"`
	UserName                   string                                                      `json:"userName"`
	URL                        string                                                      `json:"url"`
	ID                         string                                                      `json:"id"`
	Name                       string                                                      `json:"name"`
	IsBlueVerified             bool                                                        `json:"isBlueVerified"`
	VerifiedType               string                                                      `json:"verifiedType"`
	ProfilePicture             string                                                      `json:"profilePicture"`
	CoverPicture               string                                                      `json:"coverPicture"`
	Description                string                                                      `json:"description"`
	Location                   string                                                      `json:"location"`
	Followers                  int                                                         `json:"followers"`
	Following                  int                                                         `json:"following"`
	CanDM                      bool                                                        `json:"canDm"`
	CreatedAt                  string                                                      `json:"createdAt"`
	FavouritesCount            int                                                         `json:"favouritesCount"`
	HasCustomTimelines         bool                                                        `json:"hasCustomTimelines"`
	IsTranslator               bool                                                        `json:"isTranslator"`
	MediaCount                 int                                                         `json:"mediaCount"`
	StatusesCount              int                                                         `json:"statusesCount"`
	WithheldInCountries        []string                                                    `json:"withheldInCountries"`
	AffiliatesHighlightedLabel *GetTweetThreadContextReplyAuthorAffiliatesHighlightedLabel `json:"affiliatesHighlightedLabel"`
	PossiblySensitive          bool                                                        `json:"possiblySensitive"`
	PinnedTweetIDs             []string                                                    `json:"pinnedTweetIds"`
	IsAutomated                bool                                                        `json:"isAutomated"`
	AutomatedBy                string                                                      `json:"automatedBy"`
	Unavailable                bool                                                        `json:"unavailable"`
	Message                    string                                                      `json:"message"`
	UnavailableReason          string                                                      `json:"unavailableReason"`
	ProfileBio                 *GetTweetThreadContextReplyAuthorProfileBio                 `json:"profile_bio"`
}

type GetTweetThreadContextReplyEntitiesHashtag struct {
	Indices []int  `json:"indices"`
	Text    string `json:"text"`
}

type GetTweetThreadContextReplyEntitiesURL struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
	URL         string `json:"url"`
}

type GetTweetThreadContextReplyEntitiesUserMention struct {
	IDStr      string `json:"id_str"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
}

type GetTweetThreadContextReplyEntities struct {
	Hashtags     []*GetTweetThreadContextReplyEntitiesHashtag     `json:"hashtags"`
	URLs         []*GetTweetThreadContextReplyEntitiesURL         `json:"urls"`
	UserMentions []*GetTweetThreadContextReplyEntitiesUserMention `json:"user_mentions"`
}

type GetTweetThreadContextReply struct {
	Type              string                              `json:"type"`
	ID                string                              `json:"id"`
	URL               string                              `json:"url"`
	TwitterURL        string                              `json:"twitterUrl"` // added: missing from original
	Text              string                              `json:"text"`
	Source            string                              `json:"source"`
	RetweetCount      int                                 `json:"retweetCount"`
	ReplyCount        int                                 `json:"replyCount"`
	LikeCount         int                                 `json:"likeCount"`
	QuoteCount        int                                 `json:"quoteCount"`
	ViewCount         int                                 `json:"viewCount"`
	CreatedAt         string                              `json:"createdAt"`
	Lang              string                              `json:"lang"`
	BookmarkCount     int                                 `json:"bookmarkCount"`
	IsReply           bool                                `json:"isReply"`
	InReplyToID       string                              `json:"inReplyToId"`
	ConversationID    string                              `json:"conversationId"`
	DisplayTextRange  []int                               `json:"displayTextRange"`
	InReplyToUserID   string                              `json:"inReplyToUserId"`
	InReplyToUsername string                              `json:"inReplyToUsername"`
	Author            *GetTweetThreadContextReplyAuthor   `json:"author"`
	Entities          *GetTweetThreadContextReplyEntities `json:"entities"`
	QuotedTweet       *GetTweetThreadContextReply         `json:"quoted_tweet"`
	RetweetedTweet    *GetTweetThreadContextReply         `json:"retweeted_tweet"`
	IsLimitedReply    bool                                `json:"isLimitedReply"`
}

type GetTweetThreadContextResponse struct {
	Tweets      []*GetTweetThreadContextReply `json:"tweets"` // fixed: was replies
	HasNextPage bool                          `json:"has_next_page"`
	NextCursor  string                        `json:"next_cursor"`
	Status      string                        `json:"status"`
	Message     string                        `json:"msg"` // fixed: was message
}

func (t *twitterApi) GetTweetThreadContext(tweetID string, cursor *string) (*GetTweetThreadContextResponse, error) {
	if strings.TrimSpace(tweetID) == "" {
		return nil, errors.New("tweetID is required")
	}

	vals := neturl.Values{}
	vals.Set("tweetId", tweetID)
	if cursor != nil && *cursor != "" {
		vals.Set("cursor", *cursor)
	}
	url := twitterDomainURI + "/tweet/thread_context?" + vals.Encode()

	ctx1, cancel1 := context.WithTimeout(t.ctx, time.Second*10)
	defer cancel1()

	jsonData, resp, err := getDataWithHeader(ctx1, t.httpClient, url, t.headers)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			slog.Error("GetTweetThreadContext request timed out", "url", url)
			return nil, errors.New("GetTweetThreadContext request timed out")
		}
		slog.Error("GetTweetThreadContext failed", "err", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		slog.Error("GetTweetThreadContext failed", "statusCode", resp.StatusCode, "body", string(jsonData))
		return nil, errors.New("GetTweetThreadContext failed")
	}

	response := &GetTweetThreadContextResponse{}
	if err = jsoniter.Unmarshal(jsonData, &response); err != nil {
		slog.Error("GetTweetThreadContext failed", "err", err)
		return nil, err
	}
	if response.Status != "success" {
		slog.Error("GetTweetThreadContext failed", "status", response.Status, "message", response.Message)
		return nil, errors.New(response.Message)
	}

	return response, nil
}
