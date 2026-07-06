package update

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"env-editor/internal/model"
)

const (
	apiURL             = "https://api.github.com/repos/lzsbee/env-editor/releases/latest"
	fallbackReleaseURL = "https://github.com/lzsbee/env-editor/releases"
)

type githubRelease struct {
	TagName string `json:"tag_name"`
	HTMLURL string `json:"html_url"`
	Body    string `json:"body"`
}

// Check fetches the latest GitHub release and compares it with currentVersion.
func Check(currentVersion string) (model.UpdateInfo, error) {
	info := model.UpdateInfo{
		CurrentVersion: currentVersion,
		ReleaseURL:     fallbackReleaseURL,
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return info, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", fmt.Sprintf("env-editor/%s", currentVersion))

	resp, err := client.Do(req)
	if err != nil {
		return info, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return info, fmt.Errorf("github api: %s (%s)", resp.Status, strings.TrimSpace(string(body)))
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return info, err
	}

	latest := normalizeVersion(release.TagName)
	if latest == "" {
		return info, fmt.Errorf("empty release tag")
	}

	info.LatestVersion = latest
	if release.HTMLURL != "" {
		info.ReleaseURL = release.HTMLURL
	}
	info.ReleaseNotes = truncateNotes(release.Body, 280)
	info.HasUpdate = isNewer(latest, currentVersion)

	return info, nil
}

func truncateNotes(body string, max int) string {
	body = strings.TrimSpace(body)
	if body == "" {
		return ""
	}
	body = strings.ReplaceAll(body, "\r\n", "\n")
	if len(body) <= max {
		return body
	}
	return strings.TrimSpace(body[:max]) + "…"
}
