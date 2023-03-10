package database

import (
	"github.com/hatchet-dev/hatchet/internal/repository"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type ConfigFile struct {
	Kind string `mapstructure:"kind" default:"sqlite"`

	EncryptionKey string `mapstructure:"encryptionKey" default:"__random_strong_encryption_key__"`

	AutoMigrate bool `mapstructure:"autoMigrate" default:"true"`

	SQLite ConfigFileSQLite `mapstructure:"sqlite"`

	Postgres ConfigFilePostgres `mapstructure:"postgres"`
}

type ConfigFileSQLite struct {
	SQLLitePath string `mapstructure:"path" default:"/hatchet/hatchet.db"`
}

type ConfigFilePostgres struct {
	PostgresHost     string `mapstructure:"host" default:"postgres"`
	PostgresPort     int    `mapstructure:"port" default:"5432"`
	PostgresUsername string `mapstructure:"username" default:"hatchet"`
	PostgresPassword string `mapstructure:"password" default:"hatchet"`
	PostgresDbName   string `mapstructure:"dbName" default:"hatchet"`
	PostgresForceSSL bool   `mapstructure:"forceSSL" default:"false"`
}

type Config struct {
	GormDB     *gorm.DB
	Repository repository.Repository

	InstanceName string

	encryptionKey *[32]byte
}

func (c *Config) SetEncryptionKey(key *[32]byte) {
	c.encryptionKey = key
}

func (c *Config) GetEncryptionKey() *[32]byte {
	return c.encryptionKey
}

func BindAllEnv(v *viper.Viper) {
	v.BindEnv("kind", "DATABASE_KIND")
	v.BindEnv("encryptionKey", "DATABASE_ENCRYPTION_KEY")
	v.BindEnv("autoMigrate", "DATABASE_AUTO_MIGRATE")
	v.BindEnv("sqlite.path", "DATABASE_SQLITE_PATH")
	v.BindEnv("postgres.host", "DATABASE_POSTGRES_HOST")
	v.BindEnv("postgres.port", "DATABASE_POSTGRES_PORT")
	v.BindEnv("postgres.username", "DATABASE_POSTGRES_USERNAME")
	v.BindEnv("postgres.password", "DATABASE_POSTGRES_PASSWORD")
	v.BindEnv("postgres.dbName", "DATABASE_POSTGRES_DB_NAME")
	v.BindEnv("postgres.forceSSL", "DATABASE_POSTGRES_FORCE_SSL")
}
