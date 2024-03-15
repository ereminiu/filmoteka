package repositories

import "database/sql"

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (ar *AuthRepository) GetUser(username, passwordHash string) (int, error) {
	sqlQuery := `SELECT id FROM users WHERE uesrname=$1 AND password_hash=$2`
	tx, err := ar.db.Begin()
	if err != nil {
		return -1, err
	}
	row := tx.QueryRow(sqlQuery, username, passwordHash)
	var userId int
	if err := row.Scan(&userId); err != nil {
		tx.Rollback()
		return -1, err
	}
	return userId, tx.Commit()
}

func (ar *AuthRepository) CreateUser(name, username, passwordHash string) (int, error) {
	sqlQuery := `INSERT INTO users (name, uesrname, password_hash) VALUES ($1, $2, $3) RETURNING id`
	tx, err := ar.db.Begin()
	if err != nil {
		return -1, err
	}
	row := tx.QueryRow(sqlQuery, name, username, passwordHash)
	var userId int
	if err := row.Scan(&userId); err != nil {
		tx.Rollback()
		return -1, err
	}
	return userId, tx.Commit()
}
