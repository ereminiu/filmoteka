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

func TestActorRouter_AddActor(t *testing.T) {
	type mockBehavior func(s *mock_db.MockActor, actor m.Actor)

	testTable := []struct {
		name                 string
		inputBody            string
		inputActor           m.Actor
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:       "OK",
			inputBody:  `{"name":"Katya","gender":"female","birthday":"1997-07-14"}`,
			inputActor: m.Actor{Name: "Katya", Gender: "female", Birthday: "1997-07-14"},
			mockBehavior: func(s *mock_db.MockActor, actor m.Actor) {
				s.EXPECT().CreateActor(actor.Name, actor.Gender, actor.Birthday).Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1,"message":"Actor is added"}`,
		},
		{
			name:       "Empty Field",
			inputBody:  `{"name":"Katya","birthday":"1997-07-14"}`,
			inputActor: m.Actor{Name: "Katya", Birthday: "1997-07-14"},
			mockBehavior: func(s *mock_db.MockActor, actor m.Actor) {
				s.EXPECT().CreateActor(actor.Name, actor.Gender, actor.Birthday).Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1,"message":"Actor is added"}`,
		},
		{
			name: "Bad request",
			mockBehavior: func(s *mock_db.MockActor, actor m.Actor) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: "Bad Request\n",
		},
		{
			name:       "Internal Error",
			inputBody:  `{"name":"Katya","gender":"female","birthday":"1997-07-14"}`,
			inputActor: m.Actor{Name: "Katya", Gender: "female", Birthday: "1997-07-14"},
			mockBehavior: func(s *mock_db.MockActor, actor m.Actor) {
				s.EXPECT().CreateActor(actor.Name, actor.Gender, actor.Birthday).Return(-1,
					errors.New("internal error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: "internal error\n",
		},
	}
	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			actorRepos := mock_db.NewMockActor(c)
			tc.mockBehavior(actorRepos, tc.inputActor)

			actorRouter := NewActorRouter(actorRepos)

			// Test Router
			router := http.NewServeMux()
			router.HandleFunc("POST /add-actor", actorRouter.AddActor)

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/add-actor", bytes.NewBufferString(tc.inputBody))

			// Perform Request
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}

func TestActorRouter_DeleteActor(t *testing.T) {
	type mockBehavior func(s *mock_db.MockActor, actorId int)

	testTable := []struct {
		name                 string
		inputBody            string
		inputActor           int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:       "OK",
			inputBody:  `{"actor_id":1}`,
			inputActor: 1,
			mockBehavior: func(s *mock_db.MockActor, actorId int) {
				s.EXPECT().DeleteActor(actorId).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"message":"Actor is removed"}`,
		},
		{
			name:       "Wrong actor_id",
			inputBody:  `{"actor_id":1}`,
			inputActor: 1,
			mockBehavior: func(s *mock_db.MockActor, actorId int) {
				s.EXPECT().DeleteActor(actorId).Return(merrors.ErrDatabase)
			},
			expectedStatusCode:   500,
			expectedResponseBody: merrors.ErrDatabase.Error() + "\n",
		},
	}
	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			actorRepos := mock_db.NewMockActor(c)
			tc.mockBehavior(actorRepos, tc.inputActor)

			actorRouter := NewActorRouter(actorRepos)

			// Test Router
			router := http.NewServeMux()
			router.HandleFunc("DELETE /delete-actor", actorRouter.DeleteActor)

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/delete-actor", bytes.NewBufferString(tc.inputBody))

			// Perform Request
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}

func TestActorRouter_GetAllActors(t *testing.T) {
	type mockBehavior func(s *mock_db.MockActor)

	testTable := []struct {
		name                 string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			mockBehavior: func(s *mock_db.MockActor) {
				s.EXPECT().GetAllActors().Return([]m.ActorWithMovies{{
					ActorId:       26,
					ActorName:     "Valera",
					ActorGender:   "male",
					ActorBirthday: "1988-01-01",
					Movies: []m.Movie{{
						Id:          11,
						Name:        "Сумерки 4",
						Description: "Кино для настоящих мужчин",
						Date:        "2012-10-10T00:00:00Z",
						Rate:        0,
					}}},
				}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `[{"actor_id":26,"actor_name":"Valera","actor_gender":"male","actor_birthday":"1988-01-01","movies":[{"Id":11,"Name":"Сумерки 4","Description":"Кино для настоящих мужчин","Date":"2012-10-10T00:00:00Z","Rate":0}]}]`,
		},
	}
	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			actorRepos := mock_db.NewMockActor(c)
			tc.mockBehavior(actorRepos)

			actorRouter := NewActorRouter(actorRepos)

			// Test Router
			router := http.NewServeMux()
			router.HandleFunc("GET /actors-list", actorRouter.GetAllActors)

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/actors-list", nil)

			// Perform Request
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}
