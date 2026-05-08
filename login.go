// Doc https://docs.twitterapi.io/api-reference/endpoint/user_login_v2

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

type Cookies string

type loginRequest struct {
	Username   string  `json:"user_name"`
	Email      string  `json:"email"`
	Password   string  `json:"password"`
	Proxy      string  `json:"proxy"`
	TotpSecret *string `json:"totp_secret,omitempty"`
}

type loginResponse struct {
	Cookies Cookies `json:"login_cookies"`
	Status  string  `json:"status"`
	Message string  `json:"message"`
}

// Login
func (t *twitterApi) Login(username, email, password string, totpSecret *string) error {
	if t.cookies != "" {
		return nil
	}

	if t.proxy == "" {
		return errors.New("proxy is empty, please set WithProxy")
	}

	if username == "" {
		return errors.New("username is empty")
	}
	if email == "" {
		return errors.New("email is empty")
	}
	if password == "" {
		return errors.New("password is empty")
	}

	request := &loginRequest{
		Username: username,
		Email:    email,
		Password: password,
		Proxy:    t.proxy,
	}

	if totpSecret != nil {
		request.TotpSecret = totpSecret
	}

	jsonData, _ := jsoniter.Marshal(request)

	ctx1, cancel1 := context.WithTimeout(t.ctx, time.Second*5)
	defer cancel1()
	headers := maps.Clone(t.headers)
	headers["Content-Type"] = "application/json"
	jsonData, resp, err := postDataWithHeader(ctx1, t.httpClient, twitterDomainURI+"/user_login_v2", bytes.NewReader(jsonData), headers)
	if err != nil {
		slog.Error("Login failed", "err", err)
		return err
	}

	slog.Info("Login response", "jsonData", string(jsonData))

	if resp.StatusCode != http.StatusOK {
		slog.Error("Login failed", "statusCode", resp.StatusCode, "body", string(jsonData))
		return errors.New("Login failed")
	}

	response := &loginResponse{}

	if err = jsoniter.Unmarshal(jsonData, &response); err != nil {
		slog.Error("Login failed", "err", err)
		return err
	}

	if response.Status != "success" {
		slog.Error("Login failed", "status", response.Status, "message", response.Message)
		return errors.New("Login failed")
	}

	if response.Cookies == "" {
		slog.Error("Login failed", "cookies", response.Cookies, "message", response.Message)
		return errors.New("Cookie is empty")
	}

	t.cookies = response.Cookies

	return nil
}
