package logflusher

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hatchet-dev/hatchet/internal/integrations/filestorage"
	"github.com/hatchet-dev/hatchet/internal/integrations/logstorage"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/repository"
)

type LogFlusherOpts struct {
	LogStore   logstorage.LogStorageBackend
	FileStore  filestorage.FileStorageManager
	Repository repository.Repository
}

type LogFlusher struct {
	ls   logstorage.LogStorageBackend
	fs   filestorage.FileStorageManager
	repo repository.Repository
}

func NewLogFlusher(opts *LogFlusherOpts) *LogFlusher {
	return &LogFlusher{opts.LogStore, opts.FileStore, opts.Repository}
}

func (lf *LogFlusher) Flush(ctx context.Context, input FlushInput) (string, error) {
	run, err := lf.repo.Module().ReadModuleRunByID(input.ModuleID, input.ModuleRunID)

	if err != nil {
		return "", err
	}

	run.TeamID = input.TeamID

	filePath := filestorage.GetModuleRunLogsPath(input.TeamID, run.ModuleID, run.ID)
	lsPath := logstorage.GetLogStoragePath(run.TeamID, run.ModuleID, run.ID)

	logArr := logstorage.NewLogStrArr()

	err = lf.ls.ReadLogs(context.Background(), &logstorage.LogGetOpts{
		Path:  lsPath,
		Count: 0,
	}, logArr)

	if err != nil {
		return "", fmt.Errorf("could not write logs: %s", err.Error())
	}

	fileBytes, err := json.Marshal(logArr.GetLogs())

	if err != nil {
		return "", fmt.Errorf("could not convert logs to json string array: %s", err.Error())
	}

	err = lf.fs.WriteFile(filePath, fileBytes, true)

	if err != nil {
		return "", fmt.Errorf("could not write file: %s", err.Error())
	}

	// update the storage location of the log
	run.LogLocation = models.LogLocationFileStorage

	run, err = lf.repo.Module().UpdateModuleRun(run)

	if err != nil {
		return "", fmt.Errorf("could not update module run: %s", err.Error())
	}

	err = lf.ls.ClearLogs(context.Background(), lsPath)

	if err != nil {
		return "", fmt.Errorf("could not clear logs: %s", err.Error())
	}

	return "success", nil
}
