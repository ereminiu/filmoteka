package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	m "github.com/ereminiu/filmoteka/internal/models"
)

type ActorRepository struct {
	db *sql.DB
}

func NewActorRepository(db *sql.DB) *ActorRepository {
	return &ActorRepository{
		db: db,
	}
}

type RawActorWithMovies struct {
	ActorId          int    `db:"actor_id"`
	ActorName        string `db:"actor_name"`
	ActorGender      string `db:"actor_gender"`
	ActorBirthday    string `db:"actor_birthday"`
	MovieId          int    `db:"movie_id"`
	MovieName        string `db:"movie_name"`
	MovieDescription string `db:"movie_description"`
	MovieDate        string `db:"movie_date"`
	MovieRate        int    `db:"movie_rate"`
}

func (ar *ActorRepository) GetAllActors() ([]m.ActorWithMovies, error) {
	sqlQuery := `SELECT a.id AS "actor_id", 
		a.name AS "actor_name", 
		a.gender AS "actor_gender",
		a.birthday AS "actor_birthday",
		m.id AS "movie_id",
		m.name AS "movie_name",
		m.description AS "movie_description",
		m.date AS "movie_date",
		m.rate AS "movie_rate"
		FROM actors a
		JOIN actors_to_movies am
		ON am.actor_id = a.id
		JOIN movies m
		ON m.id = am.movie_id
	`
	tx, err := ar.db.Begin()
	if err != nil {
		return nil, err
	}
	rows, err := tx.Query(sqlQuery)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	actorsToMovies := make(map[m.Actor][]m.Movie)
	for rows.Next() {
		var rawActor RawActorWithMovies
		if err := rows.Scan(&rawActor.ActorId, &rawActor.ActorName, &rawActor.ActorGender, &rawActor.ActorBirthday, &rawActor.MovieId,
			&rawActor.MovieName, &rawActor.MovieDescription, &rawActor.MovieDate, &rawActor.MovieRate); err != nil {
			tx.Rollback()
			return nil, err
		}
		actor := m.Actor{Id: rawActor.ActorId, Name: rawActor.ActorName, Gender: rawActor.ActorGender,
			Birthday: rawActor.ActorBirthday}
		movie := m.Movie{Id: rawActor.MovieId, Name: rawActor.MovieName, Rate: rawActor.MovieRate,
			Date: rawActor.MovieDate, Description: rawActor.MovieDescription}
		actorsToMovies[actor] = append(actorsToMovies[actor], movie)
	}

	actors := make([]m.ActorWithMovies, 0)
	for actor, movies := range actorsToMovies {
		actors = append(actors, m.NewActorWithMovies(actor, movies))
	}

	return actors, tx.Commit()
}

func (ar *ActorRepository) CreateActor(name, gender, birthday string) (int, error) {
	sqlQuery := `INSERT INTO actors (name, gender, birthday) values ($1, $2, $3) ON CONFLICT DO NOTHING RETURNING id`
	tx, err := ar.db.Begin()
	if err != nil {
		return -1, err
	}

	var actorId int
	row := tx.QueryRow(sqlQuery, name, gender, birthday)
	err = row.Scan(&actorId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			actorId, err = ar.getActorId(tx, name, gender, birthday)
			if err != nil {
				tx.Rollback()
				return -1, err
			}
			return actorId, tx.Commit()
		}
		tx.Rollback()
		return -1, err
	}
	return actorId, tx.Commit()
}

func (ar *ActorRepository) getActorId(tx *sql.Tx, name, gender, birthday string) (int, error) {
	var actorId int
	sqlQuery := `SELECT id FROM actors WHERE name=$1 AND gender=$2 AND birthday=$3`
	row := tx.QueryRow(sqlQuery, name, gender, birthday)
	err := row.Scan(&actorId)
	if err != nil {
		tx.Rollback()
		return -1, err
	}
	return actorId, nil
}

func (ar *ActorRepository) ChangeField(actorId int, field, newValue string) error {
	sqlQuery := `UPDATE actors ` + fmt.Sprintf(`SET %s=$1 WHERE id=$2`, field)
	tx, err := ar.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(sqlQuery, newValue, actorId)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (ar *ActorRepository) DeleteActor(actorId int) error {
	tx, err := ar.db.Begin()
	if err != nil {
		return err
	}
	err = ar.deleteMoviesByActorId(tx, actorId)
	if err != nil {
		tx.Rollback()
		return err
	}
	sqlQuery := `DELETE FROM actors
		WHERE id=$1
	`
	_, err = tx.Exec(sqlQuery, actorId)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (ar *ActorRepository) getMoviesByActorId(tx *sql.Tx, actorId int) ([]int, error) {
	sqlQuery := `SELECT m.id
		FROM actors a
		JOIN actors_to_movies am
		ON am.actor_id = a.id
		JOIN movies m
		ON m.id = am.movie_id
		WHERE a.id=$1
	`
	rows, err := tx.Query(sqlQuery, actorId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	movieIds := make([]int, 0)
	for rows.Next() {
		var movieId int
		err = rows.Scan(&movieId)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		movieIds = append(movieIds, movieId)
	}
	return movieIds, nil
}

func (ar *ActorRepository) deleteMoviesByActorId(tx *sql.Tx, actorId int) error {
	movieIds, err := ar.getMoviesByActorId(tx, actorId)
	if err != nil {
		return err
	}
	for i := 0; i < len(movieIds); i++ {
		movieId := movieIds[i]
		sqlQuery := `DELETE FROM actors_to_movies WHERE actor_id=$1 AND movie_id=$2`
		_, err = tx.Exec(sqlQuery, actorId, movieId)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return nil
}
