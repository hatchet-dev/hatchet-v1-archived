package cli

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/antihax/optional"
	"github.com/hatchet-dev/hatchet/api/v1/client/swagger"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/runner"
	github_zip "github.com/hatchet-dev/hatchet/internal/integrations/git/github/zip"
)

func errorHandler(config *runner.Config, reportKind, description string) error {
	_, _, err := config.APIClient.ModulesApi.FinalizeModuleRun(
		context.Background(),
		swagger.FinalizeModuleRunRequest{
			Status:      string(types.ModuleRunStatusFailed),
			Description: description,
			ReportKind:  reportKind,
		},
		config.ConfigFile.Resources.TeamID,
		config.ConfigFile.Resources.ModuleID,
		config.ConfigFile.Resources.ModuleRunID,
	)

	if err != nil {
		return fmt.Errorf("Error reporting error: %s. Original error: %s", err.Error(), description)
	}

	return fmt.Errorf(description)
}

func successHandler(config *runner.Config, reportKind, description string) error {
	_, _, err := config.APIClient.ModulesApi.FinalizeModuleRun(
		context.Background(),
		swagger.FinalizeModuleRunRequest{
			Status:      string(types.ModuleRunStatusCompleted),
			Description: description,
			ReportKind:  reportKind,
		},
		config.ConfigFile.Resources.TeamID,
		config.ConfigFile.Resources.ModuleID,
		config.ConfigFile.Resources.ModuleRunID,
	)

	if err != nil {
		return fmt.Errorf("Error reporting success: %s. Success message: %s", err.Error(), description)
	}

	return nil
}

func downloadGithubRepoContents(config *runner.Config) error {
	resp, _, err := config.APIClient.ModulesApi.GetModuleTarballURL(
		context.Background(),
		config.ConfigFile.Resources.TeamID,
		config.ConfigFile.Resources.ModuleID,
		&swagger.ModulesApiGetModuleTarballURLOpts{
			GithubSha: optional.NewString(config.ConfigFile.Github.GithubSHA),
		},
	)

	if err != nil {
		return err
	}

	dstDir := config.ConfigFile.Github.GithubRepositoryDest

	zipDownload := &github_zip.ZIPDownloader{
		SourceURL:           resp.Url,
		ZipFolderDest:       dstDir,
		ZipName:             fmt.Sprintf("%s.zip", config.ConfigFile.Github.GithubRepositoryName),
		AssetFolderDest:     dstDir,
		RemoveAfterDownload: true,
	}

	err = zipDownload.DownloadToFile()

	if err != nil {
		panic(err)
	}

	err = zipDownload.UnzipToDir()

	if err != nil {
		panic(err)
	}

	dstFiles, err := ioutil.ReadDir(dstDir)
	var res string

	for _, info := range dstFiles {
		if info.Mode().IsDir() &&
			strings.Contains(info.Name(), strings.Replace(config.ConfigFile.Github.GithubRepositoryName, "/", "-", 1)) &&
			strings.Contains(info.Name(), config.ConfigFile.Github.GithubSHA) {
			res = filepath.Join(dstDir, info.Name())
		}
	}

	if res == "" {
		return fmt.Errorf("could not find destination folder")
	}

	fullTFPath := filepath.Join(res, config.ConfigFile.Github.GithubModulePath)
	config.SetTerraformDir(fullTFPath)

	return nil
}
