package models

type Actor struct {
	Id       int    `db:"id"`
	Name     string `db:"name"`
	Gender   string `db:"gender"`
	Birthday string `db:"birthdate"`
}
