package repositories

import (
	"github.com/DATA-DOG/go-sqlmock"
	m "github.com/ereminiu/filmoteka/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestActorRepository_GetAllActors(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error: `%s` was not expcted during database connection", err.Error())
	}
	defer db.Close()

	r := NewActorRepository(db)

	testTable := []struct {
		name      string
		mock      func()
		want      []m.ActorWithMovies
		wantError bool
	}{
		{
			name: "OK",
			mock: func() {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{
					"actor_id",
					"actor_name",
					"actor_gender",
					"actor_birthday",
					"movie_id",
					"movie_name",
					"movie_description",
					"movie_date",
					"movie_rate",
				}).AddRow(
					26,
					"Valera",
					"male",
					"1999-09-09",
					1,
					"Сумерки 4",
					"Кино для настоящих мужчин",
					"2012-12-12",
					0,
				)
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
					ON m.id = am.movie_id`

				mock.ExpectQuery(sqlQuery).WillReturnRows(rows)
				mock.ExpectCommit()
			},
			want: []m.ActorWithMovies{
				{
					ActorId:       26,
					ActorName:     "Valera",
					ActorGender:   "male",
					ActorBirthday: "1999-09-09",
					Movies: []m.Movie{
						{
							Id:          1,
							Name:        "Сумерки 4",
							Description: "Кино для настоящих мужчин",
							Date:        "2012-12-12",
							Rate:        0,
						},
					},
				},
			},
		},
	}
	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()

			got, err := r.GetAllActors()

			assert.Equal(t, tc.want, got)
			assert.NoError(t, err)
		})
	}
}

func TestActorRepository_CreateActor(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error: `%s` was not expcted during database connection", err.Error())
	}
	defer db.Close()

	r := NewActorRepository(db)

	type args struct {
		name, gender, birthday string
	}

	testTable := []struct {
		name  string
		mock  func()
		input args
		want  int
	}{
		{
			name: "OK",
			input: args{
				name:     "Valera",
				gender:   "male",
				birthday: "1999-09-09",
			},
			mock: func() {
				mock.ExpectBegin()
				rows := mock.NewRows([]string{"id"}).AddRow(26)
				sqlQuery := `INSERT` + ` INTO actors`
				mock.ExpectQuery(sqlQuery).WithArgs(
					"Valera",
					"male",
					"1999-09-09",
				).WillReturnRows(rows)
				mock.ExpectCommit()
			},
			want: 26,
		},
	}

	for _, tc := range testTable {
		tc.mock()
		got, err := r.CreateActor(tc.input.name, tc.input.gender, tc.input.birthday)

		assert.NoError(t, err)
		assert.Equal(t, tc.want, got)
	}
}
