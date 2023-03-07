package loader

import (
	"fmt"

	"github.com/hatchet-dev/hatchet/api/serverutils/erroralerter"
	"github.com/hatchet-dev/hatchet/api/v1/client/fileclient"
	"github.com/hatchet-dev/hatchet/api/v1/client/grpc"
	"github.com/hatchet-dev/hatchet/api/v1/client/swagger"
	"github.com/hatchet-dev/hatchet/internal/adapter"
	"github.com/hatchet-dev/hatchet/internal/auth/cookie"
	"github.com/hatchet-dev/hatchet/internal/auth/token"
	"github.com/hatchet-dev/hatchet/internal/config/database"
	"github.com/hatchet-dev/hatchet/internal/config/runner"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/config/shared"
	temporalconfig "github.com/hatchet-dev/hatchet/internal/config/temporal"
	"github.com/hatchet-dev/hatchet/internal/config/worker"
	"github.com/hatchet-dev/hatchet/internal/integrations/filestorage"
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
	"github.com/joeshaw/envdecode"
)

type EnvDecoderConf struct {
	ServerConfigFile           server.ConfigFile
	BackgroundWorkerConfigFile worker.BackgroundConfigFile
	RunnerWorkerConfigFile     worker.RunnerConfigFile
	RunnerConfigFile           runner.ConfigFile
	DatabaseConfigFile         database.ConfigFile
	TemporalConfigFile         temporalconfig.TemporalConfigFile
	SharedConfigFile           shared.ConfigFile
}

// ServerConfigFromEnv loads the server config file from environment variables
func ServerConfigFromEnv() (*server.ConfigFile, error) {
	var envDecoderConf EnvDecoderConf = EnvDecoderConf{}

	if err := envdecode.StrictDecode(&envDecoderConf); err != nil {
		return nil, fmt.Errorf("Failed to decode server conf: %s", err)
	}

	return &envDecoderConf.ServerConfigFile, nil
}

// BackgroundWorkerConfigFromEnv loads the background worker config file from environment variables
func BackgroundWorkerConfigFromEnv() (*worker.BackgroundConfigFile, error) {
	var envDecoderConf EnvDecoderConf = EnvDecoderConf{}

	if err := envdecode.StrictDecode(&envDecoderConf); err != nil {
		return nil, fmt.Errorf("Failed to decode server conf: %s", err)
	}

	return &envDecoderConf.BackgroundWorkerConfigFile, nil
}

// TemporalConfigFromEnv loads the temporal config file from environment variables
func TemporalConfigFromEnv() (*temporalconfig.TemporalConfigFile, error) {
	var envDecoderConf EnvDecoderConf = EnvDecoderConf{}

	if err := envdecode.StrictDecode(&envDecoderConf); err != nil {
		return nil, fmt.Errorf("Failed to decode server conf: %s", err)
	}

	return &envDecoderConf.TemporalConfigFile, nil
}

// RunnerWorkerConfigFromEnv loads the runner worker config file from environment variables
func RunnerWorkerConfigFromEnv() (*worker.RunnerConfigFile, error) {
	var envDecoderConf EnvDecoderConf = EnvDecoderConf{}

	if err := envdecode.StrictDecode(&envDecoderConf); err != nil {
		return nil, fmt.Errorf("Failed to decode server conf: %s", err)
	}

	return &envDecoderConf.RunnerWorkerConfigFile, nil
}

// RunnerConfigFromEnv loads the runner config file from environment variables
func RunnerConfigFromEnv() (*runner.ConfigFile, error) {
	var envDecoderConf EnvDecoderConf = EnvDecoderConf{}

	if err := envdecode.StrictDecode(&envDecoderConf); err != nil {
		return nil, fmt.Errorf("Failed to decode runner conf: %s", err)
	}

	return &envDecoderConf.RunnerConfigFile, nil
}

// DatabaseConfigFromEnv loads the database config file from environment variables
func DatabaseConfigFromEnv() (*database.ConfigFile, error) {
	var envDecoderConf EnvDecoderConf = EnvDecoderConf{}

	if err := envdecode.StrictDecode(&envDecoderConf); err != nil {
		return nil, fmt.Errorf("Failed to decode database conf: %s", err)
	}

	return &envDecoderConf.DatabaseConfigFile, nil
}

// SharedConfigFromEnv loads the shared config file from environment variables
func SharedConfigFromEnv() (*shared.ConfigFile, error) {
	var envDecoderConf EnvDecoderConf = EnvDecoderConf{}

	if err := envdecode.StrictDecode(&envDecoderConf); err != nil {
		return nil, fmt.Errorf("Failed to decode server conf: %s", err)
	}

	return &envDecoderConf.SharedConfigFile, nil
}

type EnvConfigLoader struct {
	version string
}

func (e *EnvConfigLoader) loadSharedConfig() (res *shared.Config, err error) {
	sharedConfig, err := SharedConfigFromEnv()

	if err != nil {
		return nil, fmt.Errorf("could not load shared config from env: %v", err)
	}

	return e.LoadSharedConfigFromConfigFile(sharedConfig)
}

func (e *EnvConfigLoader) LoadSharedConfigFromConfigFile(sharedC *shared.ConfigFile) (res *shared.Config, err error) {
	l := logger.NewConsole(sharedC.Debug)

	errorAlerter := erroralerter.NoOpAlerter{}

	var temporalClient *temporal.Client

	if sharedC.TemporalEnabled {
		temporalClient, err = temporal.NewTemporalClient(&temporal.ClientOpts{
			BroadcastAddress: sharedC.TemporalBroadcastAddress,
			HostPort:         sharedC.TemporalHostPort,
			Namespace:        sharedC.TemporalNamespace,
			BearerToken:      sharedC.TemporalBearerToken,
			ClientCertFile:   sharedC.TemporalClientTLSCertFile,
			ClientKeyFile:    sharedC.TemporalClientTLSKeyFile,
			RootCAFile:       sharedC.TemporalClientTLSRootCAFile,
			TLSServerName:    sharedC.TemporalTLSServerName,
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

func (e *EnvConfigLoader) LoadDatabaseConfig() (res *database.Config, err error) {
	dc, err := DatabaseConfigFromEnv()

	if err != nil {
		return nil, fmt.Errorf("could not load database config from env: %v", err)
	}

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

func (e *EnvConfigLoader) LoadRunnerConfigFromEnv() (res *runner.Config, err error) {
	sharedConfig, err := e.loadSharedConfig()

	if err != nil {
		return nil, fmt.Errorf("could not load shared config: %v", err)
	}

	rc, err := RunnerConfigFromEnv()

	if err != nil {
		return nil, fmt.Errorf("could not load server config from env: %v", err)
	}

	return e.LoadRunnerConfigFromConfigFile(rc, sharedConfig)
}

func (e *EnvConfigLoader) LoadRunnerConfigFromConfigFile(rc *runner.ConfigFile, sharedConfig *shared.Config) (res *runner.Config, err error) {
	grpcClient, err := grpc.NewGRPCClient(fmt.Sprintf("%s/api/v1", rc.GRPCServerAddress), rc.GRPCToken, rc.TeamID, rc.ModuleID, rc.ModuleRunID)

	if err != nil {
		return nil, fmt.Errorf("could not load GRPC client: %v", err)
	}

	clientConf := swagger.NewConfiguration()

	clientConf.AddDefaultHeader("Authorization", fmt.Sprintf("Bearer %s", rc.APIToken))

	c := swagger.NewAPIClient(clientConf)

	fileClient := fileclient.NewFileClient(rc.APIServerAddress, rc.APIToken)

	return &runner.Config{
		Config:     *sharedConfig,
		ConfigFile: rc,
		GRPCClient: grpcClient,
		APIClient:  c,
		FileClient: fileClient,
	}, nil
}

func (e *EnvConfigLoader) LoadRunnerWorkerConfigFromEnv() (res *worker.RunnerConfig, err error) {
	sharedConfig, err := e.loadSharedConfig()

	if err != nil {
		return nil, fmt.Errorf("could not load shared config: %v", err)
	}

	wc, err := RunnerWorkerConfigFromEnv()

	if err != nil {
		return nil, fmt.Errorf("could not load server config from env: %v", err)
	}

	return e.LoadRunnerWorkerConfigFromConfigFile(wc, sharedConfig)
}

func (e *EnvConfigLoader) LoadRunnerWorkerConfigFromConfigFile(
	wc *worker.RunnerConfigFile,
	sharedConfig *shared.Config,
) (res *worker.RunnerConfig, err error) {
	var provisioner provisioner.Provisioner

	if wc.ProvisionerRunnerMethod == "local" {
		provisioner = local.NewLocalProvisioner()
	}

	return &worker.RunnerConfig{
		Config:             *sharedConfig,
		DefaultProvisioner: provisioner,
	}, nil
}

func (e *EnvConfigLoader) LoadBackgroundWorkerConfigFromEnv() (res *worker.BackgroundConfig, err error) {
	sharedConfig, err := e.loadSharedConfig()

	if err != nil {
		return nil, fmt.Errorf("could not load shared config: %v", err)
	}

	dbConfig, err := e.LoadDatabaseConfig()

	if err != nil {
		return nil, fmt.Errorf("could not load database config: %v", err)
	}

	wc, err := BackgroundWorkerConfigFromEnv()

	if err != nil {
		return nil, fmt.Errorf("could not load background worker config from env: %v", err)
	}

	return e.LoadBackgroundWorkerConfigFromConfigFile(wc, dbConfig, sharedConfig)
}

func (e *EnvConfigLoader) LoadBackgroundWorkerConfigFromConfigFile(
	wc *worker.BackgroundConfigFile,
	dbConfig *database.Config,
	sharedConfig *shared.Config,
) (res *worker.BackgroundConfig, err error) {
	var storageManager filestorage.FileStorageManager

	if e.hasS3StateAppVars(wc.S3StateStore) {
		fc := wc.S3StateStore

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
	}

	tokenOpts := &token.TokenOpts{
		Issuer:   wc.TokenIssuerURL,
		Audience: wc.TokenAudience,
	}

	if wc.TokenIssuerURL == "" {
		tokenOpts.Issuer = wc.ServerURL
	}

	if len(wc.TokenAudience) == 0 {
		tokenOpts.Audience = []string{wc.ServerURL}
	}

	var logManager logstorage.LogStorageBackend

	if wc.LogStoreConfig.LogStorageKind == "redis" {
		if e.hasRedisVars(wc.LogStoreConfig) {
			rc := wc.LogStoreConfig
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
		} else {
			return nil, fmt.Errorf("redis environment variables not set")
		}
	} else if wc.LogStoreConfig.LogStorageKind == "file" {
		logManager, err = file.NewFileLogStorageManager(wc.LogStoreConfig.FileDirectory)

		if err != nil {
			return nil, err
		}
	}

	queueManager := queuemanager.NewDefaultModuleRunQueueManager(dbConfig.Repository)

	var notifier notifier.IncidentNotifier

	if wc.SendgridAPIKey != "" && wc.SendgridSenderEmail != "" && wc.SendgridIncidentTemplateID != "" {
		notifier = sendgrid.NewIncidentNotifier(&sendgrid.IncidentNotifierOpts{
			SharedOpts: &sendgrid.SharedOpts{
				APIKey:                 wc.SendgridAPIKey,
				SenderEmail:            wc.SendgridSenderEmail,
				RestrictedEmailDomains: wc.RestrictedEmailDomains,
			},
			IncidentTemplateID: wc.SendgridIncidentTemplateID,
		})
	} else {
		notifier = noop.NewNoOpIncidentNotifier()
	}

	return &worker.BackgroundConfig{
		Config:                *sharedConfig,
		DB:                    *dbConfig,
		ServerURL:             wc.ServerURL,
		TokenOpts:             tokenOpts,
		DefaultFileStore:      storageManager,
		DefaultLogStore:       logManager,
		ModuleRunQueueManager: queueManager,
		IncidentNotifier:      notifier,
	}, nil
}

func (e *EnvConfigLoader) LoadTemporalWorkerConfigFromEnv() (res *temporalconfig.Config, err error) {
	dbConfig, err := e.LoadDatabaseConfig()

	if err != nil {
		return nil, fmt.Errorf("could not load database config: %v", err)
	}

	tc, err := TemporalConfigFromEnv()

	if err != nil {
		return nil, fmt.Errorf("could not load temporal config from env: %v", err)
	}

	return e.LoadTemporalConfigFromConfigFile(tc, dbConfig)
}

func (e *EnvConfigLoader) LoadTemporalConfigFromConfigFile(
	tc *temporalconfig.TemporalConfigFile,
	dbConfig *database.Config,
) (res *temporalconfig.Config, err error) {
	authConfig := &temporalconfig.InternalAuthConfig{
		InternalNamespace:  tc.TemporalInternalNamespace,
		InternalSigningKey: []byte(tc.TemporalInternalSigningKey),
	}

	tokenOpts := &token.TokenOpts{
		Issuer:   tc.TemporalInternalTokenIssuerURL,
		Audience: tc.TemporalInternalTokenAudience,
	}

	if tc.TemporalInternalTokenIssuerURL == "" {
		tokenOpts.Issuer = tc.TemporalPublicURL
	}

	if len(tc.TemporalInternalTokenAudience) == 0 {
		tokenOpts.Audience = []string{tc.TemporalPublicURL}
	}

	authConfig.InternalTokenOpts = *tokenOpts

	return &temporalconfig.Config{
		DB:                 *dbConfig,
		ConfigFile:         tc,
		InternalAuthConfig: authConfig,
	}, nil
}

func (e *EnvConfigLoader) LoadServerConfigFromEnv() (res *server.Config, err error) {
	sharedConfig, err := e.loadSharedConfig()

	if err != nil {
		return nil, fmt.Errorf("could not load shared config: %v", err)
	}

	dbConfig, err := e.LoadDatabaseConfig()

	if err != nil {
		return nil, fmt.Errorf("could not load database config: %v", err)
	}

	sc, err := ServerConfigFromEnv()

	if err != nil {
		return nil, fmt.Errorf("could not load server config from env: %v", err)
	}

	return e.LoadServerConfigFromConfigFile(sc, dbConfig, sharedConfig)
}

func (e *EnvConfigLoader) LoadServerConfigFromConfigFile(sc *server.ConfigFile, dbConfig *database.Config, sharedConfig *shared.Config) (res *server.Config, err error) {
	authConfig := server.AuthConfig{
		BasicAuthEnabled:       sc.BasicAuthEnabled,
		RestrictedEmailDomains: sc.RestrictedEmailDomains,
	}

	serverRuntimeConfig := server.ServerRuntimeConfig{
		Port:       sc.Port,
		ServerURL:  sc.ServerURL,
		CookieName: sc.CookieName,
	}

	userSessionStore, err := cookie.NewUserSessionStore(&cookie.UserSessionStoreOpts{
		SessionRepository:   dbConfig.Repository.UserSession(),
		CookieSecrets:       sc.CookieSecrets,
		CookieAllowInsecure: sc.CookieAllowInsecure,
		CookieDomain:        sc.CookieDomain,
		CookieName:          sc.CookieName,
	})

	if err != nil {
		return nil, fmt.Errorf("could not initialize session store: %v", err)
	}

	tokenOpts := &token.TokenOpts{
		Issuer:   sc.TokenIssuerURL,
		Audience: sc.TokenAudience,
	}

	if sc.TokenIssuerURL == "" {
		tokenOpts.Issuer = sc.ServerURL
	}

	if len(sc.TokenAudience) == 0 {
		tokenOpts.Audience = []string{sc.ServerURL}
	}

	var notifier notifier.UserNotifier

	if sc.SendgridAPIKey != "" && sc.SendgridSenderEmail != "" && sc.SendgridPWResetTemplateID != "" && sc.SendgridVerifyEmailTemplateID != "" &&
		sc.SendgridInviteLinkTemplateID != "" {
		notifier = sendgrid.NewUserNotifier(&sendgrid.UserNotifierOpts{
			SharedOpts: &sendgrid.SharedOpts{
				APIKey:                 sc.SendgridAPIKey,
				SenderEmail:            sc.SendgridSenderEmail,
				RestrictedEmailDomains: sc.RestrictedEmailDomains,
			},
			PWResetTemplateID:     sc.SendgridPWResetTemplateID,
			VerifyEmailTemplateID: sc.SendgridVerifyEmailTemplateID,
			InviteLinkTemplateID:  sc.SendgridInviteLinkTemplateID,
		})

		authConfig.RequireEmailVerification = true
	} else {
		notifier = noop.NewNoOpUserNotifier()
	}

	var githubAppConf *github.GithubAppConf

	if e.hasGithubAppVars(sc) {
		var err error

		githubAppConf, err = github.NewGithubAppConf(&oauth.Config{
			ClientID:     sc.GithubAppClientID,
			ClientSecret: sc.GithubAppClientSecret,
			Scopes:       []string{"read:user"},
			BaseURL:      sc.ServerURL,
		}, sc.GithubAppName, sc.GithubAppSecretPath, sc.GithubAppWebhookSecret, sc.GithubAppID)

		if err != nil {
			return nil, err
		}
	}

	var storageManager filestorage.FileStorageManager

	if e.hasS3StateAppVars(sc.S3StateStore) {
		fc := sc.S3StateStore

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
	}

	var logManager logstorage.LogStorageBackend

	if sc.LogStoreConfig.LogStorageKind == "redis" {
		if e.hasRedisVars(sc.LogStoreConfig) {
			rc := sc.LogStoreConfig
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
		} else {
			return nil, fmt.Errorf("redis environment variables not set")
		}
	} else if sc.LogStoreConfig.LogStorageKind == "file" {
		logManager, err = file.NewFileLogStorageManager(sc.LogStoreConfig.FileDirectory)

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

func (e *EnvConfigLoader) hasGithubAppVars(sc *server.ConfigFile) bool {
	return sc.GithubAppClientID != "" &&
		sc.GithubAppClientSecret != "" &&
		sc.GithubAppName != "" &&
		sc.GithubAppWebhookSecret != "" &&
		sc.GithubAppSecretPath != "" &&
		sc.GithubAppID != ""
}

func (e *EnvConfigLoader) hasS3StateAppVars(fc shared.FileStorageConfigFile) bool {
	return fc.S3StateAWSAccessKeyID != "" &&
		fc.S3StateAWSRegion != "" &&
		fc.S3StateAWSSecretKey != "" &&
		fc.S3StateBucketName != "" &&
		fc.S3StateEncryptionKey != ""
}

func (e *EnvConfigLoader) hasRedisVars(rc shared.LogStoreConfigFile) bool {
	return rc.RedisHost != ""
}
