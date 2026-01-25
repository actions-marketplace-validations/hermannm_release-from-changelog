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
		} else {
			break
		}
	}

	if len(entryLines) == 0 {
		return "", fmt.Errorf("Changelog entry for version '%s' was empty", versionToFind)
	}

	entryLines = removeParagraphLinebreaks(entryLines)

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

// We typically format changelog files with a max line length, with line breaks to break up long
// paragraphs. In most Markdown renderers, such single linebreaks in the middle of a paragraph will
// not show up when rendered, as the paragraph will just continue on. But GitHub's rendering of
// release descriptions actually shows these linebreaks, which is annoying. So we strip away single
// line breaks from changelog descriptions (but keep double linebreaks, linebreaks between list
// items and in code blocks, and around comments).
func removeParagraphLinebreaks(lines []string) []string {
	newLines := make([]string, 0, len(lines))

	isInCodeBlock := false
	for _, line := range lines {
		if len(newLines) == 0 {
			newLines = append(newLines, line)
			continue
		}

		previousLineIndex := len(newLines) - 1
		previousLine := newLines[previousLineIndex]

		firstWord, indentLength, ok := getFirstWord(line)
		if !ok {
			newLines = append(newLines, line)
			continue
		}

		// Don't remove line breaks between list items, or around a comment
		if firstWord == "-" ||
			numberedListRegex.MatchString(firstWord) ||
			firstWord == "<!--" ||
			strings.HasSuffix(previousLine, "-->") {
			newLines = append(newLines, line)
			continue
		}

		if strings.HasPrefix(firstWord, "```") {
			isInCodeBlock = !isInCodeBlock
			newLines = append(newLines, line)
			continue
		}
		// If we're in a code block, then we don't want to remove line breaks
		if isInCodeBlock {
			newLines = append(newLines, line)
			continue
		}

		if previousLine == "" {
			newLines = append(newLines, line)
			continue
		}

		newLines[previousLineIndex] = previousLine + " " + line[indentLength:]
	}

	return newLines
}

func getFirstWord(line string) (word string, indentLength int, ok bool) {
	chars := []rune(line)

	var startIndex, endIndex int
	var startFound, endFound bool
	for i, char := range chars {
		if !startFound && char != ' ' {
			startIndex = i
			startFound = true
			continue
		}

		if startFound && char == ' ' {
			endIndex = i
			endFound = true
			break
		}
	}

	if !startFound {
		return "", 0, false
	}
	if !endFound {
		endIndex = len(chars)
	}

	return string(chars[startIndex:endIndex]), startIndex, true
}

// Matches Markdown numbered list items ("1.", "2.", "10." etc.).
var numberedListRegex = regexp.MustCompile(`\d+\.`)
