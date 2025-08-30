package changelogrelease

import (
	"os"
	"sync"
	"testing"
)

func TestActionInputFromEnv(t *testing.T) {
	setTestEnv(
		t,
		map[string]string{
			"INPUT_TAG_NAME":       "v0.4.0",
			"GITHUB_REF":           "refs/tags/v0.3.0", // Should be ignored when INPUT_TAG is set
			"INPUT_RELEASE_TITLE":  "Release 0.4.0",
			"INPUT_CHANGELOG_PATH": "dir/CHANGELOG.md",
			"GITHUB_REPOSITORY":    "hermannm/release-from-changelog",
			"INPUT_TOKEN":          "test-token",
			"GITHUB_API_URL":       "https://api.github.com",
		},
		func() {
			input, err := ActionInputFromEnv()
			assertNilError(t, err)

			expected := ActionInput{
				TagName:           "v0.4.0",
				ReleaseTitle:      "Release 0.4.0",
				ChangelogFilePath: "dir/CHANGELOG.md",
				RepoName:          "release-from-changelog",
				RepoOwner:         "hermannm",
				AuthToken:         "test-token",
				ApiURL:            "https://api.github.com",
			}
			assertDeepEqual(t, input, expected, "action input from env")
		},
	)
}

func TestOptionalInputsAndFallback(t *testing.T) {
	setTestEnv(
		t,
		map[string]string{
			// When INPUT_TAG_NAME is not set, the tag name should be parsed from this env var
			"GITHUB_REF":        "refs/tags/v0.3.0",
			"GITHUB_REPOSITORY": "hermannm/release-from-changelog",
			"INPUT_TOKEN":       "test-token",
			"GITHUB_API_URL":    "https://api.github.com",
		},
		func() {
			input, err := ActionInputFromEnv()
			assertNilError(t, err)

			expected := ActionInput{
				TagName:           "v0.3.0",
				ReleaseTitle:      "",
				ChangelogFilePath: "",
				RepoName:          "release-from-changelog",
				RepoOwner:         "hermannm",
				AuthToken:         "test-token",
				ApiURL:            "https://api.github.com",
			}
			assertDeepEqual(t, input, expected, "action input from env")
		},
	)
}

func setTestEnv(
	t *testing.T,
	envVars map[string]string,
	testFunc func(),
) {
	testEnvLock.Lock()
	defer testEnvLock.Unlock()

	previousValues := make(map[string]string, len(envVars))

	for key, value := range envVars {
		previousValues[key] = os.Getenv(key) // Returns "" if not set

		err := os.Setenv(key, value)
		assertNilError(t, err)
	}

	testFunc()

	for key, previousValue := range previousValues {
		if previousValue != "" {
			err := os.Setenv(key, previousValue)
			assertNilError(t, err)
		} else {
			err := os.Unsetenv(key)
			assertNilError(t, err)
		}
	}
}

var testEnvLock = new(sync.Mutex)
