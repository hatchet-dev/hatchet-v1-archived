package runmanager

import (
	"fmt"
	"strings"

	"github.com/hatchet-dev/hatchet/internal/models"
)

type RunManager struct {
}

type TriggerInput struct {
	// Files is a list of file paths from a Git integration
	Files []string

	// BaseBranch is the name of the base branch (for a PR)
	BaseBranch string
}

// Trigger returns true if a new module run should be triggered, false if not.
func Trigger(mod *models.Module, kind models.ModuleRunKind, in *TriggerInput) (bool, string) {
	if in.BaseBranch != "" {
		if mod.DeploymentConfig.GitRepoBranch != "" {
			if in.BaseBranch != mod.DeploymentConfig.GitRepoBranch {
				return false, fmt.Sprintf(
					"module deployment branch %s does not match base branch %s",
					mod.DeploymentConfig.GitRepoBranch,
					in.BaseBranch,
				)
			}
		}
	}

	if in.Files != nil {
		didTrigger := false
		targetPath := trimFilePath(mod.DeploymentConfig.ModulePath)

		for _, file := range in.Files {
			if strings.Contains(trimFilePath(file), targetPath) {
				didTrigger = true
			}
		}

		if !didTrigger {
			return false, fmt.Sprintf(
				"module deployment path %s does not match any changed files",
				mod.DeploymentConfig.ModulePath,
			)
		}
	}

	return true, ""
}

func trimFilePath(path string) string {
	return strings.TrimPrefix(strings.TrimPrefix(path, "./"), "/")
}
