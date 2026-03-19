package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"0ximg.sh/cli/internal/models"
)

const renderEndpoint = "https://0ximg.sh/v1/render"

type renderResponse struct {
	URL        string `json:"url"`
	PreviewURL string `json:"previewUrl"`
}

func RenderImage(req models.RenderRequest) (string, string, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return "", "", fmt.Errorf("marshal render request: %w", err)
	}

	logDebugPayload(payload)

	resp, err := http.Post(renderEndpoint, "application/json", bytes.NewReader(payload))
	if err != nil {
		return "", "", fmt.Errorf("post render request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("read render response body: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		message := strings.TrimSpace(string(body))
		if message == "" {
			message = "empty response body"
		}
		return "", "", fmt.Errorf("render request failed with status %s: %s", resp.Status, message)
	}

	var renderResp renderResponse
	if err := json.Unmarshal(body, &renderResp); err != nil {
		return "", "", fmt.Errorf("decode render response: %w", err)
	}

	imageURL := strings.TrimSpace(renderResp.URL)
	if imageURL == "" {
		return "", "", fmt.Errorf("render response missing url")
	}

	previewURL := strings.TrimSpace(renderResp.PreviewURL)
	if previewURL == "" {
		previewURL = PreviewURL(imageURL)
	}

	return imageURL, previewURL, nil
}

func logDebugPayload(payload []byte) {
	if strings.TrimSpace(os.Getenv("DEBUG")) == "" {
		return
	}

	fmt.Fprintf(os.Stderr, "DEBUG render payload: %s\n", payload)
}

func PreviewURL(imageURL string) string {
	parsedEndpoint, err := url.Parse(renderEndpoint)
	if err != nil {
		return strings.TrimSpace(imageURL)
	}

	imageURL = strings.TrimSpace(imageURL)
	if imageURL == "" {
		return ""
	}

	imageName := path.Base(imageURL)
	renderID := strings.TrimSuffix(imageName, filepath.Ext(imageName))
	if strings.TrimSpace(renderID) == "" || renderID == "." || renderID == "/" {
		return imageURL
	}

	previewBase := *parsedEndpoint
	previewBase.Path = "/renders/" + renderID
	previewBase.RawQuery = ""
	previewBase.Fragment = ""

	return previewBase.String()
}
