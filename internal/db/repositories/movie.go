package repositories

import (
	"database/sql"
	m "github.com/ereminiu/filmoteka/internal/models"
	"github.com/sirupsen/logrus"
)

type MovieRepository struct {
	db *sql.DB
}

func NewMovieRepository(db *sql.DB) *MovieRepository {
	return &MovieRepository{
		db: db,
	}
}

func (mr *MovieRepository) CreateMovie(name, description, date string, rate int, actorIds []int) (int, error) {
	sqlQuery := `INSERT INTO movies (name, description, date, rate) values ($1, $2, $3, $4) RETURNING id`
	tx, err := mr.db.Begin()
	if err != nil {
		logrus.Error(err)
		return -1, err
	}
	var movieId int
	row := tx.QueryRow(sqlQuery, name, description, date, rate)
	err = row.Scan(&movieId)
	if err != nil {
		logrus.Error(err)
		tx.Rollback()
		return -1, err
	}
	return movieId, tx.Commit()
}

func (mr *MovieRepository) ChangeField(field, newValue string) error {
	return nil
}

func (mr *MovieRepository) DeleteField(field string) error {
	// удаление информации о фильме -> удаление информации из таблицы акеторов, которые в нем играли
	return nil
}

func (mr *MovieRepository) DeleteMovie(id int) error {
	// удаление информации о фильме -> удаление информации из таблицы акеторов, которые в нем играли
	return nil
}

func (mr *MovieRepository) GetAllMovies(sortBy string) ([]m.Movie, error) {
	// возвращает список всех фильмов (название, описание, дата выпуска, рейтинг и СПИСОК АКТЕРОВ)
	return nil, nil
}

func (mr *MovieRepository) SearchMovieByPattern(pattern string) ([]m.Movie, error) {
	return nil, nil
}

func (mr *MovieRepository) SearchMovieByActorNamePattern(patter string) ([]m.Movie, error) {
	return nil, nil
}

/*
TODO: подумать, как лучше: отдельно искать фильмы по кусочку названия и по фрагменту имени или сразу по тому и тому
TODO: больше сколняюсь ко второму варианту - подумать как это будет работать

*/
