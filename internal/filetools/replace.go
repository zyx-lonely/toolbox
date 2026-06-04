package filetools

import (
	"os"
	"path/filepath"
	"strings"
)

type ReplaceResult struct {
	Path     string `json:"path"`
	Matches  int    `json:"matches"`
	Replaced int    `json:"replaced"`
	Error    string `json:"error,omitempty"`
}

func SearchAndReplace(dir, search, replace, fileTypes string) []ReplaceResult {
	var results []ReplaceResult
	exts := parseExtensions(fileTypes)

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil { return nil }
		if info.IsDir() { return nil }
		if info.Size() > 5*1024*1024 { return nil }
		if len(exts) > 0 {
			ext := strings.ToLower(filepath.Ext(path))
			if !containsExt(exts, ext) { return nil }
		}

		data, err := os.ReadFile(path)
		if err != nil { return nil }
		content := string(data)
		if !strings.Contains(content, search) { return nil }

		count := strings.Count(content, search)
		newContent := strings.ReplaceAll(content, search, replace)
		if err := os.WriteFile(path, []byte(newContent), 0644); err != nil {
			results = append(results, ReplaceResult{Path: path, Matches: count, Error: err.Error()})
			return nil
		}
		results = append(results, ReplaceResult{Path: path, Matches: count, Replaced: count})
		return nil
	})
	return results
}
