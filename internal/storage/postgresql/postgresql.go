package postgresql

import (
	"errors"
	"github.com/charmbracelet/log"
	"github.com/golang-migrate/migrate/v4"
	migratePostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Storage struct {
	DB *gorm.DB
}

func NewPostgresDB(dsn, dbname string) (*Storage, error) {
	log.Debug("Opening database connection")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	log.Debug("Successfully opened database connection")

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	log.Debug("Migrating database schema")
	driver, err := migratePostgres.WithInstance(sqlDB, &migratePostgres.Config{})
	if err != nil {
		return nil, err
	}

	migration, err := migrate.NewWithDatabaseInstance("file://internal/storage/postgresql/migrations", dbname, driver)
	if err != nil {
		return nil, err
	}

	if err = migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, err
	}
	log.Debug("Successfully migrated database schema")

	return &Storage{DB: db}, nil
}
