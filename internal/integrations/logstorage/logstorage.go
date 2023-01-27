package logstorage

import (
	"context"
	"io"
)

type LogStorageBackend interface {
	StreamLogs(ctx context.Context, path string, w io.WriteCloser) error
	PushLogLine(ctx context.Context, path string, log []byte) error
}
