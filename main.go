package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"hermannm.dev/devlog"
	"hermannm.dev/devlog/log"
	"hermannm.dev/wrap"

	"hermannm.dev/release-from-changelog/internal/changelogrelease"
)

func main() {
	log.SetDefault(devlog.NewHandler(os.Stdout, nil))

	ctx := context.Background()

	err := releaseFromChangelog(ctx)
	if err != nil {
		log.Error(ctx, err, "")
		os.Exit(1)
	}
}

func releaseFromChangelog(ctx context.Context) error {
	actionInput, err := changelogrelease.ActionInputFromEnv()
	if err != nil {
		return wrap.Error(err, "Failed to parse action input from environment variables")
	}

	release, err := changelogrelease.CreateGitHubReleaseForChangelogEntry(
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
}
