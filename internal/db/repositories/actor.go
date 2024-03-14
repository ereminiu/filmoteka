package repositories

import (
	"database/sql"
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
	sqlQuery := `INSERT INTO actors (name, gender, birthday) values ($1, $2, $3) RETURNING id`
	tx, err := ar.db.Begin()
	if err != nil {
		return -1, err
	}

	var actorId int
	row := tx.QueryRow(sqlQuery, name, gender, birthday)
	err = row.Scan(&actorId)
	if err != nil {
		tx.Rollback()
		return -1, err
	}
	return actorId, tx.Commit()
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
