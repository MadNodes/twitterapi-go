// Doc https://docs.twitterapi.io/api-reference/endpoint/get_all_community_tweets

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

type GetAllCommunityTweetsTweetAuthorAffiliatesHighlightedLabel map[string]any

type GetAllCommunityTweetsTweetAuthorProfileBioEntitiesDescriptionURL struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
	URL         string `json:"url"`
}

type GetAllCommunityTweetsTweetAuthorProfileBioEntitiesDescription struct {
	URLs []*GetAllCommunityTweetsTweetAuthorProfileBioEntitiesDescriptionURL `json:"urls"`
}

type GetAllCommunityTweetsTweetAuthorProfileBioEntitiesURLURL struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
	URL         string `json:"url"`
}

type GetAllCommunityTweetsTweetAuthorProfileBioEntitiesURL struct {
	URLs []*GetAllCommunityTweetsTweetAuthorProfileBioEntitiesURLURL `json:"urls"`
}

type GetAllCommunityTweetsTweetAuthorProfileBioEntities struct {
	Description *GetAllCommunityTweetsTweetAuthorProfileBioEntitiesDescription `json:"description"`
	URL         *GetAllCommunityTweetsTweetAuthorProfileBioEntitiesURL         `json:"url"`
}

type GetAllCommunityTweetsTweetAuthorProfileBio struct {
	Description string                                              `json:"description"`
	Entities    *GetAllCommunityTweetsTweetAuthorProfileBioEntities `json:"entities"`
}

type GetAllCommunityTweetsTweetAuthor struct {
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
	AffiliatesHighlightedLabel *GetAllCommunityTweetsTweetAuthorAffiliatesHighlightedLabel `json:"affiliatesHighlightedLabel"`
	PossiblySensitive          bool                                                        `json:"possiblySensitive"`
	PinnedTweetIDs             []string                                                    `json:"pinnedTweetIds"`
	IsAutomated                bool                                                        `json:"isAutomated"`
	AutomatedBy                string                                                      `json:"automatedBy"`
	Unavailable                bool                                                        `json:"unavailable"`
	Message                    string                                                      `json:"message"`
	UnavailableReason          string                                                      `json:"unavailableReason"`
	ProfileBio                 *GetAllCommunityTweetsTweetAuthorProfileBio                 `json:"profile_bio"`
}

type GetAllCommunityTweetsTweetEntitiesHashtag struct {
	Indices []int  `json:"indices"`
	Text    string `json:"text"`
}

type GetAllCommunityTweetsTweetEntitiesURL struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
	URL         string `json:"url"`
}

type GetAllCommunityTweetsTweetEntitiesUserMention struct {
	IDStr      string `json:"id_str"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
}

type GetAllCommunityTweetsTweetEntities struct {
	Hashtags     []*GetAllCommunityTweetsTweetEntitiesHashtag     `json:"hashtags"`
	URLs         []*GetAllCommunityTweetsTweetEntitiesURL         `json:"urls"`
	UserMentions []*GetAllCommunityTweetsTweetEntitiesUserMention `json:"user_mentions"`
}

type GetAllCommunityTweetsTweet struct {
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
	Author            *GetAllCommunityTweetsTweetAuthor   `json:"author"`
	Entities          *GetAllCommunityTweetsTweetEntities `json:"entities"`
	QuotedTweet       *GetAllCommunityTweetsTweet         `json:"quoted_tweet"`
	RetweetedTweet    *GetAllCommunityTweetsTweet         `json:"retweeted_tweet"`
	IsLimitedReply    bool                                `json:"isLimitedReply"`
}

type GetAllCommunityTweetsResponse struct {
	Tweets      []*GetAllCommunityTweetsTweet `json:"tweets"`
	HasNextPage bool                          `json:"has_next_page"`
	NextCursor  string                        `json:"next_cursor"`
}

func (t *twitterApi) GetAllCommunityTweets(query string, cursor *string) (*GetAllCommunityTweetsResponse, error) {
	if strings.TrimSpace(query) == "" {
		return nil, errors.New("query is required")
	}

	queryParts := []string{}
	queryParts = append(queryParts, "query="+query)
	if cursor != nil && *cursor != "" {
		queryParts = append(queryParts, "cursor="+*cursor)
	}
	url := twitterDomainURI + "/community/get_tweets_from_all_community"
	if len(queryParts) > 0 {
		url += "?" + strings.Join(queryParts, "&")
	}

	ctx1, cancel1 := context.WithTimeout(t.ctx, time.Second*10)
	defer cancel1()

	jsonData, resp, err := getDataWithHeader(ctx1, t.httpClient, url, t.headers)
	if err != nil {
		slog.Error("GetAllCommunityTweets failed", "err", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		slog.Error("GetAllCommunityTweets failed", "statusCode", resp.StatusCode, "body", string(jsonData))
		return nil, errors.New("GetAllCommunityTweets failed")
	}

	response := &GetAllCommunityTweetsResponse{}
	if err = jsoniter.Unmarshal(jsonData, &response); err != nil {
		slog.Error("GetAllCommunityTweets failed", "err", err)
		return nil, err
	}

	return response, nil
}
