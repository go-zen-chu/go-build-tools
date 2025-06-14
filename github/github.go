package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/go-zen-chu/go-build-tools/util"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type httpClient struct {
	cli *http.Client
}

func NewHTTPClient(cli *http.Client) HTTPClient {
	return &httpClient{
		cli: cli,
	}
}

func (c *httpClient) Do(req *http.Request) (*http.Response, error) {
	return c.cli.Do(req)
}

type GitHubRelease struct {
	Name    string  `json:"name"`
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadUrl string `json:"browser_download_url"`
}

func GetTagRelease(httpClient HTTPClient, owner, repo, gitTag string) (*GitHubRelease, error) {
	releaseUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/tags/%s", owner, repo, gitTag)

	req, err := http.NewRequest(http.MethodGet, releaseUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("new http request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	var release GitHubRelease
	err = json.Unmarshal(body, &release)
	if err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	return &release, nil
}

func GetChecksumMap(httpClient HTTPClient, release *GitHubRelease) (map[string]string, error) {
	checksumMap := make(map[string]string)
	for _, asset := range release.Assets {
		if strings.Contains(asset.Name, "checksum") {
			req, err := http.NewRequest(http.MethodGet, asset.BrowserDownloadUrl, nil)
			if err != nil {
				return nil, fmt.Errorf("new http request: %w", err)
			}
			resp, err := httpClient.Do(req)
			if err != nil {
				return nil, fmt.Errorf("fetching checksums: %w", err)
			}
			defer resp.Body.Close()
			checksum, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("reading checksums: %w", err)
			}
			lines := strings.Split(string(checksum), "\n")
			for _, line := range lines {
				if line == "" {
					continue
				}
				parts := strings.Fields(line)
				if len(parts) != 2 {
					return nil, fmt.Errorf("unexpected checksum line: %s", line)
				}
				// would be filename -> checksum
				checksumMap[parts[1]] = parts[0]
			}
		}
	}
	return checksumMap, nil
}

// GitAddAllCommitPush adds all changes, commits them with the provided message, and pushes to the remote repository.
func GitAddAllCommitPush(repoDir, userName, userEmail, remoteURL, commitMessage string) error {
	// DO NOT use single quote in commands, as it will not work with bash -c
	commands := []string{
		fmt.Sprintf("cd %s", repoDir),
		fmt.Sprintf(`git config user.name "%s"`, userName),
		fmt.Sprintf(`git config user.email "%s"`, userEmail),
		"git add --all",
		fmt.Sprintf(`git commit -m "%s"`, commitMessage),
	}
	if remoteURL != "" {
		commands = append(commands,
			fmt.Sprintf(`git remote set-url origin "%s"`, remoteURL),
		)
	}
	commands = append(commands, "git push origin HEAD")

	bashCommand := fmt.Sprintf(
		"bash -c '%s'",
		strings.Join(commands, " && "),
	)
	outMsg, errMsg, err := util.RunLongRunningCmdWithLog(bashCommand)
	if err != nil {
		return fmt.Errorf("git add, commit and push: %w\nstdout: %s\nstderr: %s", err, outMsg, errMsg)
	}
	return nil
}

func GitHubActionPush(repoDir, owner, repo, commitMessage string) error {
	userName := "GitHub Action (by go-build-tools)"
	userEmail := "aciton@github.com"

	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		return fmt.Errorf("GITHUB_TOKEN is not set")
	}

	remoteUrl := fmt.Sprintf("https://%s@github.com/%s/%s.git", githubToken, owner, repo)
	if err := GitAddAllCommitPush(repoDir, userName, userEmail, remoteUrl, commitMessage); err != nil {
		return fmt.Errorf("github action push: %w", err)
	}
	return nil
}
