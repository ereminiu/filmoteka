package app

import (
	"database/sql"
	"fmt"
	config2 "github.com/ereminiu/filmoteka/internal/config"
	"github.com/ereminiu/filmoteka/internal/db"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// TODO: add database connection

func connection() (*sql.DB, error) {
	config, err := config2.LoadConfigs("test")

	URL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.Username, config.Password, config.DBName, config.SSLMode,
	)

	database, err := sql.Open("postgres", URL)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	if err := database.Ping(); err != nil {
		logrus.Error(err)
		return nil, err
	}

	logrus.Printf("Database is ready to accept connections")

	return database, nil
}

func setUpLogger() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
}

func runServer(router *http.ServeMux) {
	srv := &http.Server{
		Handler:      router,
		Addr:         "localhost:3000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logrus.Printf("Server is starting at http://%s", srv.Addr)
	logrus.Fatalln(srv.ListenAndServe())
}

// Start HTTP server
func Run() {
	// set up logger
	setUpLogger()

	database, err := connection()
	if err != nil {
		logrus.Fatalln(err)
	}

	// Init repositories
	repos := db.NewRepositories(database)

	// Init router
	router := NewRouter(repos)

	runServer(router)
}
