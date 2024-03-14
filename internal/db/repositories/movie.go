package repositories

import (
	"database/sql"
	"errors"
	m "github.com/ereminiu/filmoteka/internal/models"
	"strconv"
)

type MovieRepository struct {
	db *sql.DB
}

func NewMovieRepository(db *sql.DB) *MovieRepository {
	return &MovieRepository{
		db: db,
	}
}

func (mr *MovieRepository) GetAllMovies(sortBy string) ([]m.MovieWithActors, error) {
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
		ON a.id = am.actor_id;
	`
	tx, err := mr.db.Begin()
	if err != nil {
		return nil, err
	}
	rows, err := tx.Query(sqlQuery)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	var movies []m.MovieWithActors
	for rows.Next() {
		var movie m.MovieWithActors
		if err := rows.Scan(&movie.MovieId, &movie.MovieName, &movie.MovieDescription, &movie.MovieDate,
			&movie.MovieRate, &movie.ActorId, &movie.ActorName, &movie.ActorGender, &movie.ActorBirthday); err != nil {
			tx.Rollback()
			return nil, err
		}
		movies = append(movies, movie)
	}
	return movies, tx.Commit()
}

func (mr *MovieRepository) CreateMovie(name, description, date string, rate int, actorIds []int) (int, error) {
	sqlQuery := `INSERT INTO movies (name, description, date, rate) values ($1, $2, $3, $4) RETURNING id`
	tx, err := mr.db.Begin()
	if err != nil {
		return -1, err
	}
	var movieId int
	row := tx.QueryRow(sqlQuery, name, description, date, rate)
	err = row.Scan(&movieId)
	if err != nil {
		tx.Rollback()
		return -1, err
	}
	return movieId, tx.Commit()
}

func (mr *MovieRepository) ChangeField(movieId int, field, newValue string) error {
	sqlQuery, validatedValue, err := generateSqlQuery(field, newValue)
	tx, err := mr.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(sqlQuery, validatedValue, movieId)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func generateSqlQuery(field, rawValue string) (string, any, error) {
	var value any
	sqlQuery := ``
	switch field {
	case "rate":
		intvalue, err := strconv.Atoi(rawValue)
		if err != nil {
			return "", "", err
		}
		sqlQuery = `
			UPDATE movies
			SET rate=$1
			WHERE id=$2
		`
		value = intvalue
	case "date":
		sqlQuery = `
			UPDATE movies
			SET date=$1
			WHERE id=$2
		`
		value = rawValue
	case "description":
		sqlQuery = `
			UPDATE movies
			SET description=$1
			WHERE id=$2
		`
		value = rawValue
	case "name":
		sqlQuery = `
			UPDATE movies
			SET name=$1
			WHERE id=$2
		`
		value = rawValue
	default:
		return "", "", errors.New("unknown field")
	}
	return sqlQuery, value, nil
}

func (mr *MovieRepository) DeleteField(field string) error {
	// удаление информации о фильме -> удаление информации из таблицы акеторов, которые в нем играли
	return nil
}

func (mr *MovieRepository) DeleteMovie(id int) error {
	// удаление информации о фильме -> удаление информации из таблицы акеторов, которые в нем играли
	return nil
}

func (mr *MovieRepository) SearchMovieByPattern(pattern string) ([]m.Movie, error) {
	return nil, nil
}

func (mr *MovieRepository) SearchMovieByActorNamePattern(pattern string) ([]m.Movie, error) {
	return nil, nil
}

/*
TODO: подумать, как лучше: отдельно искать фильмы по кусочку названия и по фрагменту имени или сразу по тому и тому
TODO: больше сколняюсь ко второму варианту - подумать как это будет работать

*/
