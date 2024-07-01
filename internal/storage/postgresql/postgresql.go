package postgresql

import (
	"errors"
	"fmt"
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
	const errFunc = "storage.postgresql.NewStorage"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errFunc, err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errFunc, err)
	}

	driver, err := migratePostgres.WithInstance(sqlDB, &migratePostgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errFunc, err)
	}

	migration, err := migrate.NewWithDatabaseInstance("file://internal/storage/postgresql/migrations", dbname, driver)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errFunc, err)
	}

	if err := migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("%s: %w", errFunc, err)
	}

	return &Storage{DB: db}, nil
}
