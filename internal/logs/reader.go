package logs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

func ReadLastLines(path string, maxReadBytes int64) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	size := info.Size()

	start := int64(0)
	if maxReadBytes > 0 && size > maxReadBytes {
		start = size - maxReadBytes
	}

	if _, err := file.Seek(start, io.SeekStart); err != nil {
		return nil, err
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if start > 0 {
		if index := bytes.IndexByte(data, '\n'); index >= 0 {
			data = data[index+1:]
		}
	}

	content := strings.ReplaceAll(string(data), "\r\n", "\n")
	rawLines := strings.Split(content, "\n")

	lines := make([]string, 0, len(rawLines))

	for _, line := range rawLines {
		line = strings.TrimSpace(line)
		if line != "" {
			lines = append(lines, line)
		}
	}

	return lines, nil
}

func ParseLogLine(line string) (LogEntry, bool) {
	var raw map[string]any

	if err := json.Unmarshal([]byte(line), &raw); err != nil {
		return LogEntry{}, false
	}

	entry := LogEntry{
		Raw: raw,
	}

	entry.Level = getLevel(raw)
	entry.Time = getTime(raw)
	entry.Service = getString(raw, "service")
	entry.Env = getString(raw, "env")
	entry.Context = getString(raw, "context")
	entry.Message = getString(raw, "message")
	entry.RequestID = getString(raw, "requestId")
	entry.UserID = getString(raw, "userId")
	entry.Stack = getString(raw, "stack")

	if entry.Message == "" {
		entry.Message = getString(raw, "msg")
	}

	if entry.RequestID == "" {
		entry.RequestID = getString(raw, "reqId")
	}

	return entry, true
}

func getString(raw map[string]any, key string) string {
	value, ok := raw[key]
	if !ok || value == nil {
		return ""
	}

	switch typed := value.(type) {
	case string:
		return typed
	case float64:
		return strconv.FormatFloat(typed, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(typed)
	default:
		return fmt.Sprintf("%v", typed)
	}
}

func getLevel(raw map[string]any) string {
	value, ok := raw["level"]
	if !ok || value == nil {
		return ""
	}

	switch typed := value.(type) {
	case string:
		return strings.ToLower(typed)
	case float64:
		switch int(typed) {
		case 10, 20:
			return "debug"
		case 30:
			return "info"
		case 40:
			return "warn"
		case 50:
			return "error"
		case 60:
			return "fatal"
		default:
			return strconv.Itoa(int(typed))
		}
	default:
		return fmt.Sprintf("%v", typed)
	}
}

func getTime(raw map[string]any) string {
	value, ok := raw["time"]
	if !ok || value == nil {
		return ""
	}

	switch typed := value.(type) {
	case string:
		return typed
	case float64:
		// Pino часто пишет time как Unix ms.
		t := time.UnixMilli(int64(typed))
		return t.UTC().Format(time.RFC3339Nano)
	default:
		return fmt.Sprintf("%v", typed)
	}
}
