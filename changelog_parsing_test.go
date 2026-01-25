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

func TestChangelogEntryWithLinebreaks(t *testing.T) {
	assertChangelogEntry(
		t,
		"testdata/CHANGELOG_1.md",
		"v0.5.0",
		`- Improve IntelliJ syntax highlighting for:
    - Go
    - Rust
    - TypeScript/JavaScript (more consistent highlighting of function calls, getters, namespaces and imports)
- Update installation instructions`,
	)
}

// Test more edge cases:
// - Long paragraph
// - Numbered lists
// - Comments
// - Code blocks
func TestChangelogEntryWithLinebreaks2(t *testing.T) {
	assertChangelogEntry(
		t,
		"testdata/CHANGELOG_2.md",
		"v0.4.0",
		`Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum lobortis auctor dolor. Phasellus justo neque, molestie ut sodales vel, posuere ut diam. Morbi finibus lacus neque, in efficitur quam commodo id.

1. Aenean laoreet ligula id justo mattis scelerisque eget eget lacus. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos.
2. Curabitur nisl turpis, malesuada ac mattis sodales, pellentesque ac est. Sed vitae placerat leo. Integer congue posuere lorem, at porta velit molestie et. Phasellus nec tempus mi. Fusce et rutrum lacus. Vestibulum pretium rhoncus urna, quis malesuada urna accumsan in.
    - Phasellus dapibus, felis viverra aliquam gravida, lorem ante varius ante, eu finibus ligula ex sit amet mi. Suspendisse placerat velit sem, vel eleifend nulla pulvinar eu.

This is a test.
<!-- Comment -->
Another test.

`+"```"+`js
const example = "test"
console.log(example)
`+"```",
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
