package logs

import (
	"os"
	"path/filepath"
	"strings"
)

func GetSources(logFiles []string) []SourceInfo {
	sources := make([]SourceInfo, 0, len(logFiles))

	for _, path := range logFiles {
		source := SourceInfo{
			Name: filepath.Base(path),
			Path: path,
		}

		info, err := os.Stat(path)
		if err != nil {
			source.Exists = false
			source.Error = err.Error()
			sources = append(sources, source)
			continue
		}

		source.Exists = true
		source.SizeBytes = info.Size()

		sources = append(sources, source)
	}

	return sources
}

func DetectSource(path string) string {
	name := strings.ToLower(filepath.Base(path))

	if strings.Contains(name, "error") {
		return "error"
	}

	if strings.Contains(name, "out") {
		return "out"
	}

	return filepath.Base(path)
}
