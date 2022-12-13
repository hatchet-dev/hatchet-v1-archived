package database

import (
	"github.com/hatchet-dev/hatchet/internal/repository"
	"gorm.io/gorm"
)

type ConfigFile struct {
	EncryptionKey string `env:"ENCRYPTION_KEY,default=__random_strong_encryption_key__"`

	PostgresHost     string `env:"PG_DB_HOST,default=postgres"`
	PostgresPort     int    `env:"PG_DB_PORT,default=5432"`
	PostgresUsername string `env:"PG_DB_USER,default=hatchet"`
	PostgresPassword string `env:"PG_DB_PASS,default=hatchet"`
	PostgresDbName   string `env:"PG_DB_NAME,default=hatchet"`
	PostgresForceSSL bool   `env:"PG_DB_FORCE_SSL,default=false"`

	SQLLite     bool   `env:"SQL_LITE,default=false"`
	SQLLitePath string `env:"SQL_LITE_PATH,default=/hatchet/hatchet.db"`
}

type Config struct {
	GormDB     *gorm.DB
	Repository repository.Repository

	encryptionKey *[32]byte
}

func (c *Config) SetEncryptionKey(key *[32]byte) {
	c.encryptionKey = key
}

func (c *Config) GetEncryptionKey() *[32]byte {
	return c.encryptionKey
}
