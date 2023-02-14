package logstorage

import (
	"bytes"
	"context"
	"fmt"
	"io"
)

type LogGetOpts struct {
	Path  string
	Count uint
}

type LogStorageBackend interface {
	// GetID retrieves an ID for the log storage backend
	GetID() string

	// StreamLogs streams logs to the WriteCloser, blocking until the context is cancelled
	StreamLogs(ctx context.Context, opts *LogGetOpts, w io.WriteCloser) error

	// ReadLogs is like StreamLogs, except it returns when all current logs have been read
	ReadLogs(ctx context.Context, opts *LogGetOpts, w io.WriteCloser) error

	// ClearLogs deletes all logs corresponding to the path
	ClearLogs(ctx context.Context, path string) error

	// PushLogLine adds a single log line to the log collection defined by the path
	PushLogLine(ctx context.Context, path string, log []byte) error
}

func GetLogStoragePath(teamID, moduleID, runID string) string {
	return fmt.Sprintf("%s/%s/%s", teamID, moduleID, runID)
}

type bufferCloser struct {
	*bytes.Buffer
}

func NewBufferCloser() *bufferCloser {
	return &bufferCloser{&bytes.Buffer{}}
}

func (bc *bufferCloser) Close() error {
	return nil
}

type logStrArr struct {
	logs []string
}

func NewLogStrArr() *logStrArr {
	return &logStrArr{
		logs: make([]string, 0),
	}
}

func (l *logStrArr) GetLogs() []string {
	return l.logs
}

func (l *logStrArr) Write(p []byte) (n int, err error) {
	l.logs = append(l.logs, string(p))

	return len(p), nil
}

func (l *logStrArr) Close() error {
	return nil
}
