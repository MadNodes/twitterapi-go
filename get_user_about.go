// Doc https://docs.twitterapi.io/api-reference/endpoint/get_user_about

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

type GetUserAboutDataAffiliatesHighlightedLabelLabelBadge struct {
	URL string `json:"url"`
}

type GetUserAboutDataAffiliatesHighlightedLabelLabelURL struct {
	URL     string `json:"url"`
	URLType string `json:"urlType"`
}

type GetUserAboutDataAffiliatesHighlightedLabelLabel struct {
	Badge                *GetUserAboutDataAffiliatesHighlightedLabelLabelBadge `json:"badge"`
	Description          string                                                `json:"description"`
	URL                  *GetUserAboutDataAffiliatesHighlightedLabelLabelURL   `json:"url"`
	UserLabelDisplayType string                                                `json:"userLabelDisplayType"`
	UserLabelType        string                                                `json:"userLabelType"`
}

type GetUserAboutDataAffiliatesHighlightedLabel struct {
	Label *GetUserAboutDataAffiliatesHighlightedLabelLabel `json:"label"`
}

type GetUserAboutDataAboutProfileUsernameChanges struct {
	Count string `json:"count"`
}

type GetUserAboutDataAboutProfile struct {
	AccountBasedIn    string                                       `json:"account_based_in"`
	LocationAccurate  bool                                         `json:"location_accurate"`
	LearnMoreURL      string                                       `json:"learn_more_url"`
	AffiliateUsername string                                       `json:"affiliate_username"`
	Source            string                                       `json:"source"`
	UsernameChanges   *GetUserAboutDataAboutProfileUsernameChanges `json:"username_changes"`
}

type GetUserAboutDataIDentityProfileLabelsHighlightedLabelLabelBadge struct {
	URL string `json:"url"`
}

type GetUserAboutDataIDentityProfileLabelsHighlightedLabelLabelURL struct {
	URL     string `json:"url"`
	URLType string `json:"urlType"`
}

type GetUserAboutDataIDentityProfileLabelsHighlightedLabelLabel struct {
	Description          string                                                           `json:"description"`
	Badge                *GetUserAboutDataIDentityProfileLabelsHighlightedLabelLabelBadge `json:"badge"`
	URL                  *GetUserAboutDataIDentityProfileLabelsHighlightedLabelLabelURL   `json:"url"`
	UserLabelDisplayType string                                                           `json:"userLabelDisplayType"`
	UserLabelType        string                                                           `json:"userLabelType"`
}

type GetUserAboutDataIDentityProfileLabelsHighlightedLabel struct {
	Label *GetUserAboutDataIDentityProfileLabelsHighlightedLabelLabel `json:"label"`
}

type GetUserAboutDataVerificationInfo struct {
	ID                 string `json:"id"`
	IsIdentityVerified bool   `json:"is_identity_verified"`
}

type GetUserAboutData struct {
	ID                                    string                                                 `json:"id"`
	Name                                  string                                                 `json:"name"`
	UserName                              string                                                 `json:"userName"`
	CreatedAt                             string                                                 `json:"createdAt"`
	ProfilePicture                        string                                                 `json:"profilePicture"` // added: from live response
	IsBlueVerified                        bool                                                   `json:"isBlueVerified"`
	IsVerified                            bool                                                   `json:"isVerified"` // added: from live response
	Protected                             bool                                                   `json:"protected"`
	AffiliatesHighlightedLabel            *GetUserAboutDataAffiliatesHighlightedLabel            `json:"affiliates_highlighted_label"`
	AboutProfile                          *GetUserAboutDataAboutProfile                          `json:"about_profile"`
	IDentityProfileLabelsHighlightedLabel *GetUserAboutDataIDentityProfileLabelsHighlightedLabel `json:"identity_profile_labels_highlighted_label"`
	VerificationInfo                      *GetUserAboutDataVerificationInfo                      `json:"verification_info"` // added: from live response
}

type GetUserAboutResponse struct {
	Data    *GetUserAboutData `json:"data"`
	Status  string            `json:"status"`
	Message string            `json:"msg"`
}

func (t *twitterApi) GetUserAbout(userName *string) (*GetUserAboutResponse, error) {
	if userName == nil || strings.TrimSpace(*userName) == "" {
		return nil, errors.New("userName is required")
	}

	vals := neturl.Values{}
	vals.Set("userName", *userName)
	url := userTwitterDomainURI + "_about?" + vals.Encode()

	ctx1, cancel1 := context.WithTimeout(t.ctx, time.Second*10)
	defer cancel1()

	jsonData, resp, err := getDataWithHeader(ctx1, t.httpClient, url, t.headers)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			slog.Error("GetUserAbout request timed out", "url", url)
			return nil, errors.New("GetUserAbout request timed out")
		}
		slog.Error("GetUserAbout failed", "err", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		slog.Error("GetUserAbout failed", "statusCode", resp.StatusCode, "body", string(jsonData))
		return nil, errors.New("GetUserAbout failed")
	}

	response := &GetUserAboutResponse{}
	if err = jsoniter.Unmarshal(jsonData, &response); err != nil {
		slog.Error("GetUserAbout failed", "err", err)
		return nil, err
	}
	if response.Status != "success" {
		slog.Error("GetUserAbout failed", "status", response.Status, "message", response.Message)
		return nil, errors.New(response.Message)
	}

	return response, nil
}
