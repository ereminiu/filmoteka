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
		name                 string
		inputMovie           inputMovie
		inputBody            string
		mockBehavior         mockBehavior
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

func TestMovieRouter_ChangeField(t *testing.T) {
	type mockBehavior func(s *mock_db.MockMovie, movieId int, field, newValue string)

	type args struct {
		MovieId  int
		Field    string
		NewValue string
	}

	testTable := []struct {
		name                 string
		args                 args
		inputBody            string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			args: args{
				MovieId:  1,
				Field:    "name",
				NewValue: "Sumerki",
			},
			inputBody: `{"movie_id":1,"field":"name","new_value":"Sumerki"}`,
			mockBehavior: func(s *mock_db.MockMovie, movieId int, field, newValue string) {
				s.EXPECT().ChangeField(movieId, field, newValue).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"message":"Field is changed"}`,
		},
		{
			name: "Bad request",
			args: args{
				MovieId:  228,
				Field:    "name",
				NewValue: "Sumerki",
			},
			inputBody: `{"movie_id":228,"field":"name","new_value":"Sumerki"}`,
			mockBehavior: func(s *mock_db.MockMovie, movieId int, field, newValue string) {
				s.EXPECT().ChangeField(movieId, field, newValue).Return(errors.New("some error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: http.StatusText(http.StatusInternalServerError) + "\n",
		}}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			movieRepos := mock_db.NewMockMovie(c)
			tc.mockBehavior(movieRepos, tc.args.MovieId, tc.args.Field, tc.args.NewValue)

			movieRouter := NewMovieRouter(movieRepos)

			// Test Router
			router := http.NewServeMux()
			router.HandleFunc("POST /change-movie-field", movieRouter.ChangeField)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/change-movie-field", bytes.NewBufferString(tc.inputBody))

			// Perform Request
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}

func TestMovieRouter_DeleteMovie(t *testing.T) {
	type mockBehavior func(s *mock_db.MockMovie, movieId int)

	testTable := []struct {
		name                 string
		movieId              int
		inputBody            string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			movieId:   1,
			inputBody: `{"movie_id":1}`,
			mockBehavior: func(s *mock_db.MockMovie, movieId int) {
				s.EXPECT().DeleteMovie(movieId).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"message":"Movie is deleted"}`,
		},
		{
			name:      "Internal Error",
			movieId:   1,
			inputBody: `{"movie_id":1}`,
			mockBehavior: func(s *mock_db.MockMovie, movieId int) {
				s.EXPECT().DeleteMovie(movieId).Return(errors.New("some error"))
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
			tc.mockBehavior(movieRepos, tc.movieId)

			movieRouter := NewMovieRouter(movieRepos)

			// Test Router
			router := http.NewServeMux()
			router.HandleFunc("DELETE /delete-movie", movieRouter.DeleteMovie)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/delete-movie", bytes.NewBufferString(tc.inputBody))

			// Perform Request
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}

func TestMovieRouter_DeleteField(t *testing.T) {
	type mockBehavior func(s *mock_db.MockMovie, movieId int, field string)

	testTable := []struct {
		name                 string
		movieId              int
		field                string
		inputBody            string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			movieId:   1,
			field:     "date",
			inputBody: `{"movie_id":1,"field":"date"}`,
			mockBehavior: func(s *mock_db.MockMovie, movieId int, field string) {
				s.EXPECT().DeleteField(movieId, field).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"message":"Field is deleted"}`,
		},
		{
			name:      "Internal Error",
			movieId:   1,
			field:     "date",
			inputBody: `{"movie_id":1,"field":"date"}`,
			mockBehavior: func(s *mock_db.MockMovie, movieId int, field string) {
				s.EXPECT().DeleteField(movieId, field).Return(errors.New("some error"))
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
			tc.mockBehavior(movieRepos, tc.movieId, tc.field)

			movieRouter := NewMovieRouter(movieRepos)

			// Test Router
			router := http.NewServeMux()
			router.HandleFunc("DELETE /delete-movie-field", movieRouter.DeleteField)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/delete-movie-field", bytes.NewBufferString(tc.inputBody))

			// Perform Request
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}

func TestMovieRouter_AddActorToMovie(t *testing.T) {
	type mockBehavior func(s *mock_db.MockMovie, movieId, actorId int)

	testTable := []struct {
		name                 string
		movieId              int
		actorId              int
		inputBody            string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			movieId:   1,
			actorId:   1,
			inputBody: `{"actor_id":1,"movie_id":1}`,
			mockBehavior: func(s *mock_db.MockMovie, movieId, actorId int) {
				s.EXPECT().AddActorToMovie(actorId, movieId).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"message":"Actor 1 added to the movie 1"}`,
		},
		{
			name:      "Internal Error",
			movieId:   1,
			actorId:   1,
			inputBody: `{"actor_id":1,"movie_id":1}`,
			mockBehavior: func(s *mock_db.MockMovie, movieId, actorId int) {
				s.EXPECT().AddActorToMovie(actorId, movieId).Return(errors.New("some error"))
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
			tc.mockBehavior(movieRepos, tc.actorId, tc.movieId)

			movieRouter := NewMovieRouter(movieRepos)

			// Test Router
			router := http.NewServeMux()
			router.HandleFunc("POST /add-actor-to-movie", movieRouter.AddActorToMovie)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/add-actor-to-movie", bytes.NewBufferString(tc.inputBody))

			// Perform Request
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}

func TestMovieRouter_SearchMovie(t *testing.T) {
	type mockBehavior func(s *mock_db.MockMovie, moviePattern, actorPattern string)

	testTable := []struct {
		name                 string
		moviePattern         string
		actorPattern         string
		inputBody            string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:         "OK",
			moviePattern: "Сум",
			actorPattern: "Val",
			inputBody:    `{"movie_pattern":"Сум","actor_pattern":"Val"}`,
			mockBehavior: func(s *mock_db.MockMovie, moviePattern, actorPattern string) {
				response := []m.MovieWithActors{
					{
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
					},
				}
				s.EXPECT().SearchMovie(moviePattern, actorPattern).Return(response, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `[{"movie_id":11,"movie_name":"Сумерки 4","movie_description":"Кино для настоящих мужчин","movie_date":"2012-10-10T00:00:00Z","movie_rate":0,"actors":[{"id":26,"name":"Valera","gender":"male","Birthday":"1988-01-01"}]}]`,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			movieRepos := mock_db.NewMockMovie(c)
			tc.mockBehavior(movieRepos, tc.moviePattern, tc.actorPattern)

			movieRouter := NewMovieRouter(movieRepos)

			// Test Router
			router := http.NewServeMux()
			router.HandleFunc("POST /search-movie", movieRouter.SearchMovie)

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/search-movie", bytes.NewBufferString(tc.inputBody))

			// Perform Request
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}
