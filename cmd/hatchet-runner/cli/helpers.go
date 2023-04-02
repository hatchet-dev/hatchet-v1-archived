package cli

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/antihax/optional"
	"github.com/hatchet-dev/hatchet/api/v1/client/swagger"
	"github.com/hatchet-dev/hatchet/internal/config/runner"
	vcs_zip "github.com/hatchet-dev/hatchet/internal/integrations/vcs/zip"
)

func downloadGithubRepoContents(config *runner.Config) error {
	resp, _, err := config.APIClient.ModulesApi.GetModuleTarballURL(
		context.Background(),
		config.ConfigFile.Resources.TeamID,
		config.ConfigFile.Resources.ModuleID,
		&swagger.ModulesApiGetModuleTarballURLOpts{
			GithubSha: optional.NewString(config.ConfigFile.VCS.VCSSHA),
		},
	)

	if err != nil {
		return err
	}

	dstDir := config.ConfigFile.VCS.VCSRepositoryDest

	// make dest directory
	err = os.MkdirAll(dstDir, os.ModePerm)

	if err != nil {
		return err
	}

	zipDownload := &vcs_zip.ZIPDownloader{
		SourceURL:           resp.Url,
		ZipFolderDest:       dstDir,
		ZipName:             fmt.Sprintf("%s.zip", config.ConfigFile.VCS.VCSRepositoryName),
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
			strings.Contains(info.Name(), strings.Replace(config.ConfigFile.VCS.VCSRepositoryName, "/", "-", 1)) &&
			strings.Contains(info.Name(), config.ConfigFile.VCS.VCSSHA) {
			res = filepath.Join(dstDir, info.Name())
		}
	}

	if res == "" {
		return fmt.Errorf("could not find destination folder")
	}

	fullTFPath := filepath.Join(res, config.ConfigFile.VCS.VCSModulePath)
	config.SetTerraformDir(fullTFPath)

	return nil
}
