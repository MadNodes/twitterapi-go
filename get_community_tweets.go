// Doc https://docs.twitterapi.io/api-reference/endpoint/get_community_tweets

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

type GetCommunityTweetsTweetAuthorAffiliatesHighlightedLabel map[string]any

type GetCommunityTweetsTweetAuthorProfileBioEntitiesDescriptionURL struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
	URL         string `json:"url"`
}

type GetCommunityTweetsTweetAuthorProfileBioEntitiesDescription struct {
	URLs []*GetCommunityTweetsTweetAuthorProfileBioEntitiesDescriptionURL `json:"urls"`
}

type GetCommunityTweetsTweetAuthorProfileBioEntitiesURLURL struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
	URL         string `json:"url"`
}

type GetCommunityTweetsTweetAuthorProfileBioEntitiesURL struct {
	URLs []*GetCommunityTweetsTweetAuthorProfileBioEntitiesURLURL `json:"urls"`
}

type GetCommunityTweetsTweetAuthorProfileBioEntities struct {
	Description *GetCommunityTweetsTweetAuthorProfileBioEntitiesDescription `json:"description"`
	URL         *GetCommunityTweetsTweetAuthorProfileBioEntitiesURL         `json:"url"`
}

type GetCommunityTweetsTweetAuthorProfileBio struct {
	Description string                                           `json:"description"`
	Entities    *GetCommunityTweetsTweetAuthorProfileBioEntities `json:"entities"`
}

type GetCommunityTweetsTweetAuthor struct {
	Type                       string                                                   `json:"type"`
	UserName                   string                                                   `json:"userName"`
	URL                        string                                                   `json:"url"`
	ID                         string                                                   `json:"id"`
	Name                       string                                                   `json:"name"`
	IsBlueVerified             bool                                                     `json:"isBlueVerified"`
	VerifiedType               string                                                   `json:"verifiedType"`
	ProfilePicture             string                                                   `json:"profilePicture"`
	CoverPicture               string                                                   `json:"coverPicture"`
	Description                string                                                   `json:"description"`
	Location                   string                                                   `json:"location"`
	Followers                  int                                                      `json:"followers"`
	Following                  int                                                      `json:"following"`
	CanDM                      bool                                                     `json:"canDm"`
	CreatedAt                  string                                                   `json:"createdAt"`
	FavouritesCount            int                                                      `json:"favouritesCount"`
	HasCustomTimelines         bool                                                     `json:"hasCustomTimelines"`
	IsTranslator               bool                                                     `json:"isTranslator"`
	MediaCount                 int                                                      `json:"mediaCount"`
	StatusesCount              int                                                      `json:"statusesCount"`
	WithheldInCountries        []string                                                 `json:"withheldInCountries"`
	AffiliatesHighlightedLabel *GetCommunityTweetsTweetAuthorAffiliatesHighlightedLabel `json:"affiliatesHighlightedLabel"`
	PossiblySensitive          bool                                                     `json:"possiblySensitive"`
	PinnedTweetIDs             []string                                                 `json:"pinnedTweetIds"`
	IsAutomated                bool                                                     `json:"isAutomated"`
	AutomatedBy                string                                                   `json:"automatedBy"`
	Unavailable                bool                                                     `json:"unavailable"`
	Message                    string                                                   `json:"message"`
	UnavailableReason          string                                                   `json:"unavailableReason"`
	ProfileBio                 *GetCommunityTweetsTweetAuthorProfileBio                 `json:"profile_bio"`
}

type GetCommunityTweetsTweetEntitiesHashtag struct {
	Indices []int  `json:"indices"`
	Text    string `json:"text"`
}

type GetCommunityTweetsTweetEntitiesURL struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
	URL         string `json:"url"`
}

type GetCommunityTweetsTweetEntitiesUserMention struct {
	IDStr      string `json:"id_str"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
}

type GetCommunityTweetsTweetEntities struct {
	Hashtags     []*GetCommunityTweetsTweetEntitiesHashtag     `json:"hashtags"`
	URLs         []*GetCommunityTweetsTweetEntitiesURL         `json:"urls"`
	UserMentions []*GetCommunityTweetsTweetEntitiesUserMention `json:"user_mentions"`
}

type GetCommunityTweetsTweet struct {
	Type              string                           `json:"type"`
	ID                string                           `json:"id"`
	URL               string                           `json:"url"`
	TwitterURL        string                           `json:"twitterUrl"` // added: missing from original
	Text              string                           `json:"text"`
	Source            string                           `json:"source"`
	RetweetCount      int                              `json:"retweetCount"`
	ReplyCount        int                              `json:"replyCount"`
	LikeCount         int                              `json:"likeCount"`
	QuoteCount        int                              `json:"quoteCount"`
	ViewCount         int                              `json:"viewCount"`
	CreatedAt         string                           `json:"createdAt"`
	Lang              string                           `json:"lang"`
	BookmarkCount     int                              `json:"bookmarkCount"`
	IsReply           bool                             `json:"isReply"`
	InReplyToID       string                           `json:"inReplyToId"`
	ConversationID    string                           `json:"conversationId"`
	DisplayTextRange  []int                            `json:"displayTextRange"`
	InReplyToUserID   string                           `json:"inReplyToUserId"`
	InReplyToUsername string                           `json:"inReplyToUsername"`
	Author            *GetCommunityTweetsTweetAuthor   `json:"author"`
	Entities          *GetCommunityTweetsTweetEntities `json:"entities"`
	QuotedTweet       *GetCommunityTweetsTweet         `json:"quoted_tweet"`
	RetweetedTweet    *GetCommunityTweetsTweet         `json:"retweeted_tweet"`
	IsLimitedReply    bool                             `json:"isLimitedReply"`
}

type GetCommunityTweetsResponse struct {
	Tweets      []*GetCommunityTweetsTweet `json:"tweets"`
	HasNext     bool                       `json:"has_next"` // added: from live response
	HasNextPage bool                       `json:"has_next_page"`
	NextCursor  string                     `json:"next_cursor"`
	Status      string                     `json:"status"`
	Message     string                     `json:"msg"`
}

func (t *twitterApi) GetCommunityTweets(communityID string, cursor *string) (*GetCommunityTweetsResponse, error) {
	if strings.TrimSpace(communityID) == "" {
		return nil, errors.New("communityID is required")
	}

	vals := neturl.Values{}
	vals.Set("community_id", communityID)
	if cursor != nil && *cursor != "" {
		vals.Set("cursor", *cursor)
	}
	url := twitterDomainURI + "/community/tweets?" + vals.Encode()

	ctx1, cancel1 := context.WithTimeout(t.ctx, time.Second*10)
	defer cancel1()

	jsonData, resp, err := getDataWithHeader(ctx1, t.httpClient, url, t.headers)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			slog.Error("GetCommunityTweets request timed out", "url", url)
			return nil, errors.New("GetCommunityTweets request timed out")
		}
		slog.Error("GetCommunityTweets failed", "err", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		slog.Error("GetCommunityTweets failed", "statusCode", resp.StatusCode, "body", string(jsonData))
		return nil, errors.New("GetCommunityTweets failed")
	}

	response := &GetCommunityTweetsResponse{}
	if err = jsoniter.Unmarshal(jsonData, &response); err != nil {
		slog.Error("GetCommunityTweets failed", "err", err)
		return nil, err
	}
	if response.Status != "success" {
		slog.Error("GetCommunityTweets failed", "status", response.Status, "message", response.Message)
		return nil, errors.New(response.Message)
	}

	return response, nil
}
