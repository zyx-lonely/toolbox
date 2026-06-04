package filetools

import (
	"bufio"
	"os"
	"strings"
)

type DiffLine struct {
	Type    string `json:"type"` // "same", "added", "removed", "modified"
	OldLine int    `json:"oldLine,omitempty"`
	NewLine int    `json:"newLine,omitempty"`
	Content string `json:"content"`
}

func DiffFiles(oldPath, newPath string) ([]DiffLine, error) {
	oldLines, err := readLines(oldPath)
	if err != nil { return nil, err }
	newLines, err := readLines(newPath)
	if err != nil { return nil, err }

	maxLen := len(oldLines)
	if len(newLines) > maxLen { maxLen = len(newLines) }

	var diffs []DiffLine
	for i := 0; i < maxLen; i++ {
		switch {
		case i >= len(oldLines):
			diffs = append(diffs, DiffLine{Type: "added", NewLine: i + 1, Content: newLines[i]})
		case i >= len(newLines):
			diffs = append(diffs, DiffLine{Type: "removed", OldLine: i + 1, Content: oldLines[i]})
		case oldLines[i] != newLines[i]:
			diffs = append(diffs, DiffLine{Type: "modified", OldLine: i + 1, NewLine: i + 1, Content: newLines[i]})
		default:
			diffs = append(diffs, DiffLine{Type: "same", OldLine: i + 1, NewLine: i + 1, Content: oldLines[i]})
		}
	}
	return diffs, nil
}

func readLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil { return nil, err }
	defer f.Close()
	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() { lines = append(lines, scanner.Text()) }
	return lines, scanner.Err()
}

var _ = strings.TrimSpace
