package changelogrelease

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"strings"

	"hermannm.dev/wrap"
)

type CreatedRelease struct {
	Name string
	URL  string
}

func CreateGitHubReleaseForChangelogEntry(
	ctx context.Context,
	input ActionInput,
	httpClient *http.Client,
) (CreatedRelease, error) {
	if err := validateTagName(input.TagName); err != nil {
		return CreatedRelease{}, err
	}

	// If release title is not provided, default to tag name
	releaseTitle := input.ReleaseTitle
	if releaseTitle == "" {
		releaseTitle = input.TagName
	}

	// If changelog file path is not provided, default to CHANGELOG.md (root of repo)
	changelogPath := input.ChangelogFilePath
	if changelogPath == "" {
		changelogPath = "CHANGELOG.md"
	}

	changelog, err := getChangelogEntry(changelogPath, input.TagName)
	if err != nil {
		return CreatedRelease{}, wrap.Error(err, "Failed to get changelog entry")
	}

	githubClient := GitHubApiClient{httpClient: httpClient, apiURL: input.ApiURL}
	release, err := githubClient.createRelease(
		ctx,
		input.TagName,
		releaseTitle,
		changelog,
		input.RepoName,
		input.RepoOwner,
		input.AuthToken,
	)
	if err != nil {
		return CreatedRelease{}, wrap.Error(err, "Failed to create GitHub release")
	}

	return release, nil
}

func validateTagName(tagName string) error {
	version := strings.TrimPrefix(tagName, "v")

	if !semanticVersioningRegex.MatchString(version) {
		return fmt.Errorf(
			"Invalid tag '%s': Expected semantic version format 'vX.Y.Z' (leading 'v' is optional)",
			tagName,
		)
	}

	return nil
}

// Regex:
// - Leading ^ and trailing $, so we always match the full string
// - \d+ to match at least 1 digit
// - \. to match dots between digits
var semanticVersioningRegex = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+$`)

// Implements [slog.LogValuer].
func (release CreatedRelease) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("name", release.Name),
		slog.String("url", release.URL),
	)
}
