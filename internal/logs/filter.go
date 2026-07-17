package logs

import (
	"encoding/json"
	"strings"
)

func MatchesFilters(entry LogEntry, query GetLogsQuery) bool {
	if query.Level != "" && !strings.EqualFold(entry.Level, query.Level) {
		return false
	}

	if query.Context != "" && !containsIgnoreCase(entry.Context, query.Context) {
		return false
	}

	if query.RequestID != "" && entry.RequestID != query.RequestID {
		return false
	}

	if query.UserID != "" && entry.UserID != query.UserID {
		return false
	}

	if query.Source != "" && query.Source != "all" && entry.Source != query.Source {
		return false
	}

	if query.Query != "" && !entryContainsQuery(entry, query.Query) {
		return false
	}

	return true
}

func entryContainsQuery(entry LogEntry, query string) bool {
	query = strings.ToLower(strings.TrimSpace(query))
	if query == "" {
		return true
	}

	rawJSON, _ := json.Marshal(entry.Raw)

	text := strings.ToLower(strings.Join([]string{
		entry.Level,
		entry.Time,
		entry.Service,
		entry.Env,
		entry.Context,
		entry.Message,
		entry.RequestID,
		entry.UserID,
		entry.Stack,
		string(rawJSON),
	}, " "))

	return strings.Contains(text, query)
}

func containsIgnoreCase(value string, search string) bool {
	return strings.Contains(
		strings.ToLower(value),
		strings.ToLower(search),
	)
}
