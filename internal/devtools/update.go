package devtools

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ReleaseInfo GitHub 发布信息
type ReleaseInfo struct {
	TagName     string `json:"tag_name"`
	Name        string `json:"name"`
	Body        string `json:"body"`
	PublishedAt string `json:"published_at"`
	HTMLURL     string `json:"html_url"`
	PreRelease  bool   `json:"prerelease"`
	DownloadURL string `json:"download_url,omitempty"`
}

// CheckUpdate 检查 GitHub Releases 更新
func CheckUpdate(currentVersion string) ReleaseInfo {
	url := "https://api.github.com/repos/user/pc-toolbox/releases/latest"

	client := &http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "pc-toolbox/"+currentVersion)

	resp, err := client.Do(req)
	if err != nil {
		return ReleaseInfo{TagName: currentVersion}
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var release ReleaseInfo
	if err := json.Unmarshal(body, &release); err != nil {
		return ReleaseInfo{TagName: currentVersion}
	}

	return release
}

var _ = fmt.Sprintf
