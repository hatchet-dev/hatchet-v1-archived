package cli

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/hatchet-dev/hatchet/internal/config/loader"
	"github.com/hatchet-dev/hatchet/internal/config/runner"
	"github.com/hatchet-dev/hatchet/internal/runner/action"
	"github.com/hatchet-dev/hatchet/internal/runner/grpcstreamer"

	"github.com/spf13/cobra"

	github_zip "github.com/hatchet-dev/hatchet/internal/integrations/git/github/zip"
)

var applyCmd = &cobra.Command{
	Use: "apply",
	Run: func(cmd *cobra.Command, args []string) {
		err := runApply()

		if err != nil {
			red := color.New(color.FgRed)
			red.Println("Error running apply:", err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
}

func runApply() error {
	configLoader := &loader.EnvConfigLoader{}
	rc, err := configLoader.LoadRunnerConfigFromEnv()

	if err != nil {
		return err
	}

	err = downloadGithubRepoContents(rc)

	if err != nil {
		return err
	}

	writer, err := getWriter(rc)

	if err != nil {
		return err
	}

	action := action.NewRunnerAction(writer)

	_, err = action.Apply(rc, map[string]interface{}{})

	if err != nil {
		return err
	}

	return nil
}

func downloadGithubRepoContents(config *runner.Config) error {
	resp, _, err := config.APIClient.ModulesApi.GetModuleTarballURL(
		context.Background(),
		config.ConfigFile.TeamID,
		config.ConfigFile.ModuleID,
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
		if info.Mode().IsDir() && strings.Contains(info.Name(), strings.Replace(config.ConfigFile.GithubRepositoryName, "/", "-", 1)) {
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
