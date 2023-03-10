package adapter

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hatchet-dev/hatchet/internal/config/database"
	"gorm.io/gorm/logger"
	gormlogger "gorm.io/gorm/logger"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// New returns a new gorm database instance
func New(configFile *database.ConfigFile) (*gorm.DB, error) {
	gormLogger := gormlogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Error,
			Colorful:      false,
		},
	)

	if configFile.Kind == "sqlite" {
		// we add DisableForeignKeyConstraintWhenMigrating since our sqlite does
		// not support foreign key constraints
		return gorm.Open(sqlite.Open(configFile.SQLite.SQLLitePath), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			FullSaveAssociations:                     true,
			Logger:                                   gormLogger,
		})
	}

	// connect to default postgres instance first
	baseDSN := fmt.Sprintf(
		"user=%s password=%s port=%d host=%s",
		configFile.Postgres.PostgresUsername,
		configFile.Postgres.PostgresPassword,
		configFile.Postgres.PostgresPort,
		configFile.Postgres.PostgresHost,
	)

	if configFile.Postgres.PostgresForceSSL {
		baseDSN = baseDSN + " sslmode=require"
	} else {
		baseDSN = baseDSN + " sslmode=disable"
	}

	targetDSN := baseDSN + " database=" + configFile.Postgres.PostgresDbName

	// open the database connection
	res, err := gorm.Open(postgres.Open(targetDSN), &gorm.Config{
		FullSaveAssociations: true,
		Logger:               gormLogger,
	})

	// retry the connection 3 times
	retryCount := 0
	timeout, _ := time.ParseDuration("5s")

	if err != nil {
		for {
			gormLogger.Warn(context.Background(), "could not connect to database. Retrying...")

			time.Sleep(timeout)
			res, err = gorm.Open(postgres.Open(targetDSN), &gorm.Config{
				FullSaveAssociations: true,
				Logger:               gormLogger,
			})

			if retryCount > 3 {
				return nil, err
			}

			if err == nil {
				return res, nil
			}

			retryCount++
		}
	}

	return res, err
}
