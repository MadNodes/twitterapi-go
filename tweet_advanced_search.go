// Doc https://docs.twitterapi.io/api-reference/endpoint/tweet_advanced_search

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

type TweetAdvancedSearchTweetAuthorAffiliatesHighlightedLabel map[any]any

type TweetAdvancedSearchTweetAuthorProfileBioEntitiesDescriptionURL struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
	URL         string `json:"url"`
}

type TweetAdvancedSearchTweetAuthorProfileBioEntitiesDescription struct {
	URLs []*TweetAdvancedSearchTweetAuthorProfileBioEntitiesDescriptionURL `json:"urls"`
}

type TweetAdvancedSearchTweetAuthorProfileBioEntitiesURLURL struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
	URL         string `json:"url"`
}

type TweetAdvancedSearchTweetAuthorProfileBioEntitiesURL struct {
	URLs []*TweetAdvancedSearchTweetAuthorProfileBioEntitiesURLURL `json:"urls"`
}

type TweetAdvancedSearchTweetAuthorProfileBioEntities struct {
	Description *TweetAdvancedSearchTweetAuthorProfileBioEntitiesDescription `json:"description"`
	URL         *TweetAdvancedSearchTweetAuthorProfileBioEntitiesURL         `json:"url"`
}

type TweetAdvancedSearchTweetAuthorProfileBio struct {
	Description string                                            `json:"description"`
	Entities    *TweetAdvancedSearchTweetAuthorProfileBioEntities `json:"entities"`
}

type TweetAdvancedSearchTweetAuthor struct {
	Type                       string                                                    `json:"type"`
	UserName                   string                                                    `json:"userName"`
	URL                        string                                                    `json:"url"`
	ID                         string                                                    `json:"id"`
	Name                       string                                                    `json:"name"`
	IsBlueVerified             bool                                                      `json:"isBlueVerified"`
	VerifiedType               string                                                    `json:"verifiedType"`
	ProfilePicture             string                                                    `json:"profilePicture"`
	CoverPicture               string                                                    `json:"coverPicture"`
	Description                string                                                    `json:"description"`
	Location                   string                                                    `json:"location"`
	Followers                  int                                                       `json:"followers"`
	Following                  int                                                       `json:"following"`
	CanDM                      bool                                                      `json:"canDm"`
	CreatedAt                  string                                                    `json:"createdAt"`
	FavouritesCount            int                                                       `json:"favouritesCount"`
	HasCustomTimelines         bool                                                      `json:"hasCustomTimelines"`
	IsTranslator               bool                                                      `json:"isTranslator"`
	MediaCount                 int                                                       `json:"mediaCount"`
	StatusesCount              int                                                       `json:"statusesCount"`
	WithheldInCountries        []string                                                  `json:"withheldInCountries"`
	AffiliatesHighlightedLabel *TweetAdvancedSearchTweetAuthorAffiliatesHighlightedLabel `json:"affiliatesHighlightedLabel"`
	PossiblySensitive          bool                                                      `json:"possiblySensitive"`
	PinnedTweetIDs             []string                                                  `json:"pinnedTweetIds"`
	IsAutomated                bool                                                      `json:"isAutomated"`
	AutomatedBy                string                                                    `json:"automatedBy"`
	Unavailable                bool                                                      `json:"unavailable"`
	Message                    string                                                    `json:"message"`
	UnavailableReason          string                                                    `json:"unavailableReason"`
	ProfileBio                 *TweetAdvancedSearchTweetAuthorProfileBio                 `json:"profile_bio"`
}

type TweetAdvancedSearchTweetEntitiesHashtag struct {
	Indices []int  `json:"indices"`
	Text    string `json:"text"`
}

type TweetAdvancedSearchTweetEntitiesURL struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
	URL         string `json:"url"`
}

type TweetAdvancedSearchTweetEntitiesUserMention struct {
	IDStr      string `json:"id_str"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
}

type TweetAdvancedSearchTweetEntities struct {
	Hashtags     []*TweetAdvancedSearchTweetEntitiesHashtag     `json:"hashtags"`
	URLs         []*TweetAdvancedSearchTweetEntitiesURL         `json:"urls"`
	UserMentions []*TweetAdvancedSearchTweetEntitiesUserMention `json:"user_mentions"`
}

type TweetAdvancedSearchTweet struct {
	Type              string                            `json:"type"`
	ID                string                            `json:"id"`
	URL               string                            `json:"url"`
	Text              string                            `json:"text"`
	Source            string                            `json:"source"`
	RetweetCount      int                               `json:"retweetCount"`
	ReplyCount        int                               `json:"replyCount"`
	LikeCount         int                               `json:"likeCount"`
	QuoteCount        int                               `json:"quoteCount"`
	ViewCount         int                               `json:"viewCount"`
	CreatedAt         string                            `json:"createdAt"`
	Lang              string                            `json:"lang"`
	BookmarkCount     int                               `json:"bookmarkCount"`
	IsReply           bool                              `json:"isReply"`
	InReplyToID       string                            `json:"inReplyToId"`
	ConversationID    string                            `json:"conversationId"`
	DisplayTextRange  []int                             `json:"displayTextRange"`
	InReplyToUserID   string                            `json:"inReplyToUserId"`
	InReplyToUsername string                            `json:"inReplyToUsername"`
	Author            *TweetAdvancedSearchTweetAuthor   `json:"author"`
	Entities          *TweetAdvancedSearchTweetEntities `json:"entities"`
	QuotedTweet       *TweetAdvancedSearchTweet         `json:"quoted_tweet"`
	RetweetedTweet    *TweetAdvancedSearchTweet         `json:"retweeted_tweet"`
	IsLimitedReply    bool                              `json:"isLimitedReply"`
}

type TweetAdvancedSearchResponse struct {
	Tweets      []*TweetAdvancedSearchTweet `json:"tweets"`
	HasNextPage bool                        `json:"has_next_page"`
	NextCursor  string                      `json:"next_cursor"`
	Status      string                      `json:"status"`
	Message     string                      `json:"msg"`
}

func (t *twitterApi) TweetAdvancedSearch(query string, cursor *string) (*TweetAdvancedSearchResponse, error) {
	if query == "" {
		return nil, errors.New("query is empty")
	}

	queryParts := []string{}
	queryParts = append(queryParts, "query="+query)
	if cursor != nil && *cursor != "" {
		queryParts = append(queryParts, "cursor="+*cursor)
	}
	url := twitterDomainURI + "/tweet/advanced_search"
	if len(queryParts) > 0 {
		url += "?" + strings.Join(queryParts, "&")
	}

	ctx1, cancel1 := context.WithTimeout(t.ctx, time.Second*10)
	defer cancel1()

	jsonData, resp, err := getDataWithHeader(ctx1, t.httpClient, url, t.headers)
	if err != nil {
		slog.Error("TweetAdvancedSearch failed", "err", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		slog.Error("TweetAdvancedSearch failed", "statusCode", resp.StatusCode, "body", string(jsonData))
		return nil, errors.New("TweetAdvancedSearch failed")
	}

	response := &TweetAdvancedSearchResponse{}
	if err = jsoniter.Unmarshal(jsonData, &response); err != nil {
		slog.Error("TweetAdvancedSearch failed", "err", err)
		return nil, err
	}
	if response.Status != "success" {
		slog.Error("TweetAdvancedSearch failed", "status", response.Status, "message", response.Message)
		return nil, errors.New("TweetAdvancedSearch failed")
	}

	return response, nil
}
