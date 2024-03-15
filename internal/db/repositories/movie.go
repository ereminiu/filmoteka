package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	m "github.com/ereminiu/filmoteka/internal/models"
)

type MovieRepository struct {
	db *sql.DB
}

func NewMovieRepository(db *sql.DB) *MovieRepository {
	return &MovieRepository{
		db: db,
	}
}

type RawMovieWithActors struct {
	MovieId          int    `db:"movie_id"`
	MovieName        string `db:"movie_name"`
	MovieDescription string `db:"movie_description"`
	MovieDate        string `db:"movie_date"`
	MovieRate        int    `db:"movie_rate"`
	ActorId          int    `db:"actor_id"`
	ActorName        string `db:"actor_name"`
	ActorGender      string `db:"actor_gender"`
	ActorBirthday    string `db:"actor_birthday"`
}

func (mr *MovieRepository) GetAllMovies(sortBy string) ([]m.MovieWithActors, error) {
	// sortBy : {}
	//sqlQuery := "SELECT m.id AS movie_id FROM movies m " + fmt.Sprintf("ORDER BY %s", sortBy)
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
	` + fmt.Sprintf("ORDER BY %s", sortBy)
	tx, err := mr.db.Begin()
	if err != nil {
		return nil, err
	}
	rows, err := tx.Query(sqlQuery)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	movieToActors := make(map[m.Movie][]m.Actor)
	for rows.Next() {
		var rawMovie RawMovieWithActors
		if err := rows.Scan(&rawMovie.MovieId, &rawMovie.MovieName, &rawMovie.MovieDescription, &rawMovie.MovieDate,
			&rawMovie.MovieRate, &rawMovie.ActorId, &rawMovie.ActorName, &rawMovie.ActorGender, &rawMovie.ActorBirthday); err != nil {
			tx.Rollback()
			return nil, err
		}
		movie := m.Movie{Id: rawMovie.MovieId, Name: rawMovie.MovieName, Description: rawMovie.MovieDescription,
			Date: rawMovie.MovieDescription, Rate: rawMovie.MovieRate}
		actor := m.Actor{Id: rawMovie.ActorId, Name: rawMovie.ActorName, Gender: rawMovie.ActorGender,
			Birthday: rawMovie.ActorBirthday}
		movieToActors[movie] = append(movieToActors[movie], actor)
	}

	movies := make([]m.MovieWithActors, 0)
	for movie, actors := range movieToActors {
		movies = append(movies, m.NewMovieWithActors(movie, actors))
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

	err = mr.addActorsToMovie(tx, movieId, actorIds)
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	return movieId, tx.Commit()
}

// TODO: make this via batch insert
func (mr *MovieRepository) addActorsToMovie(tx *sql.Tx, movieId int, actorIds []int) error {
	for i := 0; i < len(actorIds); i++ {
		actorId := actorIds[i]
		sqlQuery := `INSERT INTO actors_to_movies (actor_id, movie_id) values ($1, $2)`
		_, err := tx.Exec(sqlQuery, actorId, movieId)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return nil
}

func (mr *MovieRepository) ChangeField(movieId int, field, newValue string) error {
	sqlQuery := `UPDATE movies ` + fmt.Sprintf(`SET %s=$1 WHERE id=$2`, field)
	tx, err := mr.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(sqlQuery, newValue, movieId)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (mr *MovieRepository) DeleteField(movieId int, field any) error {
	switch field.(type) {
	case int:
		// actorId
		// check if this actor is actually in this movie
		tx, err := mr.db.Begin()
		if err != nil {
			return err
		}
		actorId := field.(int)
		err = mr.deleteActorFromMovie(tx, movieId, actorId)
		if err != nil {
			return err
		}
		return tx.Commit()
	case string:
		f := field.(string)
		emptyValue := ""
		if f == "rate" {
			emptyValue = "0"
		} else if f == "date" {
			emptyValue = "1997-07-14"
		}
		return mr.ChangeField(movieId, f, emptyValue)
	default:
		return errors.New("invalid field type")
	}
}

func (mr *MovieRepository) deleteActorFromMovie(tx *sql.Tx, movieId, actorId int) error {
	sqlQuery := `DELETE FROM actors_to_movies
		WHERE actor_id=$1 AND movie_id=$2
	`
	_, err := tx.Exec(sqlQuery, actorId, movieId)
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (mr *MovieRepository) DeleteMovie(movieId int) error {
	// удаление информации о фильме -> удаление информации из таблицы акеторов, которые в нем играли
	tx, err := mr.db.Begin()
	if err != nil {
		return err
	}
	actorIds, err := mr.getActorsByMovieId(tx, movieId)
	if err != nil {
		tx.Rollback()
		return err
	}
	for i := 0; i < len(actorIds); i++ {
		actorId := actorIds[i]
		sqlQuery := `DELETE 
			FROM actors_to_movies
			WHERE actor_id=$1 AND movie_id=$2
		`
		_, err := tx.Exec(sqlQuery, actorId, movieId)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	sqlQuery := `DELETE 
		FROM movies
		WHERE id=$1
	`
	_, err = tx.Exec(sqlQuery, movieId)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (mr *MovieRepository) getActorsByMovieId(tx *sql.Tx, movieId int) ([]int, error) {
	sqlQuery := `SELECT a.id AS "actor_id"
		FROM actors a 
		JOIN actors_to_movies am 
		ON am.actor_id = a.id
		JOIN movies m 
		ON m.id = am.movie_id
		WHERE m.id = $1
	`
	rows, err := tx.Query(sqlQuery, movieId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			tx.Rollback()
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (mr *MovieRepository) AddActorToMovie(actorId, movieId int) error {
	sqlQuery := `INSERT INTO actors_to_movies (actor_id, movie_id) VALUES ($1, $2)`
	tx, err := mr.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(sqlQuery, actorId, movieId)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
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
