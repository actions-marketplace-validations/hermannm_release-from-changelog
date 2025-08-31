# Changelog

## [v0.2.5] - 2025-08-31

- Fix cleanup of response body after GitHub request
- Remove dependency on `opt` package

## [v0.2.4] - 2025-08-28

- Improve error log on unsuccessful response from GitHub

## [v0.2.3] - 2025-08-28

- Add `FORCE_COLOR` to example workflow in README, for enabling color output
- Fix accidental duplicate log

## [v0.2.2] - 2025-08-28

- Use latest action version in release workflow, to verify that the new Go implementation works as
  expected

## [v0.2.1] - 2025-08-28

- Fix build step in Dockerfile

## [v0.2.0] - 2025-08-28

- Rewrite action from Rust to Go, to minimize third-party dependencies
    - The original motivation for writing this GitHub action myself was to minimize third-party
      dependencies. But when writing Rust, one has to depend on third-party libraries, since the
      standard library is relatively lightweight. Go, on the other hand, has a batteries-included
      standard library, letting us build the action without any dependencies besides my own utility
      libraries for Go. Not having to keep as many third-party dependencies up to date lessens the
      maintenance burden of this action.

## [v0.1.7] - 2025-08-19

- Improve documentation
- Improve log formatting of GitHub error responses

## [v0.1.6] - 2025-08-19

- Use latest action version in release workflow

## [v0.1.5] - 2025-08-19

- Set `User-Agent` header in requests to GitHub API
- Fix format of `Authorization` header in GitHub requests

## [v0.1.4] - 2025-08-19

- Fix release workflow

## [v0.1.3] - 2025-08-19

- Fix handling of default input for GitHub token

## [v0.1.2] - 2025-08-19

- Fix image version in action metadata

## [v0.1.1] - 2025-08-19

- Use this action in its own GitHub release workflow

## [v0.1.0] - 2025-08-19

- Initial implementation of changelog parsing and GitHub release creation

[Unreleased]: https://github.com/hermannm/release-from-changelog/compare/v0.2.5...HEAD

[v0.2.5]: https://github.com/hermannm/release-from-changelog/compare/v0.2.4...v0.2.5

[v0.2.4]: https://github.com/hermannm/release-from-changelog/compare/v0.2.3...v0.2.4

[v0.2.3]: https://github.com/hermannm/release-from-changelog/compare/v0.2.2...v0.2.3

[v0.2.2]: https://github.com/hermannm/release-from-changelog/compare/v0.2.1...v0.2.2

[v0.2.1]: https://github.com/hermannm/release-from-changelog/compare/v0.2.0...v0.2.1

[v0.2.0]: https://github.com/hermannm/release-from-changelog/compare/v0.1.7...v0.2.0

[v0.1.7]: https://github.com/hermannm/release-from-changelog/compare/v0.1.6...v0.1.7

[v0.1.6]: https://github.com/hermannm/release-from-changelog/compare/v0.1.5...v0.1.6

[v0.1.5]: https://github.com/hermannm/release-from-changelog/compare/v0.1.4...v0.1.5

[v0.1.4]: https://github.com/hermannm/release-from-changelog/compare/v0.1.3...v0.1.4

[v0.1.3]: https://github.com/hermannm/release-from-changelog/compare/v0.1.2...v0.1.3

[v0.1.2]: https://github.com/hermannm/release-from-changelog/compare/v0.1.1...v0.1.2

[v0.1.1]: https://github.com/hermannm/release-from-changelog/compare/v0.1.0...v0.1.1

[v0.1.0]: https://github.com/hermannm/release-from-changelog/compare/ba852f0...v0.1.0
