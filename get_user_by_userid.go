package twitterapi

import (
	"errors"
	"strings"
)

// GetUserInfoByID is a convenience wrapper around BatchGetUserInfoByUserIds
// to fetch a single user's profile by their numeric user ID.
func (t *TwitterApi) GetUserInfoByID(userId string) (*BatchGetUserInfoByUserIdsUser, error) {
	if strings.TrimSpace(userId) == "" {
		return nil, errors.New("userId is required")
	}

	resp, err := t.BatchGetUserInfoByUserIds([]string{userId})
	if err != nil {
		return nil, err
	}
	if resp == nil || len(resp.Users) == 0 {
		return nil, errors.New("user not found")
	}
	return resp.Users[0], nil
}
