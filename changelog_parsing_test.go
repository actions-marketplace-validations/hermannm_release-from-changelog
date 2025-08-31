package main

import (
	"testing"
)

func TestVersionWithLeadingVMatchesChangelogWithoutV(t *testing.T) {
	assertChangelogEntry(
		t,
		"testdata/CHANGELOG_2.md",
		"v0.2.0",
		"- Version without leading 'v'",
	)
}

func TestVersionWithoutLeadingVMatchesChangelogWithV(t *testing.T) {
	assertChangelogEntry(
		t,
		"testdata/CHANGELOG_2.md",
		"v0.3.0",
		"- Test",
	)
}

func TestChangelogEntryAtEndOfFile(t *testing.T) {
	assertChangelogEntry(
		t,
		"testdata/CHANGELOG_2.md",
		"v0.1.0",
		"- Changelog entry at end of file",
	)
}

func TestChangelogEntryAtEndOfFileWithLinks(t *testing.T) {
	assertChangelogEntry(
		t,
		"testdata/CHANGELOG_1.md",
		"v0.1.0",
		"- Initial implementation of the theme for VSCode and IntelliJ",
	)
}

func TestChangelogFromThisProject(t *testing.T) {
	assertChangelogEntry(
		t,
		"CHANGELOG.md",
		"v0.1.0",
		"- Initial implementation of changelog parsing and GitHub release creation",
	)
}

func assertChangelogEntry(
	t *testing.T,
	path string,
	versionToFind string,
	expectedEntry string,
) {
	t.Helper()

	changelogEntry, err := getChangelogEntry(path, versionToFind)
	assertNilError(t, err)

	assertEqual(t, changelogEntry, expectedEntry, "changelog entry")
}
