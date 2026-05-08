// Doc https://docs.twitterapi.io/api-reference/endpoint/get_user_last_tweets

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

type GetUserLastTweetsURL struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
	URL         string `json:"url"`
}

type GetUserLastTweetsURLs struct {
	URLs []*GetUserLastTweetsURL `json:"urls"`
}

type GetUserLastTweetsAuthorEntities struct {
	Description *GetUserLastTweetsURLs `json:"description"`
	URL         *GetUserLastTweetsURLs `json:"url"`
}

type GetUserLastTweetsProfileBioMention struct {
	IDStr      string `json:"id_str"`
	Indices    []int  `json:"indices"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
}

type GetUserLastTweetsProfileBioDescription struct {
	URLs         []*GetUserLastTweetsURL               `json:"urls"`
	UserMentions []*GetUserLastTweetsProfileBioMention `json:"user_mentions"`
}

type GetUserLastTweetsProfileBioEntities struct {
	Description *GetUserLastTweetsProfileBioDescription `json:"description"`
	URL         *GetUserLastTweetsURLs                  `json:"url"`
}

type GetUserLastTweetsProfileBio struct {
	Description string                               `json:"description"`
	Entities    *GetUserLastTweetsProfileBioEntities `json:"entities"`
}

type GetUserLastTweetsAffiliatesHighlightedLabel map[any]any

type GetUserLastTweetsAuthor struct {
	Type                       string                                       `json:"type"`
	UserName                   string                                       `json:"userName"`
	URL                        string                                       `json:"url"`
	TwitterURL                 string                                       `json:"twitterUrl"`
	ID                         string                                       `json:"id"`
	Name                       string                                       `json:"name"`
	IsVerified                 bool                                         `json:"isVerified"`
	IsBlueVerified             bool                                         `json:"isBlueVerified"`
	VerifiedType               *string                                      `json:"verifiedType"`
	ProfilePicture             string                                       `json:"profilePicture"`
	CoverPicture               string                                       `json:"coverPicture"`
	Description                string                                       `json:"description"`
	Location                   string                                       `json:"location"`
	Followers                  int                                          `json:"followers"`
	Following                  int                                          `json:"following"`
	Status                     string                                       `json:"status"`
	CanDM                      bool                                         `json:"canDm"`
	CanMediaTag                bool                                         `json:"canMediaTag"`
	CreatedAt                  string                                       `json:"createdAt"`
	Entities                   *GetUserLastTweetsAuthorEntities             `json:"entities"`
	FastFollowersCount         int                                          `json:"fastFollowersCount"`
	FavouritesCount            int                                          `json:"favouritesCount"`
	HasCustomTimelines         bool                                         `json:"hasCustomTimelines"`
	IsTranslator               bool                                         `json:"isTranslator"`
	MediaCount                 int                                          `json:"mediaCount"`
	StatusesCount              int                                          `json:"statusesCount"`
	WithheldInCountries        []string                                     `json:"withheldInCountries"`
	AffiliatesHighlightedLabel *GetUserLastTweetsAffiliatesHighlightedLabel `json:"affiliatesHighlightedLabel"`
	PossiblySensitive          bool                                         `json:"possiblySensitive"`
	PinnedTweetIDs             []string                                     `json:"pinnedTweetIds"`
	ProfileBio                 *GetUserLastTweetsProfileBio                 `json:"profile_bio"`
	IsAutomated                bool                                         `json:"isAutomated"`
	AutomatedBy                *string                                      `json:"automatedBy"`
}

type GetUserLastTweetsHashtag struct {
	Indices []int  `json:"indices"`
	Text    string `json:"text"`
}

type GetUserLastTweetsSymbol struct {
	Indices []int  `json:"indices"`
	Text    string `json:"text"`
}

type GetUserLastTweetsTimestamp struct {
	Indices []int  `json:"indices"`
	Text    string `json:"text"`
}

type GetUserLastTweetsMention struct {
	IDStr      string `json:"id_str"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
}

type GetUserLastTweetsTweetEntities struct {
	Hashtags     []*GetUserLastTweetsHashtag   `json:"hashtags"`
	Symbols      []*GetUserLastTweetsSymbol    `json:"symbols"`
	Timestamps   []*GetUserLastTweetsTimestamp `json:"timestamps"`
	URLs         []*GetUserLastTweetsURL       `json:"urls"`
	UserMentions []*GetUserLastTweetsMention   `json:"user_mentions"`
}

type GetUserLastTweetsMediaFeaturesFaces struct {
	X int `json:"x"`
	Y int `json:"y"`
	H int `json:"h"`
	W int `json:"w"`
}

type GetUserLastTweetsMediaFeatures struct {
	Large *GetUserLastTweetsMediaFacesList `json:"large"`
	Orig  *GetUserLastTweetsMediaFacesList `json:"orig"`
}

type GetUserLastTweetsMediaFacesList struct {
	Faces []*GetUserLastTweetsMediaFeaturesFaces `json:"faces"`
}

type GetUserLastTweetsExtMediaAvailability struct {
	Status string `json:"status"`
}

type GetUserLastTweetsAllowDownloadStatus struct {
	AllowDownload bool `json:"allow_download"`
}

type GetUserLastTweetsMediaResultsResult struct {
	Typename string `json:"__typename"`
	ID       string `json:"id"`
	MediaKey string `json:"media_key"`
}

type GetUserLastTweetsMediaResults struct {
	ID     string                               `json:"id"`
	Result *GetUserLastTweetsMediaResultsResult `json:"result"`
}

type GetUserLastTweetsFocusRect struct {
	X int `json:"x"`
	Y int `json:"y"`
	H int `json:"h"`
	W int `json:"w"`
}

type GetUserLastTweetsOriginalInfo struct {
	FocusRects []*GetUserLastTweetsFocusRect `json:"focus_rects"`
	Height     int                           `json:"height"`
	Width      int                           `json:"width"`
}

type GetUserLastTweetsMediaSizesLarge struct {
	H int `json:"h"`
	W int `json:"w"`
}

type GetUserLastTweetsMediaSizes struct {
	Large *GetUserLastTweetsMediaSizesLarge `json:"large"`
}

type GetUserLastTweetsExtendedEntitiesMedia struct {
	AllowDownloadStatus  *GetUserLastTweetsAllowDownloadStatus  `json:"allow_download_status"`
	DisplayURL           string                                 `json:"display_url"`
	ExpandedURL          string                                 `json:"expanded_url"`
	ExtMediaAvailability *GetUserLastTweetsExtMediaAvailability `json:"ext_media_availability"`
	Features             *GetUserLastTweetsMediaFeatures        `json:"features"`
	IDStr                string                                 `json:"id_str"`
	Indices              []int                                  `json:"indices"`
	MediaKey             string                                 `json:"media_key"`
	MediaResults         *GetUserLastTweetsMediaResults         `json:"media_results"`
	MediaURLHTTPS        string                                 `json:"media_url_https"`
	OriginalInfo         *GetUserLastTweetsOriginalInfo         `json:"original_info"`
	Sizes                *GetUserLastTweetsMediaSizes           `json:"sizes"`
	Type                 string                                 `json:"type"`
	URL                  string                                 `json:"url"`
}

type GetUserLastTweetsExtendedEntities struct {
	Media []*GetUserLastTweetsExtendedEntitiesMedia `json:"media"`
}

type GetUserLastTweetsTweet struct {
	Type              string                             `json:"type"`
	ID                string                             `json:"id"`
	URL               string                             `json:"url"`
	TwitterURL        string                             `json:"twitterUrl"`
	Text              string                             `json:"text"`
	Source            string                             `json:"source"`
	RetweetCount      int                                `json:"retweetCount"`
	ReplyCount        int                                `json:"replyCount"`
	LikeCount         int                                `json:"likeCount"`
	QuoteCount        int                                `json:"quoteCount"`
	ViewCount         int                                `json:"viewCount"`
	CreatedAt         string                             `json:"createdAt"`
	Lang              string                             `json:"lang"`
	BookmarkCount     int                                `json:"bookmarkCount"`
	IsReply           bool                               `json:"isReply"`
	InReplyToID       *string                            `json:"inReplyToId"`
	ConversationID    string                             `json:"conversationId"`
	DisplayTextRange  []int                              `json:"displayTextRange"`
	InReplyToUserID   *string                            `json:"inReplyToUserId"`
	InReplyToUsername *string                            `json:"inReplyToUsername"`
	Author            *GetUserLastTweetsAuthor           `json:"author"`
	ExtendedEntities  *GetUserLastTweetsExtendedEntities `json:"extendedEntities"`
	Card              *string                            `json:"card"`
	Place             map[string]any                     `json:"place"`
	Entities          *GetUserLastTweetsTweetEntities    `json:"entities"`
	QuotedTweet       *GetUserLastTweetsTweet            `json:"quoted_tweet"`
	RetweetedTweet    *GetUserLastTweetsTweet            `json:"retweeted_tweet"`
	IsLimitedReply    bool                               `json:"isLimitedReply"`
	CommunityInfo     *string                            `json:"communityInfo"`
	Article           *string                            `json:"article"`
}

type GetUserLastTweetsData struct {
	PinTweet *GetUserLastTweetsTweet   `json:"pin_tweet"`
	Tweets   []*GetUserLastTweetsTweet `json:"tweets"`
}

type GetUserLastTweetsResponse struct {
	Status      string                 `json:"status"`
	Code        int                    `json:"code"`
	Message     string                 `json:"msg"`
	Data        *GetUserLastTweetsData `json:"data"`
	HasNextPage bool                   `json:"has_next_page"`
	NextCursor  string                 `json:"next_cursor"`
}

func (t *twitterApi) GetUserLastTweets(userId, userName *string, includeReplies *bool, cursor *string) (*GetUserLastTweetsResponse, error) {
	if (userId == nil || *userId == "") && (userName == nil || *userName == "") {
		return nil, errors.New("userId or userName is required")
	}

	url := userTwitterDomainURI + "/last_tweets?"
	if userId != nil && *userId != "" {
		url += "userId=" + *userId
	} else {
		url += "userName=" + *userName
	}
	if includeReplies != nil {
		url += "&includeReplies=" + strconv.FormatBool(*includeReplies)
	}
	if cursor != nil && *cursor != "" {
		url += "&cursor=" + *cursor
	}

	ctx1, cancel1 := context.WithTimeout(t.ctx, time.Second*10)
	defer cancel1()

	jsonData, resp, err := getDataWithHeader(ctx1, t.httpClient, url, t.headers)
	if err != nil {
		slog.Error("GetUserLastTweets failed", "err", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		slog.Error("GetUserLastTweets failed", "statusCode", resp.StatusCode, "body", string(jsonData))
		return nil, errors.New("GetUserLastTweets failed")
	}

	response := &GetUserLastTweetsResponse{}
	if err = jsoniter.Unmarshal(jsonData, &response); err != nil {
		slog.Error("GetUserLastTweets failed", "err", err)
		return nil, err
	}
	if response.Status != "success" {
		slog.Error("GetUserLastTweets failed", "status", response.Status, "message", response.Message)
		return nil, errors.New("GetUserLastTweets failed")
	}

	return response, nil
}
