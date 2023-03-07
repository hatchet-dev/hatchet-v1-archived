package file

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/hatchet-dev/hatchet/internal/integrations/logstorage"
	"github.com/nxadm/tail"
)

type FileLogStorageManager struct {
	RootDir string
}

func NewFileLogStorageManager(rootDir string) (*FileLogStorageManager, error) {
	err := os.MkdirAll(rootDir, os.ModePerm)

	if err != nil {
		return nil, fmt.Errorf("error creating log directory: %w", err)
	}

	return &FileLogStorageManager{rootDir}, nil
}

func (f *FileLogStorageManager) GetID() string {
	return "filelocal"
}

func (f *FileLogStorageManager) StreamLogs(ctx context.Context, opts *logstorage.LogGetOpts, w io.WriteCloser) error {
	logFilePath := f.getLogFilePath(opts.Path)

	t, err := tail.TailFile(logFilePath, tail.Config{Follow: true, Poll: true})

	defer t.Cleanup()

	if err != nil {
		return fmt.Errorf("error streaming logs: %w", err)
	}

	errorchan := make(chan error)

	go func() {
		select {
		case <-ctx.Done():
			errorchan <- nil
		}
	}()

	go func(t *tail.Tail) {
		for line := range t.Lines {
			_, err = w.Write([]byte(line.Text))

			if err != nil {
				errorchan <- nil
				break
			}
		}
	}(t)

	for err = range errorchan {
		t.Stop()
	}

	return err
}

func (f *FileLogStorageManager) ReadLogs(ctx context.Context, opts *logstorage.LogGetOpts, w io.WriteCloser) error {
	logFilePath := f.getLogFilePath(opts.Path)

	file, err := os.Open(logFilePath)

	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return fmt.Errorf("error reading logs: %w", err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for {
		if !scanner.Scan() {
			return nil
		}

		select {
		case <-ctx.Done():
			return nil
		default:
			lineBytes := scanner.Bytes()
			_, err = w.Write(lineBytes)

			if err != nil {
				return err
			}
		}
	}
}

func (f *FileLogStorageManager) ClearLogs(ctx context.Context, path string) error {
	logFilePath := f.getLogFilePath(path)

	err := os.Remove(logFilePath)

	if err != nil && os.IsNotExist(err) {
		return nil
	}

	return err
}

func (f *FileLogStorageManager) PushLogLine(ctx context.Context, path string, log []byte) error {
	logFilePath := f.getLogFilePath(path)

	file, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)

	if err != nil {
		if os.IsNotExist(err) {
			fileDir := filepath.Dir(logFilePath)

			err = os.MkdirAll(fileDir, os.ModePerm)

			if err != nil {
				return fmt.Errorf("could not create log file: %w", err)
			}
		} else {
			return fmt.Errorf("error opening log file: %w", err)
		}
	}

	defer file.Close()

	fileBytes := log

	fileBytes = append(fileBytes, []byte("\n")...)

	if _, err := file.Write(fileBytes); err != nil {
		return fmt.Errorf("error pushing log: %w", err)
	}

	return nil
}

func (f *FileLogStorageManager) getLogFilePath(path string) string {
	return filepath.Join(f.RootDir, fmt.Sprintf("%s.txt", path))
}
