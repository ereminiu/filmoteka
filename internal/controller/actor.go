package controller

import (
	"encoding/json"
	"github.com/ereminiu/filmoteka/internal/db"
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

	output := struct {
		Id      int    `json:"id"`
		Message string `json:"message"`
	}{
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	output := struct {
		Message string `json:"message"`
	}{
		Message: "Actor is removed",
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
