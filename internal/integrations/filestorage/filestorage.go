package filestorage

import (
	"fmt"
)

var FileDoesNotExist error = fmt.Errorf("the specified file does not exist")

type FileStorageManager interface {
	WriteFile(path string, bytes []byte, shouldEncrypt bool) error
	ReadFile(path string, shouldDecrypt bool) ([]byte, error)
	DeleteFile(path string) error
}

func GetModuleRunLogsPath(teamID, moduleID, runID string) string {
	return fmt.Sprintf("%s/%s/%s/logs.txt", teamID, moduleID, runID)
}

func GetPlanJSONPath(teamID, moduleID, runID string) string {
	return fmt.Sprintf("%s/%s/%s/plan.json", teamID, moduleID, runID)
}

func GetPlanPrettyPath(teamID, moduleID, runID string) string {
	return fmt.Sprintf("%s/%s/%s/plan.txt", teamID, moduleID, runID)
}

func GetPlanZIPPath(teamID, moduleID, runID string) string {
	return fmt.Sprintf("%s/%s/%s/plan.zip", teamID, moduleID, runID)
}
