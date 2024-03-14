package models

type Actor struct {
	Id       int    `db:"id"`
	Name     string `db:"name" json:"name"`
	Gender   string `db:"gender" json:"gender"`
	Birthday string `db:"Birthday" json:"Birthday"`
}
