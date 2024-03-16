package controller

import (
	"encoding/json"
	"fmt"
	"github.com/ereminiu/filmoteka/internal/db"
	merrors "github.com/ereminiu/filmoteka/internal/db/errors"
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
	Name        string `json:"name"`
	Description string `json:"description"`
	Date        string `json:"date"`
	Rate        int    `json:"rate"`
	Actors      []int  `json:"actors,omitempty"`
}

// @Summary Create Movie
// @Security ApiKeyAuth
// @Tags movies
// @Description create movie
// @ID add-movie
// @Accept  json
// @Produce  json
// @Param input body createMovieInput true "movie data"
// @Success 200 {integer} integer 1
// @Failure 500 {string} string "Internal Server Error"
// @Failure 400 {string} string "Bad request"
// @Router /add-movie [post]
func (mr *MovieRouter) AddMovie(w http.ResponseWriter, r *http.Request) {
	var input createMovieInput

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		logrus.Error(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	logrus.Println(input)

	id, err := mr.r.CreateMovie(input.Name, input.Description, input.Date, input.Rate, input.Actors)
	if err != nil {
		logrus.Error(err)
		http.Error(w, merrors.ErrMovieCreation.Error(), http.StatusInternalServerError)
		return
	}

	output := outputWithId{
		Id:      id,
		Message: "Movie is added",
	}

	jsonResponse, err := json.Marshal(output)
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logrus.Printf("Success")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

type getAllMoviesInput struct {
	SortBy string `json:"sort_by"`
}

// @Summary Get All Movies
// @Tags movies
// @Description Return All Movies
// @ID movie-list
// @Accept  json
// @Produce  json
// @Param input body getAllMoviesInput true "sort_by param"
// @Success 200 {object} []models.MovieWithActors
// @Failure 500 {string} string "Internal Server Error"
// @Failure 400 {string} string "Bad request"
// @Router /movie-list [post]
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
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

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

// @Summary Change Movie Field
// @Security ApiKeyAuth
// @Tags movies
// @Description chanage movie field
// @ID change-movie-field
// @Accept  json
// @Produce  json
// @Param input body changeFieldInput true "actor_id"
// @Success 200 {object} outputWithMessage
// @Failure 500 {string} string "Internal Server Error"
// @Failure 400 {string} string "Bad request"
// @Router /change-movie-field [put]
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
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
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

// @Summary Delete Movie
// @Security ApiKeyAuth
// @Tags movies
// @Description delete movie by id
// @ID delete-movie
// @Accept  json
// @Produce  json
// @Param input body deleteMovieInput true "movie_id"
// @Success 200 {object} outputWithMessage
// @Failure 500 {string} string "Internal Server Error"
// @Failure 400 {string} string "Bad request"
// @Router /delete-movie [delete]
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

	output := outputWithMessage{Message: "Movie is deleted"}

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

type deleteFieldInput struct {
	MovieId int    `json:"movie_id"`
	Field   string `json:"field"`
	ActorId int    `json:"actor_id,omitempty"`
}

// @Summary Delete Actor
// @Security ApiKeyAuth
// @Tags movies
// @Description delete movie field by field and movie_id
// @ID delete-movie-field
// @Accept  json
// @Produce  json
// @Param input body deleteFieldInput true "movie_id field"
// @Success 200 {object} outputWithMessage
// @Failure 500 {string} string "Internal Server Error"
// @Failure 400 {string} string "Bad request"
// @Router /delete-movie-field [delete]
func (mr *MovieRouter) DeleteField(w http.ResponseWriter, r *http.Request) {
	var input deleteFieldInput

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		logrus.Error(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	logrus.Println(input)

	var field any
	if input.Field == "actor_id" {
		field = input.ActorId
	} else {
		field = input.Field
	}

	err := mr.r.DeleteField(input.MovieId, field)
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	output := outputWithMessage{Message: "Field is deleted"}

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

type addActorToMovieInput struct {
	ActorId int `json:"actor_id"`
	MovieId int `json:"movie_id"`
}

// @Summary Add actor to the movie
// @Security ApiKeyAuth
// @Tags movies
// @Description add actor to the movie by actor_id and movie_id
// @ID add-actor-to-movie
// @Accept  json
// @Produce  json
// @Param input body addActorToMovieInput true "input data"
// @Success 200 {object} outputWithMessage
// @Failure 500 {string} string "Internal Server Error"
// @Failure 400 {string} string "Bad request"
// @Router /add-actor-to-movie [post]
func (mr *MovieRouter) AddActorToMovie(w http.ResponseWriter, r *http.Request) {
	var input addActorToMovieInput

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		logrus.Error(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	logrus.Println(input)

	err := mr.r.AddActorToMovie(input.ActorId, input.MovieId)
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	output := outputWithMessage{Message: fmt.Sprintf("Actor %d added to the movie %d", input.ActorId, input.MovieId)}

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

type searchMovieInput struct {
	MoviePattern string `json:"movie_pattern,omitempty"`
	ActorPattern string `json:"actor_pattern,omitempty"`
}

// @Summary Search Movie by movie and actor patterns
// @Tags movies
// @Description Return movies containing movie and actor patterns
// @ID search-movies
// @Accept  json
// @Produce  json
// @Param input body searchMovieInput true "movie_pattern and actor_pattern"
// @Success 200 {object} []models.MovieWithActors
// @Failure 500 {string} string "Internal Server Error"
// @Failure 400 {string} string "Bad request"
// @Router /search-movies [post]
func (mr *MovieRouter) SearchMovie(w http.ResponseWriter, r *http.Request) {
	var input searchMovieInput

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		logrus.Error(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	logrus.Println(input)

	movies, err := mr.r.SearchMovie(input.MoviePattern, input.ActorPattern)
	if err != nil {
		logrus.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

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
