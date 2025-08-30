package changelogrelease

import (
	"fmt"
	"os"
	"strings"
)

type ActionInput struct {
	TagName           string
	ReleaseTitle      string // Blank if release title was not provided.
	ChangelogFilePath string // Blank if changelog file path was not provided.
	RepoName          string
	RepoOwner         string
	AuthToken         string
	ApiURL            string //nolint:staticcheck // ST1003 - APIURL is ugly
}

func ActionInputFromEnv() (ActionInput, error) {
	tagName, err := getTagNameFromEnv()
	if err != nil {
		return ActionInput{}, err
	}

	repoOwner, repoName, err := getRepoOwnerAndNameFromEnv()
	if err != nil {
		return ActionInput{}, err
	}

	authToken, err := getRequiredEnvVar("INPUT_TOKEN")
	if err != nil {
		return ActionInput{}, err
	}

	apiURL, err := getRequiredEnvVar("GITHUB_API_URL")
	if err != nil {
		return ActionInput{}, err
	}

	return ActionInput{
		TagName:           tagName,
		ReleaseTitle:      os.Getenv("INPUT_RELEASE_TITLE"),  // Returns blank if env var is not set
		ChangelogFilePath: os.Getenv("INPUT_CHANGELOG_PATH"), // Returns blank if env var is not set
		RepoName:          repoName,
		RepoOwner:         repoOwner,
		AuthToken:         authToken,
		ApiURL:            apiURL,
	}, nil
}

func getTagNameFromEnv() (string, error) {
	tagName := os.Getenv("INPUT_TAG_NAME")
	if tagName == "" {
		tagRef, err := getRequiredEnvVar("GITHUB_REF")
		if err != nil {
			return "", err
		}

		if !strings.HasPrefix(tagRef, "refs/tags/") {
			return "", fmt.Errorf(
				"Expected 'GITHUB_REF' environment variable to be on the format 'refs/tags/<tag_name>', but got '%s'",
				tagRef,
			)
		}

		tagName = strings.TrimPrefix(tagRef, "refs/tags/")
	}
	return tagName, nil
}

func getRepoOwnerAndNameFromEnv() (repoOwner string, repoName string, err error) {
	repo, err := getRequiredEnvVar("GITHUB_REPOSITORY")
	if err != nil {
		return "", "", err
	}

	repoSplit := strings.SplitN(repo, "/", 2)
	if len(repoSplit) != 2 {
		return "", "", fmt.Errorf(
			"Expected 'GITHUB_REPOSITORY' environment variable to be on the format 'repo_owner/repo_name', but got '%s'",
			repo,
		)
	}

	return repoSplit[0], repoSplit[1], nil
}

func getRequiredEnvVar(name string) (value string, err error) {
	value = os.Getenv(name)
	if value == "" {
		return "", fmt.Errorf("Expected '%s' environment variable to be set", name)
	}
	return value, nil
}
