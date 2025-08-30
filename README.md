# release-from-changelog

GitHub Action for creating a GitHub release from a `CHANGELOG.md` entry on the
["Keep a Changelog" format](https://keepachangelog.com/).

**Contents:**

- [Usage](#usage)
- [Why another action?](#why-another-action)
- [Developer's guide](#developers-guide)

## Usage

- First, create a `CHANGELOG.md` file at the root of your repo, following the
  ["Keep a Changelog" format](https://keepachangelog.com/) (see the
  [changelog in this repo](https://github.com/hermannm/release-from-changelog/blob/main/CHANGELOG.md?plain=1)
  as an example)
    - This action expects changelog entry titles to start with `## [vX.Y.Z]` (the leading `v` is
      optional)
- Then, create a GitHub Actions release workflow (in `.github/workflows/release.yml`):

<!-- @formatter:off -->
```yaml
name: Release
on:
  push:
    # Trigger release when a new semantic version tag (vX.Y.Z) is pushed
    tags:
      - v[0-9]+\.[0-9]+\.[0-9]+
env:
  FORCE_COLOR: 1 # For color output in GitHub Actions console (see https://force-color.org/)
jobs:
  create-release:
    name: Create release
    # `release-from-changelog` is a Docker action, which only works on Linux jobs
    runs-on: ubuntu-latest
    permissions:
      # Unfortunately, creating a release requires `contents: write` permission
      # https://github.com/orgs/community/discussions/68252
      contents: write
    steps:
      # https://github.com/actions/checkout
      - name: Checkout repository
        uses: actions/checkout@v5
      # Creates GitHub release from CHANGELOG.md entry for the pushed tag
      # https://github.com/hermannm/release-from-changelog
      - name: Create release from changelog
        uses: hermannm/release-from-changelog@v0.2.5
```
<!-- @formatter:on -->

- When pushing a new tag on the format `vX.Y.Z`, this action will look for a corresponding entry in
  `CHANGELOG.md`, and automatically create a GitHub release for it!

You can customize the action with various inputs:

<!-- @formatter:off -->
```yaml
- name: Create release from changelog
  uses: hermannm/release-from-changelog@v0.2.5
  with:
    # GitHub auth token to use for creating the release.
    # Uses the default `github.token` if not specified.
    # Note that the token needs `contents: write` permission in order to create the release:
    # https://github.com/orgs/community/discussions/68252
    token: <GitHub token>
    # Explicit name of the tag to create a release for.
    # If not specified, we get the tag name from the `GITHUB_REF` set by the workflow's trigger.
    # `GITHUB_REF` must be on the format `refs/tags/<tag_name>`, which it will be if `push: tags` is
    # used as the trigger for the workflow.
    tag_name: <name of existing tag on the format `vX.Y.Z` (leading `v` is optional)>
    # Title of the GitHub release to create.
    # Defaults to the tag name.
    release_title: <title>
    # Path to changelog file in this repo, which we parse to find an entry for the release tag.
    # Defaults to `CHANGELOG.md` at the root of the repo.
    changelog_path: <path/to/changelog.md>
```
<!-- @formatter:on -->

## Why another action?

There are multiple existing actions on GitHub Marketplace that does the same thing as this action,
such as [`taiki-e/create-gh-release-action`](https://github.com/taiki-e/create-gh-release-action).
The motivation for creating my own action was that creating a GitHub release unfortunately requires
`contents: write` permission (see
[this discussion](https://github.com/orgs/community/discussions/68252)). I wanted to automate
releases in my own libraries, but I didn't want to trust third-party actions with this permission.
And so I decided to create my own action, aiming to minimize third-party dependencies, and running
it in Docker instead of downloading a raw binary in the action.

As such, this action is originally designed for my own use. But feel free to use it yourself!

## Developer's guide

When publishing a new release:

- Run tests and linter ([`golangci-lint`](https://golangci-lint.run/)):
  ```
  go test ./... && golangci-lint run
  ```
- Bump the `runs.image` version in `action.yml`
- Bump the version used in the examples under [Usage](#usage) in the README
- Bump `hermannm/release-from-changelog` under `create-release` in `.github/workflows/release.yml`
  to the current latest version of the action
    - We use our own action to create releases, but obviously we can't use the version that we're
      currently publishing
- Add an entry to `CHANGELOG.md` (with the current date)
    - Remember to update the link section, and bump the version for the `[Unreleased]` link
- Create commit and tag for the release (update `TAG` variable in below command):
  ```
  TAG=vX.Y.Z && git commit -m "Release ${TAG}" && git tag -a "${TAG}" -m "Release ${TAG}" && git log --oneline -2
  ```
- Push the commit and tag:
  ```
  git push && git push --tags
  ```
    - Our release workflow will then run the tests, build the Docker image for the action, and
      create a GitHub release with the pushed tag's changelog entry
