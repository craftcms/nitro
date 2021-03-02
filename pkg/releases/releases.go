package releases

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// Finder is the interface that is used to find a specific release from the
// GitHub API using a URL, runtime.GOOS, and runtime.GOSARCH. It will return
// a Release struct or an error.
type Finder interface {
	Find(url, system, arch string) (*Release, error)
}

// Downloader takes a URL (to a Github release) and a file path to save the file to.
type Downloader interface {
	Download(url, file string) error
}

// Release represents the release asset
type Release struct {
	URL             string
	ContentType     string
	OperatingSystem string
	Version         string
}

type githubReleases struct {
	URL       string `json:"url"`
	AssetsURL string `json:"assets_url"`
	UploadURL string `json:"upload_url"`
	HTMLURL   string `json:"html_url"`
	ID        int    `json:"id"`
	Author    struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"author"`
	NodeID          string    `json:"node_id"`
	TagName         string    `json:"tag_name"`
	TargetCommitish string    `json:"target_commitish"`
	Name            string    `json:"name"`
	Draft           bool      `json:"draft"`
	Prerelease      bool      `json:"prerelease"`
	CreatedAt       time.Time `json:"created_at"`
	PublishedAt     time.Time `json:"published_at"`
	Assets          []struct {
		URL      string `json:"url"`
		ID       int    `json:"id"`
		NodeID   string `json:"node_id"`
		Name     string `json:"name"`
		Label    string `json:"label"`
		Uploader struct {
			Login             string `json:"login"`
			ID                int    `json:"id"`
			NodeID            string `json:"node_id"`
			AvatarURL         string `json:"avatar_url"`
			GravatarID        string `json:"gravatar_id"`
			URL               string `json:"url"`
			HTMLURL           string `json:"html_url"`
			FollowersURL      string `json:"followers_url"`
			FollowingURL      string `json:"following_url"`
			GistsURL          string `json:"gists_url"`
			StarredURL        string `json:"starred_url"`
			SubscriptionsURL  string `json:"subscriptions_url"`
			OrganizationsURL  string `json:"organizations_url"`
			ReposURL          string `json:"repos_url"`
			EventsURL         string `json:"events_url"`
			ReceivedEventsURL string `json:"received_events_url"`
			Type              string `json:"type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"uploader"`
		ContentType        string    `json:"content_type"`
		State              string    `json:"state"`
		Size               int       `json:"size"`
		DownloadCount      int       `json:"download_count"`
		CreatedAt          time.Time `json:"created_at"`
		UpdatedAt          time.Time `json:"updated_at"`
		BrowserDownloadURL string    `json:"browser_download_url"`
	} `json:"assets"`
	TarballURL string `json:"tarball_url"`
	ZipballURL string `json:"zipball_url"`
	Body       string `json:"body"`
}

type githubReleaseFinder struct {
	HTTPClient *http.Client
}

func (r *githubReleaseFinder) Find(url, system, arch string) (*Release, error) {
	if r.HTTPClient == nil {
		r.HTTPClient = http.DefaultClient
	}

	switch arch {
	case "amd64":
		arch = "x86_64"
	}

	// get the latest release
	resp, err := r.HTTPClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// make sure everything was ok
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code from github: %d", resp.StatusCode)
	}

	if strings.Contains(url, "latest") {
		found := githubReleases{}
		if err := json.NewDecoder(resp.Body).Decode(&found); err != nil {
			return nil, err
		}

		version := found.TagName

		// find the asset from the download
		for _, asset := range found.Assets {
			if strings.Contains(asset.Name, system) && strings.Contains(asset.Name, arch) {
				return &Release{
					URL:             asset.BrowserDownloadURL,
					ContentType:     asset.ContentType,
					OperatingSystem: system,
					Version:         version,
				}, nil
			}
		}

		return nil, fmt.Errorf("unable to find the release")
	}

	releases := []githubReleases{}
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return nil, err
	}

	release := releases[0]

	version := release.TagName

	// find the asset from the download
	for _, asset := range release.Assets {
		if strings.Contains(asset.Name, system) && strings.Contains(asset.Name, arch) {
			return &Release{
				URL:             asset.BrowserDownloadURL,
				ContentType:     asset.ContentType,
				OperatingSystem: system,
				Version:         version,
			}, nil
		}
	}

	return nil, fmt.Errorf("unable to find a release")
}

// NewFinder returns a new github release finder with the default HTTP client.
func NewFinder() Finder {
	return &githubReleaseFinder{
		HTTPClient: http.DefaultClient,
	}
}

type githubReleaseDownloader struct {
	HTTPClient *http.Client
}

func (d *githubReleaseDownloader) Download(url, file string) error {
	resp, err := d.HTTPClient.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non 200 response code")
	}

	f, err := os.OpenFile(file, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func NewDownloader() Downloader {
	return &githubReleaseDownloader{
		HTTPClient: http.DefaultClient,
	}
}
