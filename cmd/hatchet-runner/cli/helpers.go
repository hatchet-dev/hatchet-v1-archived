package cli

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/antihax/optional"
	"github.com/hatchet-dev/hatchet/api/v1/client/swagger"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/runner"
	github_zip "github.com/hatchet-dev/hatchet/internal/integrations/git/github/zip"
	"github.com/hatchet-dev/hatchet/internal/runner/grpcstreamer"
)

func errorHandler(config *runner.Config, description string) error {
	_, _, err := config.APIClient.ModulesApi.FinalizeModuleRun(
		context.Background(),
		swagger.FinalizeModuleRunRequest{
			Status:      string(types.ModuleRunStatusFailed),
			Description: description,
		},
		config.ConfigFile.TeamID,
		config.ConfigFile.ModuleID,
		config.ConfigFile.ModuleRunID,
	)

	return err
}

func successHandler(config *runner.Config, description string) error {
	_, _, err := config.APIClient.ModulesApi.FinalizeModuleRun(
		context.Background(),
		swagger.FinalizeModuleRunRequest{
			Status:      string(types.ModuleRunStatusCompleted),
			Description: description,
		},
		config.ConfigFile.TeamID,
		config.ConfigFile.ModuleID,
		config.ConfigFile.ModuleRunID,
	)

	return err
}

func downloadGithubRepoContents(config *runner.Config) error {
	resp, _, err := config.APIClient.ModulesApi.GetModuleTarballURL(
		context.Background(),
		config.ConfigFile.TeamID,
		config.ConfigFile.ModuleID,
		&swagger.ModulesApiGetModuleTarballURLOpts{
			GithubSha: optional.NewString(config.ConfigFile.GithubSHA),
		},
	)

	if err != nil {
		return err
	}

	dstDir := config.ConfigFile.GithubRepositoryDest

	zipDownload := &github_zip.ZIPDownloader{
		SourceURL:           resp.Url,
		ZipFolderDest:       dstDir,
		ZipName:             fmt.Sprintf("%s.zip", config.ConfigFile.GithubRepositoryName),
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
			strings.Contains(info.Name(), strings.Replace(config.ConfigFile.GithubRepositoryName, "/", "-", 1)) &&
			strings.Contains(info.Name(), config.ConfigFile.GithubSHA) {
			res = filepath.Join(dstDir, info.Name())
		}
	}

	if res == "" {
		return fmt.Errorf("could not find destination folder")
	}

	fullTFPath := filepath.Join(res, config.ConfigFile.GithubModulePath)
	config.SetTerraformDir(fullTFPath)

	return nil
}

func getWriter(config *runner.Config) (io.Writer, error) {
	provClient := config.GRPCClient

	grpcStream, err := grpcstreamer.NewGRPCStreamer(provClient, "")

	if err != nil {
		return nil, err
	}

	return io.MultiWriter(grpcStream, os.Stdout), nil
}
