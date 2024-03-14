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

func (ar *ActorRepository) DeleteActor(name string) error {
	return nil
}
