// Doc https://docs.twitterapi.io/api-reference/endpoint/update_avatar_v2

package twitterapi

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type updateAvatarResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// UpdateAvatar
func (t *TwitterApi) UpdateAvatar(filename string, filebody io.Reader) error {
	if filename == "" || filebody == nil {
		return errors.New("filename is empty")
	}
	if t.proxy == "" {
		return errors.New("proxy is empty, please set WithProxy")
	}
	if t.cookies == "" {
		return errors.New("cookies is empty, please login first")
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		slog.Error("UpdateAvatar failed", "err", err)
		return err
	}
	if _, err = io.Copy(part, filebody); err != nil {
		slog.Error("UpdateAvatar failed", "err", err)
		return err
	}
	if err = writer.WriteField("proxy", t.proxy); err != nil {
		slog.Error("UpdateAvatar failed", "err", err)
		return err
	}
	if err = writer.WriteField("login_cookies", string(t.cookies)); err != nil {
		slog.Error("UpdateAvatar failed", "err", err)
		return err
	}
	if err = writer.Close(); err != nil {
		slog.Error("UpdateAvatar failed", "err", err)
		return err
	}

	ctx1, cancel1 := context.WithTimeout(t.ctx, time.Minute)
	defer cancel1()

	req, err := http.NewRequestWithContext(ctx1, http.MethodPatch, twitterDomainURI+"/update_avatar_v2", body)
	if err != nil {
		slog.Error("UpdateAvatar failed", "err", err)
		return err
	}
	for k, v := range t.headers {
		req.Header.Add(k, v)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := t.httpClient.Do(req)
	if err != nil {
		slog.Error("UpdateAvatar failed", "err", err)
		return err
	}
	defer resp.Body.Close()

	jsonData, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("UpdateAvatar failed", "err", err)
		return err
	}
	if resp.StatusCode != http.StatusOK {
		slog.Error("UpdateAvatar failed", "statusCode", resp.StatusCode, "body", string(jsonData))
		return errors.New("UpdateAvatar failed")
	}

	response := &updateAvatarResponse{}
	if err = jsoniter.Unmarshal(jsonData, &response); err != nil {
		slog.Error("UpdateAvatar failed", "err", err)
		return err
	}
	if response.Status != "success" {
		slog.Error("UpdateAvatar failed", "status", response.Status, "message", response.Message)
		return errors.New("UpdateAvatar failed")
	}

	return nil
}
