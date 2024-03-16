package controller

import (
	"encoding/json"
	"github.com/ereminiu/filmoteka/internal/db"
	merrors "github.com/ereminiu/filmoteka/internal/db/errors"
	m "github.com/ereminiu/filmoteka/internal/models"
	"github.com/sirupsen/logrus"
	"net/http"
)

type ActorRouter struct {
	r db.Actor
}

func NewActorRouter(r db.Actor) *ActorRouter {
	return &ActorRouter{r: r}
}

// @Summary Create Actor
// @Security ApiKeyAuth
// @Tags actors
// @Description create actor
// @ID add-actor
// @Accept  json
// @Produce  json
// @Param input body models.Actor true "actor data"
// @Success 200 {integer} integer 1
// @Failure 500 {string} string "Internal Server Error"
// @Failure 400 {string} string "Bad request"
// @Router /add-actor [post]
func (ar *ActorRouter) AddActor(w http.ResponseWriter, r *http.Request) {
	var input m.Actor

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		logrus.Error(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	logrus.Println(input)

	id, err := ar.r.CreateActor(input.Name, input.Gender, input.Birthday)
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	output := outputWithId{
		Id:      id,
		Message: "Actor is added",
	}

	jsonResponse, err := json.Marshal(output)
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

type deleteActorInput struct {
	Id int `json:"actor_id"`
}

// @Summary Delete Actor
// @Security ApiKeyAuth
// @Tags actors
// @Description delete actor by id
// @ID delete-actor
// @Accept  json
// @Produce  json
// @Param input body deleteActorInput true "actor_id"
// @Success 200 {object} outputWithMessage
// @Failure 500 {string} string "Internal Server Error"
// @Failure 400 {string} string "Bad request"
// @Router /delete-actor [delete]
func (ar *ActorRouter) DeleteActor(w http.ResponseWriter, r *http.Request) {
	var input deleteActorInput

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		logrus.Error(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	logrus.Println(input)

	err := ar.r.DeleteActor(input.Id)
	if err != nil {
		logrus.Error(err)
		http.Error(w, merrors.ErrDatabase.Error(), http.StatusInternalServerError)
		return
	}

	output := outputWithMessage{Message: "Actor is removed"}

	jsonResponse, err := json.Marshal(output)
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// @Summary Get All Actors
// @Tags actors
// @Description get all actors
// @ID get-all-actors
// @Accept  json
// @Produce  json
// @Success 200 {object} []models.ActorWithMovies
// @Failure 500 {string} string "Internal Server Error"
// @Failure default {string} string "error"
// @Router /actor-list [get]
func (ar *ActorRouter) GetAllActors(w http.ResponseWriter, r *http.Request) {
	actors, err := ar.r.GetAllActors()
	if err != nil {
		logrus.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(actors)
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
