// Doc https://docs.twitterapi.io/api-reference/endpoint/get_space_detail

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

type GetSpaceDetailDataSettings struct {
	ConversationControls        int  `json:"conversation_controls"`
	DisallowJoin                bool `json:"disallow_join"`
	IsEmployeeOnly              bool `json:"is_employee_only"`
	IsLocked                    bool `json:"is_locked"`
	IsMuted                     bool `json:"is_muted"`
	IsSpaceAvailableForClipping bool `json:"is_space_available_for_clipping"`
	IsSpaceAvailableForReplay   bool `json:"is_space_available_for_replay"`
	NoIncognito                 bool `json:"no_incognito"`
	NarrowCastSpaceType         int  `json:"narrow_cast_space_type"`
	MaxGuestSessions            int  `json:"max_guest_sessions"`
	MaxAdminCapacity            int  `json:"max_admin_capacity"`
}

type GetSpaceDetailDataStats struct {
	TotalReplayWatched int `json:"total_replay_watched"`
	TotalLiveListeners int `json:"total_live_listeners"`
	TotalParticipants  int `json:"total_participants"`
}

type GetSpaceDetailDataCreatorAffiliatesHighlightedLabel map[string]any

type GetSpaceDetailDataCreator struct {
	ID                         string                                               `json:"id"`
	Name                       string                                               `json:"name"`
	UserName                   string                                               `json:"userName"`
	Location                   string                                               `json:"location"`
	URL                        string                                               `json:"url"`
	Description                string                                               `json:"description"`
	Protected                  bool                                                 `json:"protected"`
	IsVerified                 bool                                                 `json:"isVerified"`
	IsBlueVerified             bool                                                 `json:"isBlueVerified"`
	VerifiedType               string                                               `json:"verifiedType"`
	Followers                  int                                                  `json:"followers"`
	Following                  int                                                  `json:"following"`
	FavouritesCount            int                                                  `json:"favouritesCount"`
	StatusesCount              int                                                  `json:"statusesCount"`
	MediaCount                 int                                                  `json:"mediaCount"`
	CreatedAt                  string                                               `json:"createdAt"`
	CoverPicture               string                                               `json:"coverPicture"`
	ProfilePicture             string                                               `json:"profilePicture"`
	CanDM                      bool                                                 `json:"canDm"`
	AffiliatesHighlightedLabel *GetSpaceDetailDataCreatorAffiliatesHighlightedLabel `json:"affiliatesHighlightedLabel"`
	IsAutomated                bool                                                 `json:"isAutomated"`
	AutomatedBy                string                                               `json:"automatedBy"`
}

type GetSpaceDetailDataParticipantsAdminAffiliatesHighlightedLabel map[string]any

type GetSpaceDetailDataParticipantsAdminParticipantInfo struct {
	PeriscopeUserID string `json:"periscope_user_id"`
	StartTime       string `json:"start_time"`
	IsMutedByAdmin  bool   `json:"is_muted_by_admin"`
	IsMutedByGuest  bool   `json:"is_muted_by_guest"`
}

type GetSpaceDetailDataParticipantsAdmin struct {
	ID                         string                                                         `json:"id"`
	Name                       string                                                         `json:"name"`
	UserName                   string                                                         `json:"userName"`
	Location                   string                                                         `json:"location"`
	URL                        string                                                         `json:"url"`
	Description                string                                                         `json:"description"`
	Protected                  bool                                                           `json:"protected"`
	IsVerified                 bool                                                           `json:"isVerified"`
	IsBlueVerified             bool                                                           `json:"isBlueVerified"`
	VerifiedType               string                                                         `json:"verifiedType"`
	Followers                  int                                                            `json:"followers"`
	Following                  int                                                            `json:"following"`
	FavouritesCount            int                                                            `json:"favouritesCount"`
	StatusesCount              int                                                            `json:"statusesCount"`
	MediaCount                 int                                                            `json:"mediaCount"`
	CreatedAt                  string                                                         `json:"createdAt"`
	CoverPicture               string                                                         `json:"coverPicture"`
	ProfilePicture             string                                                         `json:"profilePicture"`
	CanDM                      bool                                                           `json:"canDm"`
	AffiliatesHighlightedLabel *GetSpaceDetailDataParticipantsAdminAffiliatesHighlightedLabel `json:"affiliatesHighlightedLabel"`
	IsAutomated                bool                                                           `json:"isAutomated"`
	AutomatedBy                string                                                         `json:"automatedBy"`
	ParticipantInfo            *GetSpaceDetailDataParticipantsAdminParticipantInfo            `json:"participant_info"`
}

type GetSpaceDetailDataParticipantsSpeaker struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	UserName string `json:"userName"`
}

type GetSpaceDetailDataParticipantsListener struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GetSpaceDetailDataParticipants struct {
	Admins    []*GetSpaceDetailDataParticipantsAdmin    `json:"admins"`
	Speakers  []*GetSpaceDetailDataParticipantsSpeaker  `json:"speakers"`
	Listeners []*GetSpaceDetailDataParticipantsListener `json:"listeners"`
}

type GetSpaceDetailData struct {
	ID             string                          `json:"id"`
	Title          string                          `json:"title"`
	State          string                          `json:"state"`
	CreatedAt      string                          `json:"created_at"`
	ScheduledStart string                          `json:"scheduled_start"`
	UpdatedAt      string                          `json:"updated_at"`
	MediaKey       string                          `json:"media_key"`
	IsSubscribed   bool                            `json:"is_subscribed"`
	Settings       *GetSpaceDetailDataSettings     `json:"settings"`
	Stats          *GetSpaceDetailDataStats        `json:"stats"`
	Creator        *GetSpaceDetailDataCreator      `json:"creator"`
	Participants   *GetSpaceDetailDataParticipants `json:"participants"`
}

type GetSpaceDetailResponse struct {
	Data    *GetSpaceDetailData `json:"data"`
	Status  string              `json:"status"`
	Message string              `json:"msg"`
}

func (t *TwitterApi) GetSpaceDetail(spaceID *string) (*GetSpaceDetailResponse, error) {
	if spaceID == nil || strings.TrimSpace(*spaceID) == "" {
		return nil, errors.New("spaceID is required")
	}

	vals := neturl.Values{}
	vals.Set("space_id", *spaceID)
	url := twitterDomainURI + "/spaces/detail?" + vals.Encode()

	jsonData, resp, err := t.getDataWithHeader(t.ctx, url, t.headers)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			slog.Error("GetSpaceDetail request timed out", "url", url)
			return nil, errors.New("GetSpaceDetail request timed out")
		}
		slog.Error("GetSpaceDetail failed", "err", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		slog.Error("GetSpaceDetail failed", "statusCode", resp.StatusCode, "body", string(jsonData))
		return nil, errors.New("GetSpaceDetail failed")
	}

	response := &GetSpaceDetailResponse{}
	if err = jsoniter.Unmarshal(jsonData, &response); err != nil {
		slog.Error("GetSpaceDetail failed", "err", err)
		return nil, err
	}
	if response.Status != "success" {
		slog.Error("GetSpaceDetail failed", "status", response.Status, "message", response.Message)
		return nil, errors.New(response.Message)
	}

	return response, nil
}
