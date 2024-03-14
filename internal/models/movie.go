package models

type Movie struct {
	Id          int    `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Date        string `db:"date"`
	Rate        int    `db:"rate"`
}
