// Doc https://docs.twitterapi.io/api-reference/endpoint/update_profile_v2

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

type updateProfileRequest struct {
	Cookies     Cookies `json:"login_cookies"`
	Proxy       string  `json:"proxy"`
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Location    string  `json:"location,omitempty"`
	URL         string  `json:"url,omitempty"`
}

type updateProfileResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// UpdateProfile
func (t *twitterApi) UpdateProfile(name, description, location, url *string) error {
	if t.proxy == "" {
		return errors.New("proxy is empty, please set WithProxy")
	}
	if t.cookies == "" {
		return errors.New("cookies is empty, please login first")
	}
	if name == nil && description == nil && location == nil && url == nil {
		return errors.New("profile update fields are empty")
	}

	request := &updateProfileRequest{
		Cookies: t.cookies,
		Proxy:   t.proxy,
	}
	if name != nil {
		request.Name = *name
	}
	if description != nil {
		request.Description = *description
	}
	if location != nil {
		request.Location = *location
	}
	if url != nil {
		request.URL = *url
	}

	jsonData, _ := jsoniter.Marshal(request)

	ctx1, cancel1 := context.WithTimeout(t.ctx, time.Second*10)
	defer cancel1()

	headers := maps.Clone(t.headers)
	headers["Content-Type"] = "application/json"
	jsonData, resp, err := patchDataWithHeader(ctx1, t.httpClient, twitterDomainURI+"/update_profile_v2", bytes.NewReader(jsonData), headers)
	if err != nil {
		slog.Error("UpdateProfile failed", "err", err)
		return err
	}

	slog.Info("UpdateProfile response", "jsonData", string(jsonData))
	if resp.StatusCode != http.StatusOK {
		slog.Error("UpdateProfile failed", "statusCode", resp.StatusCode, "body", string(jsonData))
		return errors.New("UpdateProfile failed")
	}

	response := &updateProfileResponse{}
	if err = jsoniter.Unmarshal(jsonData, &response); err != nil {
		slog.Error("UpdateProfile failed", "err", err)
		return err
	}
	if response.Status != "success" {
		slog.Error("UpdateProfile failed", "status", response.Status, "message", response.Message)
		return errors.New("UpdateProfile failed")
	}

	return nil
}
