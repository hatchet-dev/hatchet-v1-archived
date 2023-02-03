package fileclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"github.com/gorilla/schema"
	"github.com/hatchet-dev/hatchet/api/v1/types"
)

// FileClient is a Hatchet API client responsible for uploading/downloading files, because
// the auto-generated swagger clients don't properly support file uploads
type FileClient struct {
	Token   string
	BaseURL string

	HTTPClient *http.Client
}

func NewFileClient(baseURL, token string) *FileClient {
	return &FileClient{
		Token:   token,
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}

func (c *FileClient) UploadPlanZIPFile(teamID, moduleID, moduleRunID string, fileBytes []byte) (*types.APIError, error) {
	return c.sendUploadRequest(
		fmt.Sprintf("/api/v1/teams/%s/modules/%s/runs/%s/plan/zip", teamID, moduleID, moduleRunID),
		"terraform.tfplan",
		fileBytes,
	)
}

func (c *FileClient) GetPlanByCommitSHA(teamID, moduleID, moduleRunID string) (io.ReadCloser, *types.APIError, error) {
	return c.sendDownloadRequest(
		fmt.Sprintf("/api/v1/teams/%s/modules/%s/runs/%s/plan/sha", teamID, moduleID, moduleRunID),
		nil,
	)
}

func (c *FileClient) sendUploadRequest(relPath string, fileName string, fileBytes []byte) (*types.APIError, error) {
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)

	part, err := w.CreateFormFile("file", filepath.Base(fileName))
	if err != nil {
		return nil, err
	}
	_, err = part.Write(fileBytes)
	if err != nil {
		return nil, err
	}

	w.Close()

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s%s", c.BaseURL, relPath),
		body,
	)

	if err != nil {
		return nil, err
	}

	// Set the Boundary in the Content-Type
	req.Header.Add("Content-Type", w.FormDataContentType())

	// Set Content-Length
	req.Header.Add("Content-Length", fmt.Sprintf("%d", body.Len()))

	req.Header.Set("Accept", "application/json")

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))

	res, err := c.HTTPClient.Do(req)

	if err != nil {
		return nil, err
	}

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errRes types.APIError
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return &errRes, nil
		}

		return nil, fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	return nil, nil
}

func (c *FileClient) sendDownloadRequest(relPath string, data interface{}) (io.ReadCloser, *types.APIError, error) {
	vals := make(map[string][]string)
	err := schema.NewEncoder().Encode(data, vals)

	urlVals := url.Values(vals)
	encodedURLVals := urlVals.Encode()
	var req *http.Request

	if encodedURLVals != "" {
		req, err = http.NewRequest(
			"GET",
			fmt.Sprintf("%s%s?%s", c.BaseURL, relPath, encodedURLVals),
			nil,
		)
	} else {
		req, err = http.NewRequest(
			"GET",
			fmt.Sprintf("%s%s", c.BaseURL, relPath),
			nil,
		)
	}

	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Accept", "application/octet-stream")

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))

	res, err := c.HTTPClient.Do(req)

	if err != nil {
		return nil, nil, err
	}

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errRes types.APIError
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return nil, &errRes, nil
		}

		return nil, nil, fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	return res.Body, nil, nil
}
