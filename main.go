package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"hermannm.dev/devlog"
	"hermannm.dev/devlog/log"
	"hermannm.dev/wrap"
)

func main() {
	runScript(
		func(ctx context.Context) error {
			actionInput, err := actionInputFromEnv()
			if err != nil {
				return wrap.Error(err, "Failed to parse action input from environment variables")
			}

			release, err := createGitHubReleaseForChangelogEntry(
				ctx,
				actionInput,
				http.DefaultClient,
			)
			if err != nil {
				return err
			}

			log.Info(
				ctx,
				fmt.Sprintf("Successfully created release '%s'", release.Name),
				"url", release.URL,
			)

			return nil
		},
	)
}

func createGitHubReleaseForChangelogEntry(
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

	githubClient := GitHubAPIClient{httpClient: httpClient, apiURL: input.ApiURL}
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

type CreatedRelease struct {
	Name string
	URL  string
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
// - Leading ^ and trailing $, so we always match the full string.
// - \d+ to match at least 1 digit.
// - \. to match dots between digits.
var semanticVersioningRegex = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+$`)

func runScript(script func(ctx context.Context) error) {
	log.SetDefault(devlog.NewHandler(os.Stdout, nil))

	ctx := context.Background()

	err := script(ctx)
	if err != nil {
		log.Error(ctx, err, "")
		os.Exit(1)
	}
}
