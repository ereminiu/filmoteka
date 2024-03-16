package controller

import (
	"bytes"
	"errors"
	merrors "github.com/ereminiu/filmoteka/internal/db/errors"
	mock_db "github.com/ereminiu/filmoteka/internal/db/mocks"
	m "github.com/ereminiu/filmoteka/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMovieRouter_AddMovie(t *testing.T) {
	type mockBehavior func(s *mock_db.MockMovie, name, description, date string, rate int, actors []int)

	type inputMovie struct {
		Name        string
		Description string
		Date        string
		Rate        int
		Actors      []int
	}

	testTable := []struct {
		name         string
		inputMovie   inputMovie
		inputBody    string
		mockBehavior mockBehavior
		//expectedMovieId      int
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			inputMovie: inputMovie{
				Name:        "Сумерки 4",
				Description: "Настоящее кино для настощих мужчин",
				Date:        "2012-12-12",
				Rate:        0,
				Actors:      []int{1, 2, 3},
			},
			inputBody: `{"name":"Сумерки 4","description":"Настоящее кино для настощих мужчин","date":"2012-12-12","rate":0,"actors":[1,2,3]}`,
			mockBehavior: func(s *mock_db.MockMovie, name, description, date string, rate int, actors []int) {
				s.EXPECT().CreateMovie(name, description, date, rate, actors).Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1,"message":"Movie is added"}`,
		},
		{
			name: "Internal error",
			inputMovie: inputMovie{
				Name:        "Сумерки 4",
				Description: "Настоящее кино для настощих мужчин",
				Date:        "2012-12-12",
				Rate:        0,
				Actors:      []int{1, 2, 3},
			},
			inputBody: `{"name":"Сумерки 4","description":"Настоящее кино для настощих мужчин","date":"2012-12-12","rate":0,"actors":[1,2,3]}`,
			mockBehavior: func(s *mock_db.MockMovie, name, description, date string, rate int, actors []int) {
				s.EXPECT().CreateMovie(name, description, date, rate, actors).Return(-1, errors.New("any error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: merrors.ErrMovieCreation.Error() + "\n",
		},
	}
	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			movieRepos := mock_db.NewMockMovie(c)
			tc.mockBehavior(movieRepos, tc.inputMovie.Name, tc.inputMovie.Description, tc.inputMovie.Date,
				tc.inputMovie.Rate, tc.inputMovie.Actors)

			movieRouter := NewMovieRouter(movieRepos)

			// Test Router
			router := http.NewServeMux()
			router.HandleFunc("POST /add-movie", movieRouter.AddMovie)

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/add-movie", bytes.NewBufferString(tc.inputBody))

			// Perform Request
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}

func TestMovieRouter_GetAllMovies(t *testing.T) {
	type mockBehavior func(s *mock_db.MockMovie, sortBy string)

	testTable := []struct {
		name                 string
		inputSortBy          string
		inputBody            string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK",
			inputSortBy: "movie_rate",
			inputBody:   `{"sort_by":"movie_rate"}`,
			mockBehavior: func(s *mock_db.MockMovie, sortBy string) {
				response := []m.MovieWithActors{{
					MovieId:          11,
					MovieName:        "Сумерки 4",
					MovieDescription: "Кино для настоящих мужчин",
					MovieDate:        "2012-10-10T00:00:00Z",
					MovieRate:        0,
					Actors: []m.Actor{
						{
							Id:       26,
							Name:     "Valera",
							Gender:   "male",
							Birthday: "1988-01-01",
						},
					},
				}}
				s.EXPECT().GetAllMovies(sortBy).Return(response, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `[{"movie_id":11,"movie_name":"Сумерки 4","movie_description":"Кино для настоящих мужчин","movie_date":"2012-10-10T00:00:00Z","movie_rate":0,"actors":[{"id":26,"name":"Valera","gender":"male","Birthday":"1988-01-01"}]}]`,
		},
		{
			name:        "Bad Request",
			inputBody:   `{"sort_by":"movie_rrrrrate"}`,
			inputSortBy: "movie_rrrrrate",
			mockBehavior: func(s *mock_db.MockMovie, sortBy string) {
				s.EXPECT().GetAllMovies(sortBy).Return(nil, errors.New("some error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: http.StatusText(http.StatusInternalServerError) + "\n",
		},
	}
	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			movieRepos := mock_db.NewMockMovie(c)
			tc.mockBehavior(movieRepos, tc.inputSortBy)

			movieRouter := NewMovieRouter(movieRepos)

			// Test Router
			router := http.NewServeMux()
			router.HandleFunc("POST /movie-list", movieRouter.GetAllMovies)

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/movie-list", bytes.NewBufferString(tc.inputBody))

			// Perform Request
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}
