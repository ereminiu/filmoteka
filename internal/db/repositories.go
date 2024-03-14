package db

import (
	"database/sql"
	"github.com/ereminiu/filmoteka/internal/db/repositories"
	m "github.com/ereminiu/filmoteka/internal/models"
)

type Movie interface {
	CreateMovie(name, description, date string, rate int, actorIds []int) (int, error)
	ChangeField(field, newValue string) error
	DeleteField(field string) error
	DeleteMovie(id int) error
	GetAllMovies(sortBy string) ([]m.MovieWithActors, error)
	SearchMovieByPattern(pattern string) ([]m.Movie, error)
	SearchMovieByActorNamePattern(pattern string) ([]m.Movie, error)
}

type Actor interface {
	CreateActor(name, gender, birthday string) (int, error)
	ChangeField(field, newValue string) error
	DeleteField(field string) error
	DeleteActor(name string) error
}

type Repositories struct {
	Movie
	Actor
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Movie: repositories.NewMovieRepository(db),
		Actor: repositories.NewActorRepository(db),
	}
}
