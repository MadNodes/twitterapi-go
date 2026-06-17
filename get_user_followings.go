// Doc https://docs.twitterapi.io/api-reference/endpoint/get_user_followings

package twitterapi

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	neturl "net/url"
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type GetUserFollowingsUser struct {
	ID                   string  `json:"id"`
	Name                 string  `json:"name"`
	ScreenName           string  `json:"screen_name"`
	UserName             string  `json:"userName"`
	Location             string  `json:"location"`
	URL                  *string `json:"url"`
	Description          string  `json:"description"`
	Email                *string `json:"email"`
	Protected            bool    `json:"protected"`
	Verified             bool    `json:"verified"`
	FollowersCount       int     `json:"followers_count"`
	FollowingCount       int     `json:"following_count"`
	FriendsCount         int     `json:"friends_count"`
	FavouritesCount      int     `json:"favourites_count"`
	StatusesCount        int     `json:"statuses_count"`
	MediaTweetsCount     int     `json:"media_tweets_count"`
	CreatedAt            string  `json:"created_at"`
	ProfileBannerURL     *string `json:"profile_banner_url"`
	ProfileImageURLHTTPS string  `json:"profile_image_url_https"`
	CanDM                bool    `json:"can_dm"`
}

type GetUserFollowingsResponse struct {
	Followings  []*GetUserFollowingsUser `json:"followings"` // fixed: was followers
	HasNextPage bool                     `json:"has_next_page"`
	NextCursor  string                   `json:"next_cursor"`
	Status      string                   `json:"status"`
	Message     string                   `json:"msg"`
	Code        int                      `json:"code"`
}

func (t *TwitterApi) GetUserFollowings(userName string, pageSize *int, cursor *string) (*GetUserFollowingsResponse, error) {
	if strings.TrimSpace(userName) == "" {
		return nil, errors.New("userName is required")
	}

	vals := neturl.Values{}
	vals.Set("userName", userName)
	if pageSize != nil {
		vals.Set("pageSize", strconv.Itoa(*pageSize))
	}
	if cursor != nil && *cursor != "" {
		vals.Set("cursor", *cursor)
	}
	url := userTwitterDomainURI + "/followings?" + vals.Encode()

	jsonData, resp, err := t.getDataWithHeader(t.ctx, url, t.headers)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			slog.Error("GetUserFollowings request timed out", "url", url)
			return nil, errors.New("GetUserFollowings request timed out")
		}
		slog.Error("GetUserFollowings failed", "err", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		slog.Error("GetUserFollowings failed", "statusCode", resp.StatusCode, "body", string(jsonData))
		return nil, errors.New("GetUserFollowings failed")
	}

	response := &GetUserFollowingsResponse{}
	if err = jsoniter.Unmarshal(jsonData, &response); err != nil {
		slog.Error("GetUserFollowings failed", "err", err)
		return nil, err
	}
	if response.Status != "success" {
		slog.Error("GetUserFollowings failed", "status", response.Status, "message", response.Message)
		return nil, errors.New(response.Message)
	}

	return response, nil
}
