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
	return -1, nil
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
