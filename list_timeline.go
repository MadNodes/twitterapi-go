// Doc https://docs.twitterapi.io/api-reference/endpoint/list_timeline

package twitterapi

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	neturl "net/url"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type ListTimelineTweetAuthorAffiliatesHighlightedLabel map[string]any

type ListTimelineTweetAuthorProfileBioEntitiesDescriptionURL struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
	URL         string `json:"url"`
}

type ListTimelineTweetAuthorProfileBioEntitiesDescription struct {
	URLs []*ListTimelineTweetAuthorProfileBioEntitiesDescriptionURL `json:"urls"`
}

type ListTimelineTweetAuthorProfileBioEntitiesURLURL struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
	URL         string `json:"url"`
}

type ListTimelineTweetAuthorProfileBioEntitiesURL struct {
	URLs []*ListTimelineTweetAuthorProfileBioEntitiesURLURL `json:"urls"`
}

type ListTimelineTweetAuthorProfileBioEntities struct {
	Description *ListTimelineTweetAuthorProfileBioEntitiesDescription `json:"description"`
	URL         *ListTimelineTweetAuthorProfileBioEntitiesURL         `json:"url"`
}

type ListTimelineTweetAuthorProfileBio struct {
	Description string                                     `json:"description"`
	Entities    *ListTimelineTweetAuthorProfileBioEntities `json:"entities"`
}

type ListTimelineTweetAuthor struct {
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
	AffiliatesHighlightedLabel *ListTimelineTweetAuthorAffiliatesHighlightedLabel `json:"affiliatesHighlightedLabel"`
	PossiblySensitive          bool                                               `json:"possiblySensitive"`
	PinnedTweetIDs             []string                                           `json:"pinnedTweetIds"`
	IsAutomated                bool                                               `json:"isAutomated"`
	AutomatedBy                string                                             `json:"automatedBy"`
	Unavailable                bool                                               `json:"unavailable"`
	Message                    string                                             `json:"message"`
	UnavailableReason          string                                             `json:"unavailableReason"`
	ProfileBio                 *ListTimelineTweetAuthorProfileBio                 `json:"profile_bio"`
}

type ListTimelineTweetEntitiesHashtag struct {
	Indices []int  `json:"indices"`
	Text    string `json:"text"`
}

type ListTimelineTweetEntitiesURL struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
	URL         string `json:"url"`
}

type ListTimelineTweetEntitiesUserMention struct {
	IDStr      string `json:"id_str"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
}

type ListTimelineTweetEntities struct {
	Hashtags     []*ListTimelineTweetEntitiesHashtag     `json:"hashtags"`
	URLs         []*ListTimelineTweetEntitiesURL         `json:"urls"`
	UserMentions []*ListTimelineTweetEntitiesUserMention `json:"user_mentions"`
}

type ListTimelineTweet struct {
	Type              string                     `json:"type"`
	ID                string                     `json:"id"`
	URL               string                     `json:"url"`
	TwitterURL        string                     `json:"twitterUrl"` // added: missing from original
	Text              string                     `json:"text"`
	Source            string                     `json:"source"`
	RetweetCount      int                        `json:"retweetCount"`
	ReplyCount        int                        `json:"replyCount"`
	LikeCount         int                        `json:"likeCount"`
	QuoteCount        int                        `json:"quoteCount"`
	ViewCount         int                        `json:"viewCount"`
	CreatedAt         string                     `json:"createdAt"`
	Lang              string                     `json:"lang"`
	BookmarkCount     int                        `json:"bookmarkCount"`
	IsReply           bool                       `json:"isReply"`
	InReplyToID       string                     `json:"inReplyToId"`
	ConversationID    string                     `json:"conversationId"`
	DisplayTextRange  []int                      `json:"displayTextRange"`
	InReplyToUserID   string                     `json:"inReplyToUserId"`
	InReplyToUsername string                     `json:"inReplyToUsername"`
	Author            *ListTimelineTweetAuthor   `json:"author"`
	Entities          *ListTimelineTweetEntities `json:"entities"`
	QuotedTweet       *ListTimelineTweet         `json:"quoted_tweet"`
	RetweetedTweet    *ListTimelineTweet         `json:"retweeted_tweet"`
	IsLimitedReply    bool                       `json:"isLimitedReply"`
}

type ListTimelineResponse struct {
	Tweets      []*ListTimelineTweet `json:"tweets"`
	HasNextPage bool                 `json:"has_next_page"`
	NextCursor  string               `json:"next_cursor"`
	Status      string               `json:"status"`
	Message     string               `json:"msg"`
}

func (t *TwitterApi) GetListTimeline(listID string, cursor *string) (*ListTimelineResponse, error) {
	if strings.TrimSpace(listID) == "" {
		return nil, errors.New("listID is required")
	}

	vals := neturl.Values{}
	vals.Set("listId", listID)
	if cursor != nil && *cursor != "" {
		vals.Set("cursor", *cursor)
	}
	url := listTwitterDomainURI + "/tweets_timeline?" + vals.Encode()

	jsonData, resp, err := t.getDataWithHeader(t.ctx, url, t.headers)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			slog.Error("GetListTimeline request timed out", "url", url)
			return nil, errors.New("GetListTimeline request timed out")
		}
		slog.Error("GetListTimeline failed", "err", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		slog.Error("GetListTimeline failed", "statusCode", resp.StatusCode, "body", string(jsonData))
		return nil, errors.New("GetListTimeline failed")
	}

	response := &ListTimelineResponse{}
	if err = jsoniter.Unmarshal(jsonData, &response); err != nil {
		slog.Error("GetListTimeline failed", "err", err)
		return nil, err
	}
	if response.Status != "success" {
		slog.Error("GetListTimeline failed", "status", response.Status, "message", response.Message)
		return nil, errors.New(response.Message)
	}

	return response, nil
}
