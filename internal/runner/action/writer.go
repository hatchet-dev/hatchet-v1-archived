package action

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hatchet-dev/hatchet/internal/config/runner"
	"github.com/hatchet-dev/hatchet/internal/runner/grpcstreamer"

	runnertypes "github.com/hatchet-dev/hatchet/internal/runner/types"
)

func GetWriters(config *runner.Config) (io.Writer, io.Writer, error) {
	grpcStream, err := grpcstreamer.NewGRPCStreamer(config.GRPCClient)

	if err != nil {
		return nil, nil, err
	}

	return &PrettyWriter{
			file:             os.Stdout,
			additionalWriter: grpcStream,
		}, &PrettyWriter{
			file:             os.Stderr,
			additionalWriter: grpcStream,
		}, nil
}

type PrettyWriter struct {
	file *os.File

	additionalWriter io.Writer
}

func (p *PrettyWriter) Write(log []byte) (int, error) {
	if p.additionalWriter != nil {
		p.additionalWriter.Write(log)
	}

	return write(p.file, log)
}

func write(o *os.File, log []byte) (int, error) {
	if len(log) == 0 {
		return 0, nil
	}

	logLines := splitWithEscaping(string(log), "\n", "\\")

	for _, logLine := range logLines {
		rawBytes := []byte(logLine)

		rawBytes = append(rawBytes, []byte("\n")...)

		if len(logLine) == 0 {
			continue
		}

		tfLog := &runnertypes.TFLogLine{}

		err := json.Unmarshal([]byte(logLine), tfLog)

		if err != nil {
			o.Write(rawBytes)
			continue
		}

		prettyBytes := getPrettyBytes(tfLog)

		if len(prettyBytes) == 0 {
			continue
		}

		_, err = o.Write(prettyBytes)

		if err != nil {
			o.Write(rawBytes)
			continue
		}
	}

	return len(log), nil
}

func splitWithEscaping(s, separator, escape string) []string {
	s = strings.ReplaceAll(s, escape+separator, "\x00")

	tokens := strings.Split(s, separator)

	for i, token := range tokens {
		tokens[i] = strings.ReplaceAll(token, "\x00", separator)
	}

	return tokens
}

func getPrettyBytes(tfLog *runnertypes.TFLogLine) []byte {
	if tfLog.Level == "debug" || tfLog.Level == "trace" {
		return []byte{}
	}

	formattedLog := fmt.Sprintf("[%s] [%s] %s\n", tfLog.Level, tfLog.Timestamp, tfLog.Message)

	return []byte(formattedLog)
}
