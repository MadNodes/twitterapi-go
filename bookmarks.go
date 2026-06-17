// Doc https://docs.twitterapi.io/api-reference/endpoint/bookmarks_v2

package twitterapi

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"maps"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type bookmarksRequest struct {
	Cookies Cookies `json:"login_cookies"`
	Proxy   string  `json:"proxy"`
	Count   int     `json:"count,omitempty"`
	Cursor  string  `json:"cursor,omitempty"`
}

type BookmarksTweetURL struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
	URL         string `json:"url"`
}

type BookmarksTweetDescription struct {
	URLs []*BookmarksTweetURL `json:"urls"`
}

type BookmarksTweetProfileEntities struct {
	Description *BookmarksTweetDescription `json:"description"`
	URL         *BookmarksTweetDescription `json:"url"`
}

type BookmarksTweetProfileBio struct {
	Description string                         `json:"description"`
	Entities    *BookmarksTweetProfileEntities `json:"entities"`
}

type BookmarksTweetHashtag struct {
	Indices []int  `json:"indices"`
	Text    string `json:"text"`
}

type BookmarksTweetUserMention struct {
	IDStr      string `json:"id_str"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
}

type BookmarksTweetEntities struct {
	Hashtags     []*BookmarksTweetHashtag     `json:"hashtags"`
	URLs         []*BookmarksTweetURL         `json:"urls"`
	UserMentions []*BookmarksTweetUserMention `json:"user_mentions"`
}

type BookmarksTweetAffiliatesHighlightedLabel map[string]any

type BookmarksTweetAuthor struct {
	Type                       string                                    `json:"type"`
	UserName                   string                                    `json:"userName"`
	URL                        string                                    `json:"url"`
	ID                         string                                    `json:"id"`
	Name                       string                                    `json:"name"`
	IsBlueVerified             bool                                      `json:"isBlueVerified"`
	VerifiedType               string                                    `json:"verifiedType"`
	ProfilePicture             string                                    `json:"profilePicture"`
	CoverPicture               string                                    `json:"coverPicture"`
	Description                string                                    `json:"description"`
	Location                   string                                    `json:"location"`
	Followers                  int                                       `json:"followers"`
	Following                  int                                       `json:"following"`
	CanDM                      bool                                      `json:"canDm"`
	CreatedAt                  string                                    `json:"createdAt"`
	FavouritesCount            int                                       `json:"favouritesCount"`
	HasCustomTimelines         bool                                      `json:"hasCustomTimelines"`
	IsTranslator               bool                                      `json:"isTranslator"`
	MediaCount                 int                                       `json:"mediaCount"`
	StatusesCount              int                                       `json:"statusesCount"`
	WithheldInCountries        []string                                  `json:"withheldInCountries"`
	AffiliatesHighlightedLabel *BookmarksTweetAffiliatesHighlightedLabel `json:"affiliatesHighlightedLabel"`
	PossiblySensitive          bool                                      `json:"possiblySensitive"`
	PinnedTweetIDs             []string                                  `json:"pinnedTweetIds"`
	IsAutomated                bool                                      `json:"isAutomated"`
	AutomatedBy                string                                    `json:"automatedBy"`
	Unavailable                bool                                      `json:"unavailable"`
	Message                    string                                    `json:"message"`
	UnavailableReason          string                                    `json:"unavailableReason"`
	ProfileBio                 *BookmarksTweetProfileBio                 `json:"profile_bio"`
}

type BookmarksTweet struct {
	Type              string                  `json:"type"`
	ID                string                  `json:"id"`
	URL               string                  `json:"url"`
	Text              string                  `json:"text"`
	Source            string                  `json:"source"`
	RetweetCount      int                     `json:"retweetCount"`
	ReplyCount        int                     `json:"replyCount"`
	LikeCount         int                     `json:"likeCount"`
	QuoteCount        int                     `json:"quoteCount"`
	ViewCount         int                     `json:"viewCount"`
	CreatedAt         string                  `json:"createdAt"`
	Lang              string                  `json:"lang"`
	BookmarkCount     int                     `json:"bookmarkCount"`
	IsReply           bool                    `json:"isReply"`
	InReplyToID       string                  `json:"inReplyToId"`
	ConversationID    string                  `json:"conversationId"`
	DisplayTextRange  []int                   `json:"displayTextRange"`
	InReplyToUserID   string                  `json:"inReplyToUserId"`
	InReplyToUsername string                  `json:"inReplyToUsername"`
	Author            *BookmarksTweetAuthor   `json:"author"`
	Entities          *BookmarksTweetEntities `json:"entities"`
	QuotedTweet       *BookmarksTweet         `json:"quoted_tweet"`
	RetweetedTweet    *BookmarksTweet         `json:"retweeted_tweet"`
	IsLimitedReply    bool                    `json:"isLimitedReply"`
}

type BookmarksResponse struct {
	Tweets      []*BookmarksTweet `json:"tweets"`
	HasNextPage bool              `json:"has_next_page"`
	NextCursor  *string           `json:"next_cursor"`
	Detail      string            `json:"detail,omitempty"`
}

// GetBookmarks
func (t *TwitterApi) GetBookmarks(count int, cursor *string) (*BookmarksResponse, error) {
	if t.proxy == "" {
		return nil, errors.New("proxy is empty, please set WithProxy")
	}

	if t.cookies == "" {
		return nil, errors.New("cookies is empty, please login first")
	}

	if count <= 0 {
		count = 20
	}

	request := &bookmarksRequest{
		Cookies: t.cookies,
		Proxy:   t.proxy,
		Count:   count,
	}

	if cursor != nil {
		request.Cursor = *cursor
	}

	jsonData, _ := jsoniter.Marshal(request)

	url := twitterDomainURI + "/bookmarks_v2"

	ctx1, cancel1 := context.WithTimeout(t.ctx, time.Second*10)
	defer cancel1()
	headers := maps.Clone(t.headers)
	headers["Content-Type"] = "application/json"
	jsonData, resp, err := postDataWithHeader(ctx1, t.httpClient, url, bytes.NewReader(jsonData), headers)
	if err != nil {
		slog.Error("GetBookmarks failed", "err", err)
		return nil, err
	}

	//slog.Info("GetBookmarks response", "jsonData", string(jsonData))

	if resp.StatusCode != http.StatusOK {
		slog.Error("GetBookmarks failed", "statusCode", resp.StatusCode, "body", string(jsonData))
		return nil, errors.New("GetBookmarks failed")
	}

	response := &BookmarksResponse{}

	if err = jsoniter.Unmarshal(jsonData, &response); err != nil {
		slog.Error("GetBookmarks failed", "err", err)
		return nil, err
	}

	return response, nil
}
