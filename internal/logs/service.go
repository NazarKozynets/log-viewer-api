package logs

import (
	"sort"
	"time"
)

type Service struct {
	logFiles     []string
	maxReadBytes int64
	defaultLimit int
	maxLimit     int
}

func NewService(
	logFiles []string,
	maxReadBytes int64,
	defaultLimit int,
	maxLimit int,
) *Service {
	return &Service{
		logFiles:     logFiles,
		maxReadBytes: maxReadBytes,
		defaultLimit: defaultLimit,
		maxLimit:     maxLimit,
	}
}

func (s *Service) GetSources() SourcesResponse {
	return SourcesResponse{
		Sources: GetSources(s.logFiles),
	}
}

func (s *Service) GetLogs(query GetLogsQuery) LogsResponse {
	query.Limit = s.normalizeLimit(query.Limit)

	entries := make([]LogEntry, 0)

	for _, path := range s.logFiles {
		source := DetectSource(path)

		if query.Source != "" && query.Source != "all" && query.Source != source {
			continue
		}

		lines, err := ReadLastLines(path, s.maxReadBytes)
		if err != nil {
			continue
		}

		for _, line := range lines {
			entry, ok := ParseLogLine(line)
			if !ok {
				continue
			}

			entry.Source = source

			if !MatchesFilters(entry, query) {
				continue
			}

			entries = append(entries, entry)
		}
	}

	sort.SliceStable(entries, func(i, j int) bool {
		left := parseLogTime(entries[i].Time)
		right := parseLogTime(entries[j].Time)

		return left.After(right)
	})

	if len(entries) > query.Limit {
		entries = entries[:query.Limit]
	}

	return LogsResponse{
		Items: entries,
		Count: len(entries),
		Limit: query.Limit,
	}
}

func (s *Service) normalizeLimit(limit int) int {
	if limit <= 0 {
		return s.defaultLimit
	}

	if limit > s.maxLimit {
		return s.maxLimit
	}

	return limit
}

func parseLogTime(value string) time.Time {
	if value == "" {
		return time.Time{}
	}

	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		return time.Time{}
	}

	return parsed
}
