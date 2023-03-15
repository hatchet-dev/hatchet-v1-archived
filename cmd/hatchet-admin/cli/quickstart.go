package cli

import (
	_ "embed"
	"io/ioutil"

	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	hatchettoken "github.com/hatchet-dev/hatchet/internal/auth/token"
	"github.com/hatchet-dev/hatchet/internal/cliutils"
	"github.com/hatchet-dev/hatchet/internal/config/database"
	"github.com/hatchet-dev/hatchet/internal/config/loader"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/config/shared"
	"github.com/hatchet-dev/hatchet/internal/config/temporal"
	"github.com/hatchet-dev/hatchet/internal/config/worker"
	"github.com/hatchet-dev/hatchet/internal/encryption"
	github_zip "github.com/hatchet-dev/hatchet/internal/integrations/git/github/zip"
	"github.com/hatchet-dev/hatchet/internal/temporal/server/authorizer/token"
	"sigs.k8s.io/yaml"

	"github.com/spf13/cobra"
)

var certDir string
var staticDir string
var generatedConfigDir string

var quickstartCmd = &cobra.Command{
	Use:   "quickstart",
	Short: "Command used to setup a Hatchet instance",
	Run: func(cmd *cobra.Command, args []string) {
		err := runQuickstart()

		if err != nil {
			red := color.New(color.FgRed)
			red.Printf("Error running [%s]:%s\n", cmd.Use, err.Error())
			os.Exit(1)
		}
	},
}

type generatedConfigFiles struct {
	sc     *server.ConfigFile
	shared *shared.ConfigFile
	dc     *database.ConfigFile
	bwc    *worker.BackgroundConfigFile
	tc     *temporal.TemporalConfigFile
}

func init() {
	rootCmd.AddCommand(quickstartCmd)

	quickstartCmd.PersistentFlags().StringVar(
		&certDir,
		"cert-dir",
		"./certs",
		"path to the directory where certificates should be stored",
	)

	quickstartCmd.PersistentFlags().StringVar(
		&generatedConfigDir,
		"generated-config-dir",
		"./generated",
		"path to the directory where the generated config should be written",
	)

	quickstartCmd.PersistentFlags().StringVar(
		&generatedConfigDir,
		"static-dir",
		"./static",
		"path to the directory where the static assets should be served from",
	)
}

func runQuickstart() error {
	generated, err := loadBaseConfigFiles()

	if err != nil {
		return fmt.Errorf("could not get base config files: %w", err)
	}

	err = setupCerts(generated)

	if err != nil {
		return fmt.Errorf("could not setup certs: %w", err)
	}

	err = generateKeys(generated)

	if err != nil {
		return fmt.Errorf("could not generate server keys: %w", err)
	}

	err = generateBearerToken(generated)

	if err != nil {
		return fmt.Errorf("could not generate internal bearer token: %w", err)
	}

	err = writeGeneratedConfig(generated)

	if err != nil {
		return fmt.Errorf("could not write generated config files: %w", err)
	}

	return nil
}

func loadBaseConfigFiles() (*generatedConfigFiles, error) {
	res := &generatedConfigFiles{}

	configFileBytes, err := ioutil.ReadFile(filepath.Join(configDirectory, "shared.yaml"))

	if err != nil {
		return nil, err
	}

	res.shared, err = loader.LoadSharedConfigFile(configFileBytes)

	if err != nil {
		return nil, err
	}

	configFileBytes, err = ioutil.ReadFile(filepath.Join(configDirectory, "database.yaml"))

	if err != nil {
		return nil, err
	}

	res.dc, err = loader.LoadDatabaseConfigFile(configFileBytes)

	if err != nil {
		return nil, err
	}

	configFileBytes, err = ioutil.ReadFile(filepath.Join(configDirectory, "server.yaml"))

	if err != nil {
		return nil, err
	}

	res.sc, err = loader.LoadServerConfigFile(configFileBytes)

	if err != nil {
		return nil, err
	}

	configFileBytes, err = ioutil.ReadFile(filepath.Join(configDirectory, "background_worker.yaml"))

	if err != nil {
		return nil, err
	}

	res.bwc, err = loader.LoadBackgroundWorkerConfigFile(configFileBytes)

	if err != nil {
		return nil, err
	}

	configFileBytes, err = ioutil.ReadFile(filepath.Join(configDirectory, "temporal.yaml"))

	if err != nil {
		return nil, err
	}

	res.tc, err = loader.LoadTemporalConfigFile(configFileBytes)

	if err != nil {
		return nil, err
	}

	return res, nil
}

//go:embed certs/cluster-cert.conf
var ClusterCertConf []byte

//go:embed certs/internal-admin-client-cert.conf
var InternalAdminClientCertConf []byte

//go:embed certs/worker-client-cert.conf
var WorkerClientCertConf []byte

//go:embed certs/generate-certs.sh
var GenerateCertsScript string

func setupCerts(generated *generatedConfigFiles) error {
	color.New(color.FgGreen).Printf("Generating certificates in cert directory %s\n", certDir)

	// verify that bash and openssl are installed on the system
	if !cliutils.CommandExists("openssl") {
		return fmt.Errorf("openssl must be installed and available in your $PATH")
	}

	if !cliutils.CommandExists("bash") {
		return fmt.Errorf("bash must be installed and available in your $PATH")
	}

	cwd, err := os.Getwd()

	if err != nil {
		return err
	}

	// write certificate config files to system
	fullPathCertDir := filepath.Join(cwd, certDir)

	err = os.MkdirAll(fullPathCertDir, os.ModePerm)

	if err != nil {
		return fmt.Errorf("could not create cert directory: %w", err)
	}

	err = os.WriteFile(filepath.Join(fullPathCertDir, "./cluster-cert.conf"), ClusterCertConf, 0666)

	if err != nil {
		return fmt.Errorf("could not create cluster-cert.conf file: %w", err)
	}

	err = os.WriteFile(filepath.Join(fullPathCertDir, "./internal-admin-client-cert.conf"), InternalAdminClientCertConf, 0666)

	if err != nil {
		return fmt.Errorf("could not create internal-admin-client-cert.conf file: %w", err)
	}

	err = os.WriteFile(filepath.Join(fullPathCertDir, "./worker-client-cert.conf"), WorkerClientCertConf, 0666)

	if err != nil {
		return fmt.Errorf("could not create worker-client-cert.conf file: %w", err)
	}

	// run openssl commands
	c := exec.Command("bash", "-s", "-", fullPathCertDir)

	c.Stdin = strings.NewReader(GenerateCertsScript)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	err = c.Run()

	if err != nil {
		return err
	}

	generated.shared.Temporal.Client.TemporalEnabled = true
	generated.shared.Temporal.Client.TemporalClientTLSRootCAFile = filepath.Join(fullPathCertDir, "ca.cert")
	generated.shared.Temporal.Client.TemporalClientTLSCertFile = filepath.Join(fullPathCertDir, "client-internal-admin.pem")
	generated.shared.Temporal.Client.TemporalClientTLSKeyFile = filepath.Join(fullPathCertDir, "client-internal-admin.key")
	generated.shared.Temporal.Client.TemporalTLSServerName = "cluster"

	generated.tc.Frontend.TemporalFrontendTLSRootCAFile = filepath.Join(fullPathCertDir, "ca.cert")
	generated.tc.Frontend.TemporalFrontendTLSCertFile = filepath.Join(fullPathCertDir, "cluster.pem")
	generated.tc.Frontend.TemporalFrontendTLSKeyFile = filepath.Join(fullPathCertDir, "cluster.key")
	generated.tc.Frontend.TemporalFrontendTLSServerName = "cluster"

	generated.tc.Internode.TemporalInternodeTLSRootCAFile = filepath.Join(fullPathCertDir, "ca.cert")
	generated.tc.Internode.TemporalInternodeTLSCertFile = filepath.Join(fullPathCertDir, "cluster.pem")
	generated.tc.Internode.TemporalInternodeTLSKeyFile = filepath.Join(fullPathCertDir, "cluster.key")
	generated.tc.Internode.TemporalInternodeTLSServerName = "cluster"

	generated.tc.Worker.TemporalWorkerTLSRootCAFile = filepath.Join(fullPathCertDir, "ca.cert")
	generated.tc.Worker.TemporalWorkerTLSCertFile = filepath.Join(fullPathCertDir, "cluster.pem")
	generated.tc.Worker.TemporalWorkerTLSKeyFile = filepath.Join(fullPathCertDir, "cluster.key")
	generated.tc.Worker.TemporalWorkerTLSServerName = "cluster"

	generated.tc.UI.TemporalUITLSRootCAFile = filepath.Join(fullPathCertDir, "ca.cert")
	generated.tc.UI.TemporalUITLSCertFile = filepath.Join(fullPathCertDir, "client-internal-admin.pem")
	generated.tc.UI.TemporalUITLSKeyFile = filepath.Join(fullPathCertDir, "client-internal-admin.key")
	generated.tc.UI.TemporalUITLSServerName = "cluster"

	return nil
}

func generateKeys(generated *generatedConfigFiles) error {
	color.New(color.FgGreen).Printf("Generating encryption keys for Hatchet server\n")

	cookieHashKey, err := encryption.GenerateRandomBytes(8)

	if err != nil {
		return fmt.Errorf("could not generate hash key for instance: %w", err)
	}

	cookieBlockKey, err := encryption.GenerateRandomBytes(8)

	if err != nil {
		return fmt.Errorf("could not generate block key for instance: %w", err)
	}

	generated.sc.Auth.Cookie.Secrets = []string{
		cookieHashKey,
		cookieBlockKey,
	}

	generated.bwc.Auth.Cookie.Secrets = []string{
		cookieHashKey,
		cookieBlockKey,
	}

	databaseEncryptionKey, err := encryption.GenerateRandomBytes(16)

	if err != nil {
		return fmt.Errorf("could not generate database encryption key for instance: %w", err)
	}

	generated.dc.EncryptionKey = databaseEncryptionKey

	fileStoreEncryptionKey, err := encryption.GenerateRandomBytes(16)

	if err != nil {
		return fmt.Errorf("could not generate file storage encryption key for instance: %w", err)
	}

	generated.sc.FileStore.Local.FileEncryptionKey = fileStoreEncryptionKey
	generated.bwc.FileStore.Local.FileEncryptionKey = fileStoreEncryptionKey

	temporalInternalSigningKey, err := encryption.GenerateRandomBytes(16)

	if err != nil {
		return fmt.Errorf("could not generate temporal internal signing key for instance: %w", err)
	}

	generated.tc.TemporalInternalSigningKey = temporalInternalSigningKey

	return nil
}

func generateBearerToken(generated *generatedConfigFiles) error {
	color.New(color.FgGreen).Printf("Generating internal bearer token for Hatchet server\n")

	internalAuthConfig := &temporal.InternalAuthConfig{
		InternalNamespace:  generated.tc.TemporalInternalNamespace,
		InternalSigningKey: []byte(generated.tc.TemporalInternalSigningKey),
		InternalTokenOpts: hatchettoken.TokenOpts{
			Issuer:   generated.tc.TemporalBroadcastAddress,
			Audience: []string{generated.tc.TemporalBroadcastAddress},
		},
	}

	token, err := token.GenerateInternalToken(internalAuthConfig)

	if err != nil {
		return err
	}

	generated.shared.Temporal.Client.TemporalBearerToken = token

	return nil
}

func downloadStaticFiles() error {
	color.New(color.FgGreen).Printf("Downloading static files into directory %s\n", staticDir)

	err := os.MkdirAll(generatedConfigDir, os.ModePerm)

	if err != nil {
		return fmt.Errorf("could not create generated config directory: %w", err)
	}

	downloadURL, err := github_zip.GetHatchetStaticAssetsDownloadURL(Version)

	if err != nil {
		return err
	}

	zipDownloader := github_zip.ZIPDownloader{
		SourceURL:           downloadURL,
		ZipFolderDest:       ".",
		ZipName:             "./static.zip",
		AssetFolderDest:     "./static",
		RemoveAfterDownload: true,
	}

	err = zipDownloader.DownloadToFile()

	if err != nil {
		return err
	}

	err = zipDownloader.UnzipToDir()

	if err != nil {
		return err
	}

	return nil
}

func writeGeneratedConfig(generated *generatedConfigFiles) error {
	color.New(color.FgGreen).Printf("Generating config files %s\n", certDir)

	err := os.MkdirAll(generatedConfigDir, os.ModePerm)

	if err != nil {
		return fmt.Errorf("could not create generated config directory: %w", err)
	}

	sharedConfigBytes, err := yaml.Marshal(generated.shared)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join(generatedConfigDir, "./shared.yaml"), sharedConfigBytes, 0666)

	if err != nil {
		return fmt.Errorf("could not write shared.yaml file: %w", err)
	}

	databaseConfigBytes, err := yaml.Marshal(generated.dc)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join(generatedConfigDir, "./database.yaml"), databaseConfigBytes, 0666)

	if err != nil {
		return fmt.Errorf("could not write database.yaml file: %w", err)
	}

	serverConfigBytes, err := yaml.Marshal(generated.sc)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join(generatedConfigDir, "./server.yaml"), serverConfigBytes, 0666)

	if err != nil {
		return fmt.Errorf("could not write server.yaml file: %w", err)
	}

	backgroundWorkerConfigBytes, err := yaml.Marshal(generated.bwc)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join(generatedConfigDir, "./background_worker.yaml"), backgroundWorkerConfigBytes, 0666)

	if err != nil {
		return fmt.Errorf("could not write background_worker.yaml file: %w", err)
	}

	temporalConfigBytes, err := yaml.Marshal(generated.tc)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join(generatedConfigDir, "./temporal.yaml"), temporalConfigBytes, 0666)

	if err != nil {
		return fmt.Errorf("could not write temporal.yaml file: %w", err)
	}

	return nil
}
