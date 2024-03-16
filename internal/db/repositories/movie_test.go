package repositories

import (
	"github.com/DATA-DOG/go-sqlmock"
	m "github.com/ereminiu/filmoteka/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMovieRepository_SearchMovie(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error: `%s` was not expcted during database connection", err.Error())
	}
	defer db.Close()

	r := NewMovieRepository(db)

	testTable := []struct {
		name         string
		moviePattern string
		actorPattern string
		mock         func()
		want         []m.MovieWithActors
	}{
		{
			name:         "OK",
			moviePattern: "Сум",
			actorPattern: "Val",
			mock: func() {
				mock.ExpectBegin()
				rows := mock.NewRows([]string{
					"movie_id",
					"movie_name",
					"movie_description",
					"movie_date",
					"movie_rate",
					"actor_id",
					"actor_name",
					"actor_gender",
					"actor_birthday",
				}).AddRow(
					1,
					"Сумерки 4",
					"Кино для настоящих мужчин",
					"2012-12-12",
					10,
					26,
					"Valera",
					"male",
					"1999-09-09",
				)
				sqlQuery := `SELECT m.id AS "movie_id",
				   m.name AS "movie_name",
				   m.description AS "movie_description",
				   m.date AS "movie_date",
				   m.rate AS "movie_rate",
				   a.id AS "actor_id",
				   a.name AS "actor_name",
				   a.gender AS "actor_gender",
				   a.birthday AS "actor_birthday"
					FROM movies m
					LEFT JOIN actors_to_movies am
					ON am.movie_id = m.id
					LEFT JOIN actors a
					ON a.id = am.actor_id
					WHERE m.name LIKE '%Сум%' AND a.name LIKE '%Val%'
				`
				mock.ExpectQuery(sqlQuery).WillReturnRows(rows)
				mock.ExpectCommit()
			},
			want: []m.MovieWithActors{{
				MovieId:          1,
				MovieName:        "Сумерки 4",
				MovieDescription: "Кино для настоящих мужчин",
				MovieDate:        "2012-12-12",
				MovieRate:        10,
				Actors: []m.Actor{
					{
						Id:       26,
						Name:     "Valera",
						Gender:   "male",
						Birthday: "1999-09-09",
					},
				},
			}},
		},
	}
	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()

			got, err := r.SearchMovie(tc.moviePattern, tc.actorPattern)

			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}
