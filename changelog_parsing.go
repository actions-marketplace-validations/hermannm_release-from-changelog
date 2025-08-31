package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"hermannm.dev/errclose"
	"hermannm.dev/wrap"
)

func getChangelogEntry(
	changelogPath string,
	versionToFind string,
) (changelogEntry string, returnedErr error) {
	absolutePath, err := filepath.Abs(changelogPath)
	if err != nil {
		return "", wrap.Errorf(
			err, "Failed to get absolute path for changelog file path '%s'", changelogPath,
		)
	}
	file, err := os.Open(absolutePath)
	if err != nil {
		return "", wrap.Errorf(err, "Failed to open changelog file at path '%s'", changelogPath)
	}
	defer errclose.Closef(file, &returnedErr, "changelog file at path '%s'", changelogPath)

	var entryLines []string
	foundEntry := false
	targetTitles := getTargetTitles(versionToFind)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if !foundEntry {
			isTargetTitle := slices.ContainsFunc(
				targetTitles,
				func(targetTitle string) bool {
					return strings.HasPrefix(line, targetTitle)
				},
			)

			if isTargetTitle {
				foundEntry = true

				// Check next line - if it's blank, we don't want to include it in the changelog
				if scanner.Scan() {
					nextLine := scanner.Text()
					if nextLine != "" {
						entryLines = append(entryLines, nextLine)
					}
				}
			}

			continue
		}

		if changelogEntryEnded(line) {
			break
		}

		entryLines = append(entryLines, line)
	}
	if err := scanner.Err(); err != nil {
		return "", wrap.Error(err, "Error while reading changelog file")
	}

	if !foundEntry {
		return "", fmt.Errorf(
			"No changelog entry found for version '%s' in changelog file '%s' (looking for titles starting with one of: %v)",
			versionToFind,
			changelogPath,
			strings.Join(targetTitles, ", "),
		)
	}

	// Remove trailing blank lines from changelog
	for i, line := range slices.Backward(entryLines) {
		if line == "" {
			entryLines = slices.Delete(entryLines, i, i+1)
		}
	}

	if len(entryLines) == 0 {
		return "", fmt.Errorf("Changelog entry for version '%s' was empty", versionToFind)
	}

	return strings.Join(entryLines, "\n"), nil
}

func getTargetTitles(targetVersion string) []string {
	var versionWithPrefix string
	var versionWithoutPrefix string

	if strings.HasPrefix(targetVersion, "v") {
		versionWithPrefix = targetVersion
		versionWithoutPrefix = strings.TrimPrefix(targetVersion, "v")
	} else {
		versionWithoutPrefix = targetVersion
		versionWithPrefix = "v" + targetVersion
	}

	return []string{
		"## [" + versionWithPrefix + "]",
		"## [" + versionWithoutPrefix + "]",
	}
}

// A changelog entry has ended if we find:
// - A higher-level title (#)
// - A new changelog entry at the same title level (##)
// - The start of the link section at the end of the changelog
//   - Example: [v0.1.0]: <link>
func changelogEntryEnded(line string) bool {
	return strings.HasPrefix(line, "# ") ||
		strings.HasPrefix(line, "## ") ||
		tagLinkRegex.MatchString(line)
}

// Regex:
// - Leading ^ to match beginning of line.
// - \[ and \] to match square brackets around link text.
// - [^\[\]]+ to match link text: all characters _except_ [ or ].
// - : to match trailing colon.
var tagLinkRegex = regexp.MustCompile(`^\[[^\[\]]+]:`)
