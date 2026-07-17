package logs

type SourceInfo struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	Exists    bool   `json:"exists"`
	SizeBytes int64  `json:"sizeBytes"`
	Error     string `json:"error,omitempty"`
}

type LogEntry struct {
	Level     string         `json:"level"`
	Time      string         `json:"time"`
	Service   string         `json:"service,omitempty"`
	Env       string         `json:"env,omitempty"`
	Context   string         `json:"context,omitempty"`
	Message   string         `json:"message"`
	RequestID string         `json:"requestId,omitempty"`
	UserID    string         `json:"userId,omitempty"`
	Stack     string         `json:"stack,omitempty"`
	Source    string         `json:"source,omitempty"`
	Raw       map[string]any `json:"raw,omitempty"`
}
