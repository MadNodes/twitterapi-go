// Doc https://docs.twitterapi.io/api-reference/endpoint/get_user_timeline

package twitterapi

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type GetUserTimelineURL struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
	URL         string `json:"url"`
}

type GetUserTimelineURLs struct {
	URLs []*GetUserTimelineURL `json:"urls"`
}

type GetUserTimelineAuthorEntities struct {
	Description *GetUserTimelineURLs `json:"description"`
	URL         *GetUserTimelineURLs `json:"url"`
}

type GetUserTimelineProfileBioEntities struct {
	Description *GetUserTimelineURLs `json:"description"`
	URL         *GetUserTimelineURLs `json:"url"`
}

type GetUserTimelineProfileBio struct {
	Description string                             `json:"description"`
	Entities    *GetUserTimelineProfileBioEntities `json:"entities"`
}

type GetUserTimelineAffiliatesHighlightedLabel map[any]any

type GetUserTimelineAuthor struct {
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
	Entities                   *GetUserTimelineAuthorEntities             `json:"entities"`
	FastFollowersCount         int                                        `json:"fastFollowersCount"`
	FavouritesCount            int                                        `json:"favouritesCount"`
	HasCustomTimelines         bool                                       `json:"hasCustomTimelines"`
	IsTranslator               bool                                       `json:"isTranslator"`
	MediaCount                 int                                        `json:"mediaCount"`
	StatusesCount              int                                        `json:"statusesCount"`
	WithheldInCountries        []string                                   `json:"withheldInCountries"`
	AffiliatesHighlightedLabel *GetUserTimelineAffiliatesHighlightedLabel `json:"affiliatesHighlightedLabel"`
	PossiblySensitive          bool                                       `json:"possiblySensitive"`
	PinnedTweetIDs             []string                                   `json:"pinnedTweetIds"`
	ProfileBio                 *GetUserTimelineProfileBio                 `json:"profile_bio"`
	IsAutomated                bool                                       `json:"isAutomated"`
	AutomatedBy                *string                                    `json:"automatedBy"`
}

type GetUserTimelineHashtag struct {
	Indices []int  `json:"indices"`
	Text    string `json:"text"`
}

type GetUserTimelineSymbol struct {
	Indices []int  `json:"indices"`
	Text    string `json:"text"`
}

type GetUserTimelineMention struct {
	IDStr      string `json:"id_str"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
}

type GetUserTimelineTweetEntities struct {
	Hashtags     []*GetUserTimelineHashtag `json:"hashtags"`
	Symbols      []*GetUserTimelineSymbol  `json:"symbols"`
	URLs         []*GetUserTimelineURL     `json:"urls"`
	UserMentions []*GetUserTimelineMention `json:"user_mentions"`
}

type GetUserTimelineMediaFeaturesFaces struct {
	X int `json:"x"`
	Y int `json:"y"`
	H int `json:"h"`
	W int `json:"w"`
}

type GetUserTimelineMediaFeatures struct {
	Large *GetUserTimelineMediaFacesList `json:"large"`
	Orig  *GetUserTimelineMediaFacesList `json:"orig"`
}

type GetUserTimelineMediaFacesList struct {
	Faces []*GetUserTimelineMediaFeaturesFaces `json:"faces"`
}

type GetUserTimelineExtMediaAvailability struct {
	Status string `json:"status"`
}

type GetUserTimelineMediaResultsResult struct {
	Typename string `json:"__typename"`
	ID       string `json:"id"`
	MediaKey string `json:"media_key"`
}

type GetUserTimelineMediaResults struct {
	ID     string                             `json:"id"`
	Result *GetUserTimelineMediaResultsResult `json:"result"`
}

type GetUserTimelineAdditionalMediaInfo struct {
	Monetizable bool `json:"monetizable"`
}

type GetUserTimelineFocusRect struct {
	X int `json:"x"`
	Y int `json:"y"`
	H int `json:"h"`
	W int `json:"w"`
}

type GetUserTimelineOriginalInfo struct {
	FocusRects []*GetUserTimelineFocusRect `json:"focus_rects"`
	Height     int                         `json:"height"`
	Width      int                         `json:"width"`
}

type GetUserTimelineMediaSizesLarge struct {
	H int `json:"h"`
	W int `json:"w"`
}

type GetUserTimelineMediaSizes struct {
	Large *GetUserTimelineMediaSizesLarge `json:"large"`
}

type GetUserTimelineVideoVariant struct {
	Bitrate     int    `json:"bitrate"`
	ContentType string `json:"content_type"`
	URL         string `json:"url"`
}

type GetUserTimelineVideoInfo struct {
	AspectRatio    []int                          `json:"aspect_ratio"`
	DurationMillis int                            `json:"duration_millis"`
	Variants       []*GetUserTimelineVideoVariant `json:"variants"`
}

type GetUserTimelineExtendedEntitiesMedia struct {
	AdditionalMediaInfo  *GetUserTimelineAdditionalMediaInfo  `json:"additional_media_info"`
	DisplayURL           string                               `json:"display_url"`
	ExpandedURL          string                               `json:"expanded_url"`
	ExtMediaAvailability *GetUserTimelineExtMediaAvailability `json:"ext_media_availability"`
	Features             *GetUserTimelineMediaFeatures        `json:"features"`
	IDStr                string                               `json:"id_str"`
	Indices              []int                                `json:"indices"`
	MediaKey             string                               `json:"media_key"`
	MediaResults         *GetUserTimelineMediaResults         `json:"media_results"`
	MediaURLHTTPS        string                               `json:"media_url_https"`
	OriginalInfo         *GetUserTimelineOriginalInfo         `json:"original_info"`
	Sizes                *GetUserTimelineMediaSizes           `json:"sizes"`
	Type                 string                               `json:"type"`
	URL                  string                               `json:"url"`
	VideoInfo            *GetUserTimelineVideoInfo            `json:"video_info"`
}

type GetUserTimelineExtendedEntities struct {
	Media []*GetUserTimelineExtendedEntitiesMedia `json:"media"`
}

type GetUserTimelineTweet struct {
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
	Author            *GetUserTimelineAuthor           `json:"author"`
	ExtendedEntities  *GetUserTimelineExtendedEntities `json:"extendedEntities"`
	Card              *string                          `json:"card"`
	Place             map[string]any                   `json:"place"`
	Entities          *GetUserTimelineTweetEntities    `json:"entities"`
	QuotedTweet       *GetUserTimelineTweet            `json:"quoted_tweet"`
	RetweetedTweet    *GetUserTimelineTweet            `json:"retweeted_tweet"`
	IsLimitedReply    bool                             `json:"isLimitedReply"`
	CommunityInfo     *string                          `json:"communityInfo"`
	Article           *string                          `json:"article"`
}

type GetUserTimelineData struct {
	Tweets []*GetUserTimelineTweet `json:"tweets"`
}

type GetUserTimelineResponse struct {
	Status      string               `json:"status"`
	Code        int                  `json:"code"`
	Message     string               `json:"msg"`
	Data        *GetUserTimelineData `json:"data"`
	HasNextPage bool                 `json:"has_next_page"`
	NextCursor  string               `json:"next_cursor"`
}

func (t *twitterApi) GetUserTimeline(userId string, includeReplies, includeParentTweet *bool, cursor *string) (*GetUserTimelineResponse, error) {
	if userId == "" {
		return nil, errors.New("userId is empty")
	}

	url := userTwitterDomainURI + "/tweet_timeline?userId=" + userId
	if includeReplies != nil {
		url += "&includeReplies=" + strconv.FormatBool(*includeReplies)
	}
	if includeParentTweet != nil {
		url += "&includeParentTweet=" + strconv.FormatBool(*includeParentTweet)
	}
	if cursor != nil && *cursor != "" {
		url += "&cursor=" + *cursor
	}

	ctx1, cancel1 := context.WithTimeout(t.ctx, time.Second*10)
	defer cancel1()

	jsonData, resp, err := getDataWithHeader(ctx1, t.httpClient, url, t.headers)
	if err != nil {
		slog.Error("GetUserTimeline failed", "err", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		slog.Error("GetUserTimeline failed", "statusCode", resp.StatusCode, "body", string(jsonData))
		return nil, errors.New("GetUserTimeline failed")
	}

	response := &GetUserTimelineResponse{}
	if err = jsoniter.Unmarshal(jsonData, &response); err != nil {
		slog.Error("GetUserTimeline failed", "err", err)
		return nil, err
	}
	if response.Status != "success" {
		slog.Error("GetUserTimeline failed", "status", response.Status, "message", response.Message)
		return nil, errors.New("GetUserTimeline failed")
	}

	return response, nil
}
