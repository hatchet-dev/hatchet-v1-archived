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
	vcs_zip "github.com/hatchet-dev/hatchet/internal/integrations/vcs/zip"
	"github.com/hatchet-dev/hatchet/internal/temporal/server/authorizer/token"
	"sigs.k8s.io/yaml"

	"github.com/spf13/cobra"
)

var certDir string
var staticDir string
var generatedConfigDir string
var skip []string
var overwrite bool

const (
	StageCerts  string = "certs"
	StageKeys   string = "keys"
	StageTokens string = "tokens"
	StageStatic string = "static"
)

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
	rwc    *worker.RunnerConfigFile
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
		&staticDir,
		"static-dir",
		"./static",
		"path to the directory where the static assets should be served from",
	)

	quickstartCmd.PersistentFlags().StringArrayVar(
		&skip,
		"skip",
		[]string{},
		"a list of steps to skip. possible values are \"certs\", \"static\", \"keys\", or \"tokens\"",
	)

	quickstartCmd.PersistentFlags().BoolVar(
		&overwrite,
		"overwrite",
		true,
		"whether generated files should be overwritten, if they exist",
	)
}

func runQuickstart() error {
	generated, err := loadBaseConfigFiles()

	if err != nil {
		return fmt.Errorf("could not get base config files: %w", err)
	}

	if !shouldSkip(StageCerts) {
		err = setupCerts(generated)

		if err != nil {
			return fmt.Errorf("could not setup certs: %w", err)
		}
	}

	if !shouldSkip(StageKeys) {
		err = generateKeys(generated)

		if err != nil {
			return fmt.Errorf("could not generate server keys: %w", err)
		}
	}

	if !shouldSkip(StageTokens) {
		err = generateBearerToken(generated)

		if err != nil {
			return fmt.Errorf("could not generate internal bearer token: %w", err)
		}
	}

	if !shouldSkip(StageStatic) {
		err = downloadStaticFiles(generated)

		if err != nil {
			return fmt.Errorf("could not download static files: %w", err)
		}
	}

	err = writeGeneratedConfig(generated)

	if err != nil {
		return fmt.Errorf("could not write generated config files: %w", err)
	}

	return nil
}

func shouldSkip(stage string) bool {
	for _, skipStage := range skip {
		if stage == skipStage {
			return true
		}
	}

	return false
}

func loadBaseConfigFiles() (*generatedConfigFiles, error) {
	res := &generatedConfigFiles{}
	var err error

	res.shared, err = loader.LoadSharedConfigFile(getFiles("shared.yaml")...)

	if err != nil {
		return nil, err
	}

	res.dc, err = loader.LoadDatabaseConfigFile(getFiles("database.yaml")...)

	if err != nil {
		return nil, err
	}

	res.sc, err = loader.LoadServerConfigFile(getFiles("server.yaml")...)

	if err != nil {
		return nil, err
	}

	res.bwc, err = loader.LoadBackgroundWorkerConfigFile(getFiles("background_worker.yaml")...)

	if err != nil {
		return nil, err
	}

	res.rwc, err = loader.LoadRunnerWorkerConfigFile(getFiles("runner_worker.yaml")...)

	if err != nil {
		return nil, err
	}

	res.tc, err = loader.LoadTemporalConfigFile(getFiles("temporal.yaml")...)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func getFiles(name string) [][]byte {
	files := [][]byte{}

	basePath := filepath.Join(configDirectory, name)

	if fileExists(basePath) {
		configFileBytes, err := ioutil.ReadFile(basePath)

		if err != nil {
			panic(err)
		}

		files = append(files, configFileBytes)
	}

	generatedPath := filepath.Join(generatedConfigDir, name)

	if fileExists(generatedPath) {
		generatedFileBytes, err := ioutil.ReadFile(filepath.Join(generatedConfigDir, name))

		if err != nil {
			panic(err)
		}

		files = append(files, generatedFileBytes)
	}

	return files
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

	// if we determine we should write the root CA file, we regenerate ALL certificates
	if shouldWriteConfig(generated.shared.Temporal.Client.TemporalClientTLSRootCAFile) {
		generated.shared.Temporal.Client.TemporalEnabled = true

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
	} else {
		fmt.Println("skipping certificate generation because root CA is already set")
	}

	return nil
}

func shouldWriteConfig(conf string) bool {
	return overwrite || conf == ""
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

	if overwrite || (generated.sc.Auth.Cookie.Secrets == nil || len(generated.sc.Auth.Cookie.Secrets) == 0) {
		generated.sc.Auth.Cookie.Secrets = []string{
			cookieHashKey,
			cookieBlockKey,
		}
	}

	if overwrite || (generated.bwc.Auth.Cookie.Secrets == nil || len(generated.bwc.Auth.Cookie.Secrets) == 0) {
		generated.bwc.Auth.Cookie.Secrets = []string{
			cookieHashKey,
			cookieBlockKey,
		}
	}

	databaseEncryptionKey, err := encryption.GenerateRandomBytes(16)

	if err != nil {
		return fmt.Errorf("could not generate database encryption key for instance: %w", err)
	}

	if shouldWriteConfig(generated.dc.EncryptionKey) {
		generated.dc.EncryptionKey = databaseEncryptionKey
	}

	fileStoreEncryptionKey, err := encryption.GenerateRandomBytes(16)

	if err != nil {
		return fmt.Errorf("could not generate file storage encryption key for instance: %w", err)
	}

	if shouldWriteConfig(generated.sc.FileStore.Local.FileEncryptionKey) {
		generated.sc.FileStore.Local.FileEncryptionKey = fileStoreEncryptionKey
	}

	if shouldWriteConfig(generated.bwc.FileStore.Local.FileEncryptionKey) {
		generated.bwc.FileStore.Local.FileEncryptionKey = fileStoreEncryptionKey
	}

	temporalInternalSigningKey, err := encryption.GenerateRandomBytes(16)

	if err != nil {
		return fmt.Errorf("could not generate temporal internal signing key for instance: %w", err)
	}

	if shouldWriteConfig(generated.tc.TemporalInternalSigningKey) {
		generated.tc.TemporalInternalSigningKey = temporalInternalSigningKey
	}

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

	if shouldWriteConfig(generated.shared.Temporal.Client.TemporalBearerToken) {
		generated.shared.Temporal.Client.TemporalBearerToken = token
	}

	return nil
}

func downloadStaticFiles(generated *generatedConfigFiles) error {
	if shouldWriteConfig(generated.sc.Runtime.StaticFileServerPath) {
		color.New(color.FgGreen).Printf("Downloading static files into directory %s\n", staticDir)
		cwd, err := os.Getwd()

		if err != nil {
			return err
		}

		err = os.MkdirAll(generatedConfigDir, os.ModePerm)

		if err != nil {
			return fmt.Errorf("could not create generated config directory: %w", err)
		}

		downloadURL, err := vcs_zip.GetHatchetStaticAssetsDownloadURL(Version)

		if err != nil {
			return err
		}

		zipDownloader := vcs_zip.ZIPDownloader{
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

		generated.sc.Runtime.RunStaticFileServer = true

		generated.sc.Runtime.StaticFileServerPath = filepath.Join(cwd, staticDir)
	} else {

	}

	return nil
}

func writeGeneratedConfig(generated *generatedConfigFiles) error {
	color.New(color.FgGreen).Printf("Generating config files %s\n", certDir)

	err := os.MkdirAll(generatedConfigDir, os.ModePerm)

	if err != nil {
		return fmt.Errorf("could not create generated config directory: %w", err)
	}

	sharedPath := filepath.Join(generatedConfigDir, "./shared.yaml")
	sharedConfigBytes, err := yaml.Marshal(generated.shared)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(sharedPath, sharedConfigBytes, 0666)

	if err != nil {
		return fmt.Errorf("could not write shared.yaml file: %w", err)
	}

	databasePath := filepath.Join(generatedConfigDir, "./database.yaml")

	databaseConfigBytes, err := yaml.Marshal(generated.dc)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(databasePath, databaseConfigBytes, 0666)

	if err != nil {
		return fmt.Errorf("could not write database.yaml file: %w", err)
	}

	serverPath := filepath.Join(generatedConfigDir, "./server.yaml")

	serverConfigBytes, err := yaml.Marshal(generated.sc)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(serverPath, serverConfigBytes, 0666)

	if err != nil {
		return fmt.Errorf("could not write server.yaml file: %w", err)
	}

	backgroundWorkerPath := filepath.Join(generatedConfigDir, "./background_worker.yaml")

	backgroundWorkerConfigBytes, err := yaml.Marshal(generated.bwc)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(backgroundWorkerPath, backgroundWorkerConfigBytes, 0666)

	if err != nil {
		return fmt.Errorf("could not write background_worker.yaml file: %w", err)
	}

	runnerWorkerPath := filepath.Join(generatedConfigDir, "./runner_worker.yaml")

	runnerWorkerConfigBytes, err := yaml.Marshal(generated.rwc)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(runnerWorkerPath, runnerWorkerConfigBytes, 0666)

	if err != nil {
		return fmt.Errorf("could not write runner_worker.yaml file: %w", err)
	}

	temporalPath := filepath.Join(generatedConfigDir, "./temporal.yaml")

	temporalConfigBytes, err := yaml.Marshal(generated.tc)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(temporalPath, temporalConfigBytes, 0666)

	if err != nil {
		return fmt.Errorf("could not write temporal.yaml file: %w", err)
	}

	return nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
