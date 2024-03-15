package repositories

import (
	"database/sql"
	"errors"
)

type ActorRepository struct {
	db *sql.DB
}

func NewActorRepository(db *sql.DB) *ActorRepository {
	return &ActorRepository{
		db: db,
	}
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

func (ar *ActorRepository) ChangeField(field, newValue string) error {
	return nil
}

func (ar *ActorRepository) DeleteField(field string) error {
	return nil
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
