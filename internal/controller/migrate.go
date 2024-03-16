package controller

import (
	"encoding/json"
	"github.com/ereminiu/filmoteka/internal/db"
	"github.com/sirupsen/logrus"
	"net/http"
)

type MigrateRouter struct {
	m *db.Migrator
}

func NewMigrateRouter(m *db.Migrator) *MigrateRouter {
	return &MigrateRouter{m: m}
}

type migrateOutput struct {
	Version uint `json:"version"`
	Dirty   bool `json:"dirty"`
}

// @Summary Migrate Up
// @Tags migrations
// @Description Apply to Migrate Up
// @ID migrate-up
// @Produce  json
// @Success 200 {object} migrateOutput
// @Failure 500 {string} string "Internal Server Error"
// @Failure 400 {string} string "Bad request"
// @Router /migrate-up [post]
func (mr *MigrateRouter) MigrateUp(w http.ResponseWriter, r *http.Request) {
	logrus.Println("Migrate Up event")

	version, dirty, err := mr.m.Up()
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	output := migrateOutput{
		Version: version,
		Dirty:   dirty,
	}

	jsonResponse, err := json.Marshal(output)
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logrus.Println("Success")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// @Summary Migrate Down
// @Tags migrations
// @Description Apply to Migrate Down
// @ID migrate-down
// @Produce  json
// @Success 200 {object} migrateOutput
// @Failure 500 {string} string "Internal Server Error"
// @Failure 400 {string} string "Bad request"
// @Router /migrate-down [post]
func (mr *MigrateRouter) MigrateDown(w http.ResponseWriter, r *http.Request) {
	logrus.Println("Migrate Down event")

	version, dirty, err := mr.m.Down()
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	output := migrateOutput{
		Version: version,
		Dirty:   dirty,
	}

	jsonResponse, err := json.Marshal(output)
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logrus.Println("Success")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

type migrateForceInput struct {
	Version int `json:"version"`
}

// @Summary Migrate Force
// @Tags migrations
// @Description Apply to Migrate Force
// @ID migrate-force
// @Accept json
// @Produce  json
// @Param input body migrateForceInput true "version param"
// @Success 200 {object} migrateOutput
// @Failure 500 {string} string "Internal Server Error"
// @Failure 400 {string} string "Bad request"
// @Router /migrate-force [post]
func (mr *MigrateRouter) MigrateForce(w http.ResponseWriter, r *http.Request) {
	var input migrateForceInput

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		logrus.Error(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	logrus.Println(input)
	logrus.Println("Migrate Force event")

	version, dirty, err := mr.m.Force(input.Version)
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	output := migrateOutput{
		Version: version,
		Dirty:   dirty,
	}

	jsonResponse, err := json.Marshal(output)
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logrus.Println("Success")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
