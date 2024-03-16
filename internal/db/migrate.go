package db

import (
	"errors"
	"fmt"
	"github.com/ereminiu/filmoteka/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
	"log"
)

type Migrator struct {
	m *migrate.Migrate
}

func NewMigrator(cfg *config.Config) (*Migrator, error) {
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.SSLMode,
	)
	m, err := migrate.New("file://internal/migrations", databaseURL)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return &Migrator{m: m}, nil
}

func (mig *Migrator) Up() (uint, bool, error) {
	if err := mig.m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal(err)
	}
	version, dirty, err := mig.m.Version()
	if err != nil {
		logrus.Error(err)
	}
	logrus.Printf("Applied migration: %d, Dirty: %t", version, dirty)
	return version, dirty, err
}

func (mig *Migrator) Down() (uint, bool, error) {
	if err := mig.m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal(err)
	}
	version, dirty, err := mig.m.Version()
	if err != nil {
		logrus.Error(err)
	}
	logrus.Printf("Applied migration: %d, Dirty: %t", version, dirty)
	return version, dirty, err
}

func (mig *Migrator) Force(v int) (uint, bool, error) {
	if err := mig.m.Force(v); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal(err)
	}
	version, dirty, err := mig.m.Version()
	if err != nil {
		logrus.Error(err)
	}
	logrus.Printf("Applied migration: %d, Dirty: %t", version, dirty)
	return version, dirty, err
}
