package logs

type GetLogsQuery struct {
	Level     string
	Context   string
	RequestID string
	UserID    string
	Query     string
	Source    string
	Limit     int
}

type LogsResponse struct {
	Items []LogEntry `json:"items"`
	Count int        `json:"count"`
	Limit int        `json:"limit"`
}

type SourcesResponse struct {
	Sources []SourceInfo `json:"sources"`
}
