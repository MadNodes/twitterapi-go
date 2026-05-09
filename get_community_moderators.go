// Doc https://docs.twitterapi.io/api-reference/endpoint/get_community_moderators

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

// GetCommunityModeratorsMember matches the live API response, which uses
// snake_case fields instead of the camelCase fields used by the other
// community endpoints.
type GetCommunityModeratorsMember struct {
	ID                   string  `json:"id"`
	Name                 string  `json:"name"`
	ScreenName           string  `json:"screen_name"`
	Location             string  `json:"location"`
	URL                  string  `json:"url"`
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
	ProfileBannerURL     string  `json:"profile_banner_url"`
	ProfileImageURLHTTPS string  `json:"profile_image_url_https"`
	CanDM                bool    `json:"can_dm"`
	IsBlueVerified       bool    `json:"isBlueVerified"`
}

type GetCommunityModeratorsResponse struct {
	Moderators  []*GetCommunityModeratorsMember `json:"moderators"` // fixed: was members
	HasNext     bool                            `json:"has_next"`   // added: from live response
	HasNextPage bool                            `json:"has_next_page"`
	NextCursor  string                          `json:"next_cursor"`
	Status      string                          `json:"status"`
	Message     string                          `json:"msg"`
}

func (t *twitterApi) GetCommunityModerators(communityID string, cursor *string) (*GetCommunityModeratorsResponse, error) {
	if strings.TrimSpace(communityID) == "" {
		return nil, errors.New("communityID is required")
	}

	vals := neturl.Values{}
	vals.Set("community_id", communityID)
	if cursor != nil && *cursor != "" {
		vals.Set("cursor", *cursor)
	}
	url := twitterDomainURI + "/community/moderators?" + vals.Encode()

	ctx1, cancel1 := context.WithTimeout(t.ctx, time.Second*10)
	defer cancel1()

	jsonData, resp, err := getDataWithHeader(ctx1, t.httpClient, url, t.headers)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			slog.Error("GetCommunityModerators request timed out", "url", url)
			return nil, errors.New("GetCommunityModerators request timed out")
		}
		slog.Error("GetCommunityModerators failed", "err", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		slog.Error("GetCommunityModerators failed", "statusCode", resp.StatusCode, "body", string(jsonData))
		return nil, errors.New("GetCommunityModerators failed")
	}

	response := &GetCommunityModeratorsResponse{}
	if err = jsoniter.Unmarshal(jsonData, &response); err != nil {
		slog.Error("GetCommunityModerators failed", "err", err)
		return nil, err
	}
	if response.Status != "success" {
		slog.Error("GetCommunityModerators failed", "status", response.Status, "message", response.Message)
		return nil, errors.New(response.Message)
	}

	return response, nil
}
