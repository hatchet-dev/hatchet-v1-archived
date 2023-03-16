package github_zip

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/google/go-github/v49/github"
)

type ZIPDownloader struct {
	SourceURL string

	ZipFolderDest   string
	ZipName         string
	AssetFolderDest string

	RemoveAfterDownload bool
}

func (z *ZIPDownloader) DownloadToFile() error {
	// Get the data
	resp, err := http.Get(z.SourceURL)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath.Join(z.ZipFolderDest, z.ZipName))

	if err != nil {
		return err
	}

	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)

	return err
}

func (z *ZIPDownloader) UnzipToDir() error {
	r, err := zip.OpenReader(filepath.Join(z.ZipFolderDest, z.ZipName))

	if err != nil {
		return err
	}

	defer r.Close()

	for _, f := range r.File {
		// Store filename/path for returning and using later on
		fpath := filepath.Join(z.AssetFolderDest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(z.AssetFolderDest)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: illegal file path", fpath)
		}

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// delete file if exists
		os.Remove(fpath)

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}

	if z.RemoveAfterDownload {
		os.Remove(filepath.Join(z.ZipFolderDest, z.ZipName))
	}

	return nil
}

const HatchetRepositoryOwner = "hatchet-dev"
const HatchetRepositoryName = "hatchet"
const HatchetStaticAssetsName = "hatchet-static"

func GetHatchetStaticAssetsDownloadURL(version string) (string, error) {
	return getHatchetAssetDownloadURL(version, HatchetStaticAssetsName, false)
}

func getHatchetAssetDownloadURL(releaseTag, assetName string, isPlatformDependent bool) (string, error) {
	client := github.NewClient(nil)

	rel, _, err := client.Repositories.GetReleaseByTag(
		context.Background(),
		HatchetRepositoryOwner,
		HatchetRepositoryName,
		releaseTag,
	)

	if err != nil {
		return "", fmt.Errorf("release %s does not exist: %w", releaseTag, err)
	}

	re := getDownloadRegexp(assetName, isPlatformDependent)
	releaseURL := ""

	// iterate through the assets
	for _, asset := range rel.Assets {
		downloadURL := asset.GetBrowserDownloadURL()

		if re.MatchString(asset.GetName()) {
			releaseURL = downloadURL
			break
		}
	}

	return releaseURL, nil
}

func getDownloadRegexp(assetName string, isPlatformDependent bool) *regexp.Regexp {
	if !isPlatformDependent {
		return regexp.MustCompile(fmt.Sprintf(`%s_.*\.zip`, assetName))
	}

	switch os := runtime.GOOS; os {
	case "darwin":
		return regexp.MustCompile(fmt.Sprintf(`%s_.*_Darwin_x86_64\.zip`, assetName))
	case "linux":
		return regexp.MustCompile(fmt.Sprintf(`%s_.*_Linux_x86_64\.zip`, assetName))
	default:
		return regexp.MustCompile(fmt.Sprintf(`%s_.*\.zip`, assetName))
	}
}
