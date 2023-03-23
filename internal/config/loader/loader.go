package loader

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/creasty/defaults"
	"github.com/spf13/viper"

	"github.com/hatchet-dev/hatchet/api/serverutils/erroralerter"
	"github.com/hatchet-dev/hatchet/api/v1/client/fileclient"
	"github.com/hatchet-dev/hatchet/api/v1/client/grpc"
	"github.com/hatchet-dev/hatchet/api/v1/client/swagger"
	"github.com/hatchet-dev/hatchet/internal/adapter"
	"github.com/hatchet-dev/hatchet/internal/auth/cookie"
	"github.com/hatchet-dev/hatchet/internal/auth/token"
	"github.com/hatchet-dev/hatchet/internal/config/cli"
	"github.com/hatchet-dev/hatchet/internal/config/database"
	"github.com/hatchet-dev/hatchet/internal/config/runner"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/config/shared"
	temporalconfig "github.com/hatchet-dev/hatchet/internal/config/temporal"
	"github.com/hatchet-dev/hatchet/internal/config/worker"
	"github.com/hatchet-dev/hatchet/internal/integrations/filestorage"
	filelocal "github.com/hatchet-dev/hatchet/internal/integrations/filestorage/local"
	"github.com/hatchet-dev/hatchet/internal/integrations/filestorage/s3"
	"github.com/hatchet-dev/hatchet/internal/integrations/logstorage"
	"github.com/hatchet-dev/hatchet/internal/integrations/logstorage/file"
	"github.com/hatchet-dev/hatchet/internal/integrations/logstorage/redis"
	"github.com/hatchet-dev/hatchet/internal/integrations/oauth"
	"github.com/hatchet-dev/hatchet/internal/integrations/oauth/github"
	"github.com/hatchet-dev/hatchet/internal/logger"
	"github.com/hatchet-dev/hatchet/internal/notifier"
	"github.com/hatchet-dev/hatchet/internal/notifier/noop"
	"github.com/hatchet-dev/hatchet/internal/notifier/sendgrid"
	"github.com/hatchet-dev/hatchet/internal/provisioner"
	"github.com/hatchet-dev/hatchet/internal/provisioner/local"
	"github.com/hatchet-dev/hatchet/internal/queuemanager"
	"github.com/hatchet-dev/hatchet/internal/repository/gorm"
	"github.com/hatchet-dev/hatchet/internal/temporal"
)

// LoadSharedConfigFile loads the shared config file via viper
func LoadSharedConfigFile(files ...[]byte) (*shared.ConfigFile, error) {
	configFile := &shared.ConfigFile{}
	f := shared.BindAllEnv

	_, err := loadConfigFromViper(f, configFile, files...)

	return configFile, err
}

// LoadServerConfigFile loads the server config file via viper
func LoadServerConfigFile(files ...[]byte) (*server.ConfigFile, error) {
	configFile := &server.ConfigFile{}
	f := server.BindAllEnv

	_, err := loadConfigFromViper(f, configFile, files...)

	return configFile, err
}

// LoadRunnerConfigFile loads the runner config file via viper
func LoadRunnerConfigFile(files ...[]byte) (*runner.ConfigFile, error) {
	configFile := &runner.ConfigFile{}
	f := runner.BindAllEnv

	_, err := loadConfigFromViper(f, configFile, files...)

	return configFile, err
}

// LoadDatabaseConfigFile loads the database config file via viper
func LoadDatabaseConfigFile(files ...[]byte) (*database.ConfigFile, error) {
	configFile := &database.ConfigFile{}
	f := database.BindAllEnv

	_, err := loadConfigFromViper(f, configFile, files...)

	return configFile, err
}

// LoadBackgroundWorkerConfigFile loads the background worker config file via viper
func LoadBackgroundWorkerConfigFile(files ...[]byte) (*worker.BackgroundConfigFile, error) {
	configFile := &worker.BackgroundConfigFile{}
	f := worker.BindAllBackgroundEnv

	_, err := loadConfigFromViper(f, configFile, files...)

	return configFile, err
}

// LoadRunnerWorkerConfigFile loads the runner worker config file via viper
func LoadRunnerWorkerConfigFile(files ...[]byte) (*worker.RunnerConfigFile, error) {
	configFile := &worker.RunnerConfigFile{}
	f := worker.BindAllRunnerEnv

	_, err := loadConfigFromViper(f, configFile, files...)

	return configFile, err
}

// LoadTemporalConfigFile loads the temporal config file via viper
func LoadTemporalConfigFile(files ...[]byte) (*temporalconfig.TemporalConfigFile, error) {
	configFile := &temporalconfig.TemporalConfigFile{}
	f := temporalconfig.BindAllEnv

	_, err := loadConfigFromViper(f, configFile, files...)

	return configFile, err
}

// LoadCLIConfigFile loads the CLI config file via viper
func LoadCLIConfigFile(files ...[]byte) (*cli.ConfigFile, *viper.Viper, error) {
	configFile := &cli.ConfigFile{}
	f := cli.BindAllEnv

	v, err := loadConfigFromViper(f, configFile, files...)

	return configFile, v, err
}

func loadConfigFromViper(bindFunc func(v *viper.Viper), configFile interface{}, files ...[]byte) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigType("yaml")
	bindFunc(v)

	for _, f := range files {
		err := v.MergeConfig(bytes.NewBuffer(f))

		if err != nil {
			return nil, fmt.Errorf("could not load viper config: %w", err)
		}
	}

	defaults.Set(configFile)

	err := v.Unmarshal(configFile)

	if err != nil {
		return nil, fmt.Errorf("could not unmarshal viper config: %w", err)
	}

	return v, nil
}

type ConfigLoader struct {
	version, directory string
}

func NewConfigLoader(version, directory string) *ConfigLoader {
	return &ConfigLoader{version, directory}
}

// LoadSharedConfig loads the shared configuration
func (c *ConfigLoader) LoadSharedConfig() (res *shared.Config, err error) {
	sharedFilePath := filepath.Join(c.directory, "shared.yaml")
	configFileBytes, err := getConfigBytes(sharedFilePath)

	if err != nil {
		return nil, err
	}

	cf, err := LoadSharedConfigFile(configFileBytes...)

	if err != nil {
		return nil, err
	}

	return GetSharedConfigFromConfigFile(cf)
}

// LoadDatabaseConfig loads the database configuration
func (c *ConfigLoader) LoadDatabaseConfig() (res *database.Config, err error) {
	sharedFilePath := filepath.Join(c.directory, "database.yaml")
	configFileBytes, err := getConfigBytes(sharedFilePath)

	if err != nil {
		return nil, err
	}

	cf, err := LoadDatabaseConfigFile(configFileBytes...)

	if err != nil {
		return nil, err
	}

	return GetDatabaseConfigFromConfigFile(cf)
}

// LoadServerConfig loads the server configuration
func (c *ConfigLoader) LoadServerConfig() (res *server.Config, err error) {
	databaseConfig, err := c.LoadDatabaseConfig()

	if err != nil {
		return nil, fmt.Errorf("could not load database config for server: %w", err)
	}

	sharedConfig, err := c.LoadSharedConfig()

	if err != nil {
		return nil, fmt.Errorf("could not load shared config for server: %w", err)
	}

	sharedFilePath := filepath.Join(c.directory, "server.yaml")
	configFileBytes, err := getConfigBytes(sharedFilePath)

	if err != nil {
		return nil, err
	}

	cf, err := LoadServerConfigFile(configFileBytes...)

	if err != nil {
		return nil, err
	}

	return GetServerConfigFromConfigFile(c.version, cf, databaseConfig, sharedConfig)
}

// LoadTemporalConfig loads the temporal server configuration
func (c *ConfigLoader) LoadTemporalConfig() (res *temporalconfig.Config, err error) {
	databaseConfig, err := c.LoadDatabaseConfig()

	if err != nil {
		return nil, fmt.Errorf("could not load database config for temporal: %w", err)
	}

	sharedFilePath := filepath.Join(c.directory, "temporal.yaml")
	configFileBytes, err := getConfigBytes(sharedFilePath)

	if err != nil {
		return nil, err
	}

	cf, err := LoadTemporalConfigFile(configFileBytes...)

	if err != nil {
		return nil, err
	}

	return GetTemporalConfigFromConfigFile(cf, databaseConfig)
}

// LoadRunnerConfig loads the runner configuration
func (c *ConfigLoader) LoadRunnerConfig() (res *runner.Config, err error) {
	sharedConfig, err := c.LoadSharedConfig()

	if err != nil {
		return nil, fmt.Errorf("could not load shared config for runner: %w", err)
	}

	sharedFilePath := filepath.Join(c.directory, "runner.yaml")
	configFileBytes, err := getConfigBytes(sharedFilePath)

	if err != nil {
		return nil, err
	}

	cf, err := LoadRunnerConfigFile(configFileBytes...)

	if err != nil {
		return nil, err
	}

	return GetRunnerConfigFromConfigFile(cf, sharedConfig)
}

// LoadRunnerConfig loads the cli configuration
func (c *ConfigLoader) LoadCLIConfig() (res *cli.Config, v *viper.Viper, err error) {
	sharedConfig, err := c.LoadSharedConfig()

	if err != nil {
		return nil, nil, fmt.Errorf("could not load shared config for runner: %w", err)
	}

	sharedFilePath := filepath.Join(c.directory, "hatchet.yaml")
	configFileBytes, err := getConfigBytes(sharedFilePath)

	if err != nil {
		return nil, nil, err
	}

	cf, v, err := LoadCLIConfigFile(configFileBytes...)

	if err != nil {
		return nil, nil, err
	}

	res, err = GetCLIConfigFromConfigFile(cf, sharedConfig)

	if err != nil {
		return nil, nil, err
	}

	return res, v, nil
}

// LoadBackgroundWorkerConfig loads the background worker configuration
func (c *ConfigLoader) LoadBackgroundWorkerConfig() (res *worker.BackgroundConfig, err error) {
	sharedConfig, err := c.LoadSharedConfig()

	if err != nil {
		return nil, fmt.Errorf("could not load shared config for background worker: %w", err)
	}

	databaseConfig, err := c.LoadDatabaseConfig()

	if err != nil {
		return nil, fmt.Errorf("could not load database config for background worker: %w", err)
	}

	sharedFilePath := filepath.Join(c.directory, "background_worker.yaml")
	configFileBytes, err := getConfigBytes(sharedFilePath)

	if err != nil {
		return nil, err
	}

	cf, err := LoadBackgroundWorkerConfigFile(configFileBytes...)

	if err != nil {
		return nil, err
	}

	return GetBackgroundWorkerConfigFromConfigFile(cf, databaseConfig, sharedConfig)
}

// LoadRunnerWorkerConfig loads the runner worker configuration
func (c *ConfigLoader) LoadRunnerWorkerConfig() (res *worker.RunnerConfig, err error) {
	sharedConfig, err := c.LoadSharedConfig()

	if err != nil {
		return nil, fmt.Errorf("could not load shared config for background worker: %w", err)
	}

	sharedFilePath := filepath.Join(c.directory, "runner_worker.yaml")
	configFileBytes, err := getConfigBytes(sharedFilePath)

	if err != nil {
		return nil, err
	}

	cf, err := LoadRunnerWorkerConfigFile(configFileBytes...)

	if err != nil {
		return nil, err
	}

	return GetRunnerWorkerConfigFromConfigFile(cf, sharedConfig)
}

func getConfigBytes(configFilePath string) ([][]byte, error) {
	configFileBytes := make([][]byte, 0)

	if fileExists(configFilePath) {
		fileBytes, err := ioutil.ReadFile(configFilePath)

		if err != nil {
			return nil, fmt.Errorf("could not read config file at path %s: %w", configFilePath, err)
		}

		configFileBytes = append(configFileBytes, fileBytes)
	}

	return configFileBytes, nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func GetSharedConfigFromConfigFile(sharedC *shared.ConfigFile) (res *shared.Config, err error) {
	l := logger.NewConsole(sharedC.Debug)

	errorAlerter := erroralerter.NoOpAlerter{}

	var temporalClient *temporal.Client

	if sharedC.Temporal.Client.TemporalEnabled {
		temporalClient, err = temporal.NewTemporalClient(&temporal.ClientOpts{
			BroadcastAddress: sharedC.Temporal.Client.TemporalBroadcastAddress,
			HostPort:         sharedC.Temporal.Client.TemporalHostPort,
			Namespace:        sharedC.Temporal.Client.TemporalNamespace,
			BearerToken:      sharedC.Temporal.Client.TemporalBearerToken,
			ClientCertFile:   sharedC.Temporal.Client.TemporalClientTLSCertFile,
			ClientKeyFile:    sharedC.Temporal.Client.TemporalClientTLSKeyFile,
			RootCAFile:       sharedC.Temporal.Client.TemporalClientTLSRootCAFile,
			TLSServerName:    sharedC.Temporal.Client.TemporalTLSServerName,
		})

		if err != nil {
			return nil, err
		}
	}

	return &shared.Config{
		Logger:         *l,
		ErrorAlerter:   errorAlerter,
		TemporalClient: temporalClient,
	}, nil
}

func GetDatabaseConfigFromConfigFile(dc *database.ConfigFile) (res *database.Config, err error) {
	db, err := adapter.New(dc)

	if err != nil {
		return nil, fmt.Errorf("could not load database from adapter: %v", err)
	}

	var key [32]byte

	for i, b := range []byte(dc.EncryptionKey) {
		key[i] = b
	}

	repo := gorm.NewRepository(db, &key)

	res = &database.Config{
		GormDB:     db,
		Repository: repo,
	}

	res.SetEncryptionKey(&key)

	return res, nil
}

func GetServerConfigFromConfigFile(version string, sc *server.ConfigFile, dbConfig *database.Config, sharedConfig *shared.Config) (res *server.Config, err error) {
	authConfig := server.AuthConfig{
		BasicAuthEnabled:       sc.Auth.BasicAuthEnabled,
		RestrictedEmailDomains: sc.Auth.RestrictedEmailDomains,
	}

	serverRuntimeConfig := server.ServerRuntimeConfig{
		Version:                             version,
		Port:                                sc.Runtime.Port,
		ServerURL:                           sc.Runtime.ServerURL,
		BroadcastGRPCAddress:                sc.Runtime.BroadcastGRPCAddress,
		CookieName:                          sc.Auth.Cookie.Name,
		RunBackgroundWorker:                 sc.Runtime.RunBackgroundWorker,
		RunRunnerWorker:                     sc.Runtime.RunRunnerWorker,
		RunTemporalServer:                   sc.Runtime.RunTemporalServer,
		RunStaticFileServer:                 sc.Runtime.RunStaticFileServer,
		StaticFileServerPath:                sc.Runtime.StaticFileServerPath,
		PermittedModuleDeploymentMechanisms: sc.Runtime.PermittedModuleDeploymentMechanisms,
	}

	userSessionStore, err := cookie.NewUserSessionStore(&cookie.UserSessionStoreOpts{
		SessionRepository:   dbConfig.Repository.UserSession(),
		CookieSecrets:       sc.Auth.Cookie.Secrets,
		CookieAllowInsecure: sc.Auth.Cookie.Insecure,
		CookieDomain:        sc.Auth.Cookie.Domain,
		CookieName:          sc.Auth.Cookie.Name,
	})

	if err != nil {
		return nil, fmt.Errorf("could not initialize session store: %v", err)
	}

	tokenOpts := &token.TokenOpts{
		Issuer:   sc.Auth.Token.TokenIssuerURL,
		Audience: sc.Auth.Token.TokenAudience,
	}

	if sc.Auth.Token.TokenIssuerURL == "" {
		tokenOpts.Issuer = sc.Runtime.ServerURL
	}

	if len(sc.Auth.Token.TokenAudience) == 0 {
		tokenOpts.Audience = []string{sc.Runtime.ServerURL}
	}

	var notifier notifier.UserNotifier

	if sg := sc.Notification.Sendgrid; sg.SendgridAPIKey != "" && sg.SendgridSenderEmail != "" && sg.SendgridPWResetTemplateID != "" && sg.SendgridVerifyEmailTemplateID != "" &&
		sg.SendgridInviteLinkTemplateID != "" {
		notifier = sendgrid.NewUserNotifier(&sendgrid.UserNotifierOpts{
			SharedOpts: &sendgrid.SharedOpts{
				APIKey:                 sg.SendgridAPIKey,
				SenderEmail:            sg.SendgridSenderEmail,
				RestrictedEmailDomains: sc.Auth.RestrictedEmailDomains,
			},
			PWResetTemplateID:     sg.SendgridPWResetTemplateID,
			VerifyEmailTemplateID: sg.SendgridVerifyEmailTemplateID,
			InviteLinkTemplateID:  sg.SendgridInviteLinkTemplateID,
		})

		authConfig.RequireEmailVerification = true
	} else {
		notifier = noop.NewNoOpUserNotifier(&sharedConfig.Logger)
	}

	var githubAppConf *github.GithubAppConf

	if sc.VCS.Kind == "github" {
		var err error

		githubAppConf, err = github.NewGithubAppConf(&oauth.Config{
			ClientID:     sc.VCS.Github.GithubAppClientID,
			ClientSecret: sc.VCS.Github.GithubAppClientSecret,
			Scopes:       []string{"read:user"},
			BaseURL:      sc.Runtime.ServerURL,
		}, sc.VCS.Github.GithubAppName, sc.VCS.Github.GithubAppSecretPath, sc.VCS.Github.GithubAppWebhookSecret, sc.VCS.Github.GithubAppID)

		if err != nil {
			return nil, err
		}
	}

	var storageManager filestorage.FileStorageManager

	if sc.FileStore.FileStorageKind == "s3" {
		fc := sc.FileStore.S3

		var stateEncryptionKey [32]byte

		for i, b := range []byte(fc.S3StateEncryptionKey) {
			stateEncryptionKey[i] = b
		}

		storageManager, err = s3.NewS3StorageClient(&s3.S3Options{
			AWSRegion:      fc.S3StateAWSRegion,
			AWSAccessKeyID: fc.S3StateAWSAccessKeyID,
			AWSSecretKey:   fc.S3StateAWSSecretKey,
			AWSBucketName:  fc.S3StateBucketName,
			EncryptionKey:  &stateEncryptionKey,
		})

		if err != nil {
			return nil, err
		}
	} else if sc.FileStore.FileStorageKind == "local" {
		fc := sc.FileStore.Local

		var stateEncryptionKey [32]byte

		for i, b := range []byte(fc.FileEncryptionKey) {
			stateEncryptionKey[i] = b
		}

		storageManager, err = filelocal.NewLocalFileStorageManager(fc.FileDirectory, &stateEncryptionKey)

		if err != nil {
			return nil, err
		}
	}

	var logManager logstorage.LogStorageBackend

	if sc.LogStore.LogStorageKind == "redis" {
		rc := sc.LogStore.Redis
		logManager, err = redis.NewRedisLogStorageManager(&redis.InitOpts{
			RedisHost:     rc.RedisHost,
			RedisPort:     rc.RedisPort,
			RedisUsername: rc.RedisUsername,
			RedisPassword: rc.RedisPassword,
			RedisDB:       rc.RedisDB,
		})

		if err != nil {
			return nil, err
		}
	} else if sc.LogStore.LogStorageKind == "file" {
		logManager, err = file.NewFileLogStorageManager(sc.LogStore.File.FileDirectory)

		if err != nil {
			return nil, err
		}
	}

	queueManager := queuemanager.NewDefaultModuleRunQueueManager(dbConfig.Repository)

	return &server.Config{
		DB:                    *dbConfig,
		Config:                *sharedConfig,
		AuthConfig:            authConfig,
		ServerRuntimeConfig:   serverRuntimeConfig,
		UserSessionStore:      userSessionStore,
		TokenOpts:             tokenOpts,
		UserNotifier:          notifier,
		GithubApp:             githubAppConf,
		DefaultFileStore:      storageManager,
		DefaultLogStore:       logManager,
		ModuleRunQueueManager: queueManager,
	}, nil
}

func GetTemporalConfigFromConfigFile(
	tc *temporalconfig.TemporalConfigFile,
	dbConfig *database.Config,
) (res *temporalconfig.Config, err error) {
	authConfig := &temporalconfig.InternalAuthConfig{
		InternalNamespace:  tc.TemporalInternalNamespace,
		InternalSigningKey: []byte(tc.TemporalInternalSigningKey),
	}

	tokenOpts := &token.TokenOpts{
		Issuer:   tc.Token.TokenIssuerURL,
		Audience: tc.Token.TokenAudience,
	}

	if tc.Token.TokenIssuerURL == "" {
		tokenOpts.Issuer = tc.TemporalPublicURL
	}

	if len(tc.Token.TokenAudience) == 0 {
		tokenOpts.Audience = []string{tc.TemporalPublicURL}
	}

	authConfig.InternalTokenOpts = *tokenOpts

	return &temporalconfig.Config{
		DB:                 *dbConfig,
		ConfigFile:         tc,
		InternalAuthConfig: authConfig,
	}, nil
}

func GetRunnerConfigFromConfigFile(rc *runner.ConfigFile, sharedConfig *shared.Config) (res *runner.Config, err error) {
	grpcClient, err := grpc.NewGRPCClient(
		fmt.Sprintf("%s/api/v1", rc.GRPC.GRPCServerAddress),
		rc.GRPC.GRPCToken,
		rc.Resources.TeamID,
		rc.Resources.ModuleID,
		rc.Resources.ModuleRunID,
	)

	if err != nil {
		return nil, fmt.Errorf("could not load GRPC client: %v", err)
	}

	clientConf := swagger.NewConfiguration()

	clientConf.AddDefaultHeader("Authorization", fmt.Sprintf("Bearer %s", rc.API.APIToken))

	c := swagger.NewAPIClient(clientConf)

	fileClient := fileclient.NewFileClient(rc.API.APIServerAddress, rc.API.APIToken)

	return &runner.Config{
		Config:     *sharedConfig,
		ConfigFile: rc,
		GRPCClient: grpcClient,
		APIClient:  c,
		FileClient: fileClient,
	}, nil
}

func GetCLIConfigFromConfigFile(cc *cli.ConfigFile, sharedConfig *shared.Config) (res *cli.Config, err error) {
	clientConf := swagger.NewConfiguration()

	clientConf.AddDefaultHeader("Authorization", fmt.Sprintf("Bearer %s", cc.APIToken))

	c := swagger.NewAPIClient(clientConf)

	fileClient := fileclient.NewFileClient(cc.Address, cc.APIToken)

	return &cli.Config{
		Config:     *sharedConfig,
		ConfigFile: cc,
		APIClient:  c,
		FileClient: fileClient,
	}, nil
}

func GetBackgroundWorkerConfigFromConfigFile(
	wc *worker.BackgroundConfigFile,
	dbConfig *database.Config,
	sharedConfig *shared.Config,
) (res *worker.BackgroundConfig, err error) {
	var storageManager filestorage.FileStorageManager

	if wc.FileStore.FileStorageKind == "s3" {
		fc := wc.FileStore.S3

		var stateEncryptionKey [32]byte

		for i, b := range []byte(fc.S3StateEncryptionKey) {
			stateEncryptionKey[i] = b
		}

		storageManager, err = s3.NewS3StorageClient(&s3.S3Options{
			AWSRegion:      fc.S3StateAWSRegion,
			AWSAccessKeyID: fc.S3StateAWSAccessKeyID,
			AWSSecretKey:   fc.S3StateAWSSecretKey,
			AWSBucketName:  fc.S3StateBucketName,
			EncryptionKey:  &stateEncryptionKey,
		})

		if err != nil {
			return nil, err
		}
	} else if wc.FileStore.FileStorageKind == "local" {
		fc := wc.FileStore.Local

		var stateEncryptionKey [32]byte

		for i, b := range []byte(fc.FileEncryptionKey) {
			stateEncryptionKey[i] = b
		}

		storageManager, err = filelocal.NewLocalFileStorageManager(fc.FileDirectory, &stateEncryptionKey)

		if err != nil {
			return nil, err
		}
	}

	tokenOpts := &token.TokenOpts{
		Issuer:   wc.Auth.Token.TokenIssuerURL,
		Audience: wc.Auth.Token.TokenAudience,
	}

	if wc.Auth.Token.TokenIssuerURL == "" {
		tokenOpts.Issuer = wc.ServerURL
	}

	if len(wc.Auth.Token.TokenAudience) == 0 {
		tokenOpts.Audience = []string{wc.ServerURL}
	}

	var logManager logstorage.LogStorageBackend

	if wc.LogStore.LogStorageKind == "redis" {
		rc := wc.LogStore.Redis

		logManager, err = redis.NewRedisLogStorageManager(&redis.InitOpts{
			RedisHost:     rc.RedisHost,
			RedisPort:     rc.RedisPort,
			RedisUsername: rc.RedisUsername,
			RedisPassword: rc.RedisPassword,
			RedisDB:       rc.RedisDB,
		})

		if err != nil {
			return nil, err
		}
	} else if wc.LogStore.LogStorageKind == "file" {
		logManager, err = file.NewFileLogStorageManager(wc.LogStore.File.FileDirectory)

		if err != nil {
			return nil, err
		}
	}

	queueManager := queuemanager.NewDefaultModuleRunQueueManager(dbConfig.Repository)

	var notifier notifier.IncidentNotifier

	if wc.Notifier.Kind == "sendgrid" {
		notifier = sendgrid.NewIncidentNotifier(&sendgrid.IncidentNotifierOpts{
			SharedOpts: &sendgrid.SharedOpts{
				APIKey:                 wc.Notifier.Sendgrid.SendgridAPIKey,
				SenderEmail:            wc.Notifier.Sendgrid.SendgridSenderEmail,
				RestrictedEmailDomains: wc.Auth.RestrictedEmailDomains,
			},
			IncidentTemplateID: wc.Notifier.Sendgrid.SendgridIncidentTemplateID,
		})
	} else {
		notifier = noop.NewNoOpIncidentNotifier(&sharedConfig.Logger)
	}

	return &worker.BackgroundConfig{
		Config:                *sharedConfig,
		DB:                    *dbConfig,
		ServerURL:             wc.ServerURL,
		BroadcastGRPCAddress:  wc.BroadcastGRPCAddress,
		TokenOpts:             tokenOpts,
		DefaultFileStore:      storageManager,
		DefaultLogStore:       logManager,
		ModuleRunQueueManager: queueManager,
		IncidentNotifier:      notifier,
	}, nil
}

func GetRunnerWorkerConfigFromConfigFile(
	wc *worker.RunnerConfigFile,
	sharedConfig *shared.Config,
) (res *worker.RunnerConfig, err error) {
	var provisioner provisioner.Provisioner

	if wc.Provisioner.Kind == "local" {
		provisioner = local.NewLocalProvisioner()
	}

	return &worker.RunnerConfig{
		Config:             *sharedConfig,
		DefaultProvisioner: provisioner,
	}, nil
}
