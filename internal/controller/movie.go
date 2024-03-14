package controller

import (
	"encoding/json"
	"github.com/ereminiu/filmoteka/internal/db"
	"github.com/ereminiu/filmoteka/internal/models"
	"github.com/sirupsen/logrus"
	"net/http"
)

type MovieRouter struct {
	r db.Movie
}

func NewMovieRouter(r db.Movie) *MovieRouter {
	return &MovieRouter{r: r}
}

type createMovieInput struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Date        string         `json:"date"`
	Rate        int            `json:"rate"`
	Actors      []models.Actor `json:"actors,omitempty"`
}

func (mr *MovieRouter) AddMovie(w http.ResponseWriter, r *http.Request) {
	var input createMovieInput

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		logrus.Error(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	logrus.Println(input)

	id, err := mr.r.CreateMovie(input.Name, input.Description, input.Date, input.Rate, make([]int, 0))
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
		Message: "Movie is added",
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

type getAllMoviesInput struct {
	SortBy string `json:"sort_by"`
}

func (mr *MovieRouter) GetAllMovies(w http.ResponseWriter, r *http.Request) {
	var input getAllMoviesInput

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		logrus.Error(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	logrus.Println(input)

	movies, err := mr.r.GetAllMovies(input.SortBy)
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: переложить output в более читаемую форму
	jsonResponse, err := json.Marshal(movies)
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

type changeFieldInput struct {
	MovieId  int    `json:"movie_id"`
	Field    string `json:"field"`
	NewValue string `json:"new_value"`
}

func (mr *MovieRouter) ChangeField(w http.ResponseWriter, r *http.Request) {
	var input changeFieldInput

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		logrus.Error(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	logrus.Println(input)

	err := mr.r.ChangeField(input.MovieId, input.Field, input.NewValue)
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	output := struct {
		Message string `json:"message"`
	}{
		Message: "Field is changed",
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

type deleteMovieInput struct {
	MovieId int `json:"movie_id"`
}

func (mr *MovieRouter) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	var input deleteMovieInput

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		logrus.Error(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	logrus.Println(input)

	err := mr.r.DeleteMovie(input.MovieId)
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	output := struct {
		Message string `json:"message"`
	}{
		Message: "Movie is deleted",
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
