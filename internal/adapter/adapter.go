package adapter

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/hatchet-dev/hatchet/internal/config/database"
	"github.com/hatchet-dev/hatchet/internal/migrate"
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

	var db *gorm.DB
	var err error

	if configFile.Kind == "sqlite" {
		// we create the data directory if it does not exist
		sqliteDir := filepath.Dir(configFile.SQLite.SQLLitePath)

		err = os.MkdirAll(sqliteDir, os.ModePerm)

		if err != nil {
			return nil, fmt.Errorf("could not create sqlite directory: %w", err)
		}

		// we add DisableForeignKeyConstraintWhenMigrating since our sqlite does
		// not support foreign key constraints
		db, err = gorm.Open(sqlite.Open(configFile.SQLite.SQLLitePath), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			FullSaveAssociations:                     true,
			Logger:                                   gormLogger,
		})

		if err != nil {
			return nil, err
		}
	} else if configFile.Kind == "postgres" {
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
		db, err = gorm.Open(postgres.Open(targetDSN), &gorm.Config{
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
				db, err = gorm.Open(postgres.Open(targetDSN), &gorm.Config{
					FullSaveAssociations: true,
					Logger:               gormLogger,
				})

				if retryCount > 3 {
					return nil, err
				}

				if err == nil {
					return db, nil
				}

				retryCount++
			}
		}
	}

	if configFile.AutoMigrate {
		err = migrate.AutoMigrate(db, false)

		if err != nil {
			return nil, err
		}
	}

	return db, err
}
