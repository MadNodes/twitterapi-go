// Doc https://docs.twitterapi.io/api-reference/endpoint/get_user_mention

package twitterapi

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type GetUserMentionsURL struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
	URL         string `json:"url"`
}

type GetUserMentionsURLs struct {
	URLs []*GetUserMentionsURL `json:"urls"`
}

type GetUserMentionsAuthorEntities struct {
	Description *GetUserMentionsURLs `json:"description"`
	URL         *GetUserMentionsURLs `json:"url"`
}

type GetUserMentionsProfileBioMention struct {
	IDStr      string `json:"id_str"`
	Indices    []int  `json:"indices"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
}

type GetUserMentionsProfileBioDescription struct {
	Hashtags     []*GetUserMentionsHashtag           `json:"hashtags"`
	URLs         []*GetUserMentionsURL               `json:"urls"`
	UserMentions []*GetUserMentionsProfileBioMention `json:"user_mentions"`
}

type GetUserMentionsProfileBioEntities struct {
	Description *GetUserMentionsProfileBioDescription `json:"description"`
	URL         *GetUserMentionsURLs                  `json:"url"`
}

type GetUserMentionsProfileBio struct {
	Description string                             `json:"description"`
	Entities    *GetUserMentionsProfileBioEntities `json:"entities"`
}

type GetUserMentionsAffiliatesHighlightedLabel map[string]any

type GetUserMentionsAuthor struct {
	Type                       string                                     `json:"type"`
	UserName                   string                                     `json:"userName"`
	URL                        string                                     `json:"url"`
	TwitterURL                 string                                     `json:"twitterUrl"`
	ID                         string                                     `json:"id"`
	Name                       string                                     `json:"name"`
	IsVerified                 bool                                       `json:"isVerified"`
	IsBlueVerified             bool                                       `json:"isBlueVerified"`
	VerifiedType               *string                                    `json:"verifiedType"`
	ProfilePicture             string                                     `json:"profilePicture"`
	CoverPicture               string                                     `json:"coverPicture"`
	Description                string                                     `json:"description"`
	Location                   string                                     `json:"location"`
	Followers                  int                                        `json:"followers"`
	Following                  int                                        `json:"following"`
	Status                     string                                     `json:"status"`
	CanDM                      bool                                       `json:"canDm"`
	CanMediaTag                bool                                       `json:"canMediaTag"`
	CreatedAt                  string                                     `json:"createdAt"`
	Entities                   *GetUserMentionsAuthorEntities             `json:"entities"`
	FastFollowersCount         int                                        `json:"fastFollowersCount"`
	FavouritesCount            int                                        `json:"favouritesCount"`
	HasCustomTimelines         bool                                       `json:"hasCustomTimelines"`
	IsTranslator               bool                                       `json:"isTranslator"`
	MediaCount                 int                                        `json:"mediaCount"`
	StatusesCount              int                                        `json:"statusesCount"`
	WithheldInCountries        []string                                   `json:"withheldInCountries"`
	AffiliatesHighlightedLabel *GetUserMentionsAffiliatesHighlightedLabel `json:"affiliatesHighlightedLabel"`
	PossiblySensitive          bool                                       `json:"possiblySensitive"`
	PinnedTweetIDs             []string                                   `json:"pinnedTweetIds"`
	ProfileBio                 *GetUserMentionsProfileBio                 `json:"profile_bio"`
	IsAutomated                bool                                       `json:"isAutomated"`
	AutomatedBy                *string                                    `json:"automatedBy"`
}

type GetUserMentionsHashtag struct {
	Indices []int  `json:"indices"`
	Text    string `json:"text"`
}

type GetUserMentionsSymbol struct {
	Indices []int  `json:"indices"`
	Text    string `json:"text"`
}

type GetUserMentionsTimestamp struct {
	Indices []int  `json:"indices"`
	Text    string `json:"text"`
}

type GetUserMentionsMention struct {
	IDStr      string `json:"id_str"`
	Indices    []int  `json:"indices"` // added: from live response
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
}

type GetUserMentionsTweetEntities struct {
	Hashtags     []*GetUserMentionsHashtag   `json:"hashtags"`
	Symbols      []*GetUserMentionsSymbol    `json:"symbols"`
	Timestamps   []*GetUserMentionsTimestamp `json:"timestamps"`
	URLs         []*GetUserMentionsURL       `json:"urls"`
	UserMentions []*GetUserMentionsMention   `json:"user_mentions"`
}

type GetUserMentionsMediaFeaturesFaces struct {
	X int `json:"x"`
	Y int `json:"y"`
	H int `json:"h"`
	W int `json:"w"`
}

type GetUserMentionsMediaFeatures struct {
	Large *GetUserMentionsMediaFacesList `json:"large"`
	Orig  *GetUserMentionsMediaFacesList `json:"orig"`
}

type GetUserMentionsMediaFacesList struct {
	Faces []*GetUserMentionsMediaFeaturesFaces `json:"faces"`
}

type GetUserMentionsExtMediaAvailability struct {
	Status string `json:"status"`
}

type GetUserMentionsAllowDownloadStatus struct {
	AllowDownload bool `json:"allow_download"`
}

type GetUserMentionsMediaResultsResult struct {
	Typename string `json:"__typename"`
	ID       string `json:"id"`
	MediaKey string `json:"media_key"`
}

type GetUserMentionsMediaResults struct {
	ID     string                             `json:"id"`
	Result *GetUserMentionsMediaResultsResult `json:"result"`
}

type GetUserMentionsFocusRect struct {
	X int `json:"x"`
	Y int `json:"y"`
	H int `json:"h"`
	W int `json:"w"`
}

type GetUserMentionsOriginalInfo struct {
	FocusRects []*GetUserMentionsFocusRect `json:"focus_rects"`
	Height     int                         `json:"height"`
	Width      int                         `json:"width"`
}

type GetUserMentionsMediaSizesLarge struct {
	H int `json:"h"`
	W int `json:"w"`
}

type GetUserMentionsMediaSizes struct {
	Large *GetUserMentionsMediaSizesLarge `json:"large"`
}

type GetUserMentionsExtendedEntitiesMedia struct {
	AdditionalMediaInfo  json.RawMessage                      `json:"additional_media_info"` // polymorphic: see API docs
	AllowDownloadStatus  *GetUserMentionsAllowDownloadStatus  `json:"allow_download_status"`
	DisplayURL           string                               `json:"display_url"`
	ExpandedURL          string                               `json:"expanded_url"`
	ExtMediaAvailability *GetUserMentionsExtMediaAvailability `json:"ext_media_availability"`
	Features             *GetUserMentionsMediaFeatures        `json:"features"`
	IDStr                string                               `json:"id_str"`
	Indices              []int                                `json:"indices"`
	MediaKey             string                               `json:"media_key"`
	MediaResults         *GetUserMentionsMediaResults         `json:"media_results"`
	MediaURLHTTPS        string                               `json:"media_url_https"`
	OriginalInfo         *GetUserMentionsOriginalInfo         `json:"original_info"`
	Sizes                *GetUserMentionsMediaSizes           `json:"sizes"`
	SourceStatusIDStr    string                               `json:"source_status_id_str"` // added: from live response
	SourceUserIDStr      string                               `json:"source_user_id_str"`   // added: from live response
	Type                 string                               `json:"type"`
	URL                  string                               `json:"url"`
	VideoInfo            json.RawMessage                      `json:"video_info"` // polymorphic: see API docs
}

type GetUserMentionsExtendedEntities struct {
	Media []*GetUserMentionsExtendedEntitiesMedia `json:"media"`
}

type GetUserMentionsTweet struct {
	Type              string                           `json:"type"`
	ID                string                           `json:"id"`
	URL               string                           `json:"url"`
	TwitterURL        string                           `json:"twitterUrl"`
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
	InReplyToID       *string                          `json:"inReplyToId"`
	ConversationID    string                           `json:"conversationId"`
	DisplayTextRange  []int                            `json:"displayTextRange"`
	InReplyToUserID   *string                          `json:"inReplyToUserId"`
	InReplyToUsername *string                          `json:"inReplyToUsername"`
	Author            *GetUserMentionsAuthor           `json:"author"`
	ExtendedEntities  *GetUserMentionsExtendedEntities `json:"extendedEntities"`
	Card              json.RawMessage                  `json:"card"` // polymorphic: see API docs
	Place             map[string]any                   `json:"place"`
	Entities          *GetUserMentionsTweetEntities    `json:"entities"`
	QuotedTweet       *GetUserMentionsTweet            `json:"quoted_tweet"`
	RetweetedTweet    *GetUserMentionsTweet            `json:"retweeted_tweet"`
	IsLimitedReply    bool                             `json:"isLimitedReply"`
	CommunityInfo     *string                          `json:"communityInfo"`
	Article           json.RawMessage                  `json:"article"` // polymorphic: see API docs
}

type GetUserMentionsResponse struct {
	Tweets      []*GetUserMentionsTweet `json:"tweets"`
	HasNextPage bool                    `json:"has_next_page"`
	NextCursor  string                  `json:"next_cursor"`
	Status      string                  `json:"status"`
	Message     string                  `json:"msg"`
}

func (t *twitterApi) GetUserMentions(userName string, sinceTime, untilTime *int64, cursor *string) (*GetUserMentionsResponse, error) {
	if strings.TrimSpace(userName) == "" {
		return nil, errors.New("userName is required")
	}

	url := userTwitterDomainURI + "/mentions?userName=" + userName
	if sinceTime != nil {
		url += "&sinceTime=" + strconv.FormatInt(*sinceTime, 10)
	}
	if untilTime != nil {
		url += "&untilTime=" + strconv.FormatInt(*untilTime, 10)
	}
	if cursor != nil && *cursor != "" {
		url += "&cursor=" + *cursor
	}

	ctx1, cancel1 := context.WithTimeout(t.ctx, time.Second*10)
	defer cancel1()

	jsonData, resp, err := getDataWithHeader(ctx1, t.httpClient, url, t.headers)
	if err != nil {
		slog.Error("GetUserMentions failed", "err", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		slog.Error("GetUserMentions failed", "statusCode", resp.StatusCode, "body", string(jsonData))
		return nil, errors.New("GetUserMentions failed")
	}

	response := &GetUserMentionsResponse{}
	if err = jsoniter.Unmarshal(jsonData, &response); err != nil {
		slog.Error("GetUserMentions failed", "err", err)
		return nil, err
	}
	if response.Status != "success" {
		slog.Error("GetUserMentions failed", "status", response.Status, "message", response.Message)
		return nil, errors.New("GetUserMentions failed")
	}

	return response, nil
}
