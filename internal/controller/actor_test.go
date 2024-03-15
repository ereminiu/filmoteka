package controller

import (
	"bytes"
	"errors"
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
		name                string
		inputBody           string
		inputActor          m.Actor
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:       "OK",
			inputBody:  `{"name":"Katya","gender":"female","birthday":"1997-07-14"}`,
			inputActor: m.Actor{Name: "Katya", Gender: "female", Birthday: "1997-07-14"},
			mockBehavior: func(s *mock_db.MockActor, actor m.Actor) {
				s.EXPECT().CreateActor(actor.Name, actor.Gender, actor.Birthday).Return(1, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"id":1,"message":"Actor is added"}`,
		},
		{
			name:       "Empty Field",
			inputBody:  `{"name":"Katya","birthday":"1997-07-14"}`,
			inputActor: m.Actor{Name: "Katya", Birthday: "1997-07-14"},
			mockBehavior: func(s *mock_db.MockActor, actor m.Actor) {
				s.EXPECT().CreateActor(actor.Name, actor.Gender, actor.Birthday).Return(1, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"id":1,"message":"Actor is added"}`,
		},
		{
			name: "Bad request",
			mockBehavior: func(s *mock_db.MockActor, actor m.Actor) {
			},
			expectedStatusCode:  400,
			expectedRequestBody: "Bad Request\n",
		},
		{
			name:       "Internal Error",
			inputBody:  `{"name":"Katya","gender":"female","birthday":"1997-07-14"}`,
			inputActor: m.Actor{Name: "Katya", Gender: "female", Birthday: "1997-07-14"},
			mockBehavior: func(s *mock_db.MockActor, actor m.Actor) {
				s.EXPECT().CreateActor(actor.Name, actor.Gender, actor.Birthday).Return(-1,
					errors.New("internal error"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: "internal error\n",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			actorRouter := mock_db.NewMockActor(c)
			testCase.mockBehavior(actorRouter, testCase.inputActor)

			movieRouter := NewActorRouter(actorRouter)

			// Test Router
			router := http.NewServeMux()
			router.HandleFunc("POST /add-actor", movieRouter.AddActor)

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/add-actor", bytes.NewBufferString(testCase.inputBody))

			// Perform Request
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}
