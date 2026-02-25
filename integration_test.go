package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateGitHubReleaseForChangelogEntry(t *testing.T) {
	token := "test-token"
	repoOwner := "hermannm"

	githubServer := mockServer(
		func(res http.ResponseWriter, req *http.Request) {
			expectedRequestBody := `{"tag_name":"v0.4.0","name":"v0.4.0","body":"- Overhaul UI colors in IntelliJ to improve consistency\n- Improve IntelliJ syntax highlighting for:\n    - TypeScript/JavaScript\n    - C/C++\n    - Markdown\n    - XML\n- Improve VSCode syntax highlighting for Rust"}`
			assertRequestBody(t, req, expectedRequestBody)

			assertHeader(t, req, "Authorization", "Bearer "+token)
			assertHeader(t, req, "User-Agent", repoOwner)
			assertHeader(t, req, "Accept", "application/vnd.github+json")
			assertHeader(t, req, "X-GitHub-Api-Version", "2022-11-28")

			res.WriteHeader(http.StatusCreated) // GitHub responds with Created on create release
			_, err := io.WriteString(res, githubCreateReleaseResponse)
			assertNilError(t, err)
		},
	)
	defer githubServer.Close()

	release, err := createGitHubReleaseForChangelogEntry(
		context.Background(),
		ActionInput{
			TagName:           "v0.4.0",
			ReleaseTitle:      "",
			ChangelogFilePath: "testdata/CHANGELOG_1.md",
			RepoName:          "gruvbox-plain",
			RepoOwner:         repoOwner,
			AuthToken:         token,
			ApiURL:            githubServer.URL,
		},
		githubServer.Client(),
	)
	assertNilError(t, err)
	assertEqual(t, release.Name, "v0.4.0", "release name")
	assertEqual(
		t,
		release.URL,
		"https://github.com/hermannm/gruvbox-plain/releases/v0.4.0",
		"release URL",
	)
}

func assertRequestBody(t *testing.T, req *http.Request, expected string) {
	t.Helper()

	requestBody, err := io.ReadAll(req.Body)
	assertNilError(t, err)
	assertEqual(t, string(requestBody), expected, "request body")
}

func assertHeader(t *testing.T, req *http.Request, headerName string, expectedValue string) {
	t.Helper()

	assertEqual(t, req.Header.Get(headerName), expectedValue, headerName+" header")
}

func mockServer(handler func(http.ResponseWriter, *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(handler))
}

// Based on example response from:
// https://docs.github.com/en/rest/releases/releases?apiVersion=2022-11-28#create-a-release
var githubCreateReleaseResponse = `{
  "url": "https://api.github.com/repos/hermannm/gruvbox-plain/releases/4",
  "html_url": "https://github.com/hermannm/gruvbox-plain/releases/v0.4.0",
  "assets_url": "https://api.github.com/repos/hermannm/gruvbox-plain/releases/4/assets",
  "upload_url": "https://uploads.github.com/repos/hermannm/gruvbox-plain/releases/4/assets{?name,label}",
  "tarball_url": "https://api.github.com/repos/hermannm/gruvbox-plain/tarball/v0.4.0",
  "zipball_url": "https://api.github.com/repos/hermannm/gruvbox-plain/zipball/v0.4.0",
  "discussion_url": "https://github.com/hermannm/gruvbox-plain/discussions/90",
  "id": 1,
  "node_id": "MDc6UmVsZWFzZTE=",
  "tag_name": "v0.4.0",
  "target_commitish": "master",
  "name": "v0.4.0",
  "body": "Description of the release",
  "draft": false,
  "prerelease": false,
  "immutable": false,
  "created_at": "2013-02-27T19:35:32Z",
  "published_at": "2013-02-27T19:35:32Z",
  "author": {
    "login": "octocat",
    "id": 1,
    "node_id": "MDQ6VXNlcjE=",
    "avatar_url": "https://github.com/images/error/octocat_happy.gif",
    "gravatar_id": "",
    "url": "https://api.github.com/users/octocat",
    "html_url": "https://github.com/octocat",
    "followers_url": "https://api.github.com/users/octocat/followers",
    "following_url": "https://api.github.com/users/octocat/following{/other_user}",
    "gists_url": "https://api.github.com/users/octocat/gists{/gist_id}",
    "starred_url": "https://api.github.com/users/octocat/starred{/owner}{/repo}",
    "subscriptions_url": "https://api.github.com/users/octocat/subscriptions",
    "organizations_url": "https://api.github.com/users/octocat/orgs",
    "repos_url": "https://api.github.com/users/octocat/repos",
    "events_url": "https://api.github.com/users/octocat/events{/privacy}",
    "received_events_url": "https://api.github.com/users/octocat/received_events",
    "type": "User",
    "site_admin": false
  },
  "assets": [
    {
      "url": "https://api.github.com/repos/hermannm/gruvbox-plain/releases/assets/1",
      "browser_download_url": "https://github.com/hermannm/gruvbox-plain/releases/download/v0.4.0/example.zip",
      "id": 1,
      "node_id": "MDEyOlJlbGVhc2VBc3NldDE=",
      "name": "example.zip",
      "label": "short description",
      "state": "uploaded",
      "content_type": "application/zip",
      "size": 1024,
      "digest": "sha256:2151b604e3429bff440b9fbc03eb3617bc2603cda96c95b9bb05277f9ddba255",
      "download_count": 42,
      "created_at": "2013-02-27T19:35:32Z",
      "updated_at": "2013-02-27T19:35:32Z",
      "uploader": {
        "login": "octocat",
        "id": 1,
        "node_id": "MDQ6VXNlcjE=",
        "avatar_url": "https://github.com/images/error/octocat_happy.gif",
        "gravatar_id": "",
        "url": "https://api.github.com/users/octocat",
        "html_url": "https://github.com/octocat",
        "followers_url": "https://api.github.com/users/octocat/followers",
        "following_url": "https://api.github.com/users/octocat/following{/other_user}",
        "gists_url": "https://api.github.com/users/octocat/gists{/gist_id}",
        "starred_url": "https://api.github.com/users/octocat/starred{/owner}{/repo}",
        "subscriptions_url": "https://api.github.com/users/octocat/subscriptions",
        "organizations_url": "https://api.github.com/users/octocat/orgs",
        "repos_url": "https://api.github.com/users/octocat/repos",
        "events_url": "https://api.github.com/users/octocat/events{/privacy}",
        "received_events_url": "https://api.github.com/users/octocat/received_events",
        "type": "User",
        "site_admin": false
      }
    }
  ]
}`
