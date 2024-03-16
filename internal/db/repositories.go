package db

import (
	"database/sql"
	"github.com/ereminiu/filmoteka/internal/db/repositories"
	m "github.com/ereminiu/filmoteka/internal/models"
)

//go:generate mockgen -source=repositories.go -destination=mocks/mock.go

type Movie interface {
	CreateMovie(name, description, date string, rate int, actorIds []int) (int, error)
	ChangeField(movieId int, field, newValue string) error
	DeleteField(movieId int, field any) error
	DeleteMovie(movieId int) error
	GetAllMovies(sortBy string) ([]m.MovieWithActors, error)
	AddActorToMovie(actorId, movieId int) error
	SearchMovie(moviePattern, actorPattern string) ([]m.MovieWithActors, error)
}

type Actor interface {
	GetAllActors() ([]m.ActorWithMovies, error)
	CreateActor(name, gender, birthday string) (int, error)
	ChangeField(field, newValue string) error
	DeleteField(field string) error
	DeleteActor(actorId int) error
}

type Authorization interface {
	CreateUser(name, username, passwordHash string) (int, error)
	GetUser(username, passwordHash string) (int, error)
}

type Repositories struct {
	Authorization
	Movie
	Actor
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Authorization: repositories.NewAuthRepository(db),
		Movie:         repositories.NewMovieRepository(db),
		Actor:         repositories.NewActorRepository(db),
	}
}
