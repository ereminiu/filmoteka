package main

import (
	"errors"
	"fmt"
	"github.com/ereminiu/filmoteka/internal/config"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// load configs
	cfg, err := config.LoadConfigs("test")
	if err != nil {
		log.Fatalln(err)
	}

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
		log.Fatal(err)
	}
	// m.Down() - to discard changes
	// m.Force() - to fix dirty version of migrations
	if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal(err)
	}

	version, dirty, err := m.Version()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Applied migration: %d, Dirty: %t\n", version, dirty)
}
