package controller

import (
	"bytes"
	"fmt"
	"github.com/ereminiu/filmoteka/internal/controller/lib"
	merrors "github.com/ereminiu/filmoteka/internal/db/errors"
	mock_db "github.com/ereminiu/filmoteka/internal/db/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthRouter_GenerateToken(t *testing.T) {
	type mockBehavior func(s *mock_db.MockAuthorization, username, password string)

	type input struct {
		name     string
		username string
		password string
	}

	testTable := []struct {
		name               string
		input              input
		inputBody          string
		mockBehavior       mockBehavior
		expectedStatusCode int
		wantError          bool
		expectedError      error
	}{
		{
			name: "OK",
			input: input{
				name:     "Kai'sa",
				username: "kaisa",
				password: "qwerty",
			},
			inputBody: `{"name":"Kai'sa'","username":"kaisa","password":"qwerty"}`,
			mockBehavior: func(s *mock_db.MockAuthorization, username, password string) {
				passwordHash := lib.GeneratePasswordHash(password)
				s.EXPECT().GetUser(username, passwordHash).Return(1, nil)
			},
			expectedStatusCode: 200,
		},
		{
			name: "Authentification failed",
			input: input{
				name:     "Kai'sa",
				username: "kaisa",
				password: "qwerty",
			},
			inputBody: `{"name":"Kai'sa'","username":"kaisa","password":"qwerty"}`,
			mockBehavior: func(s *mock_db.MockAuthorization, username, password string) {
				passwordHash := lib.GeneratePasswordHash(password)
				s.EXPECT().GetUser(username, passwordHash).Return(-1, merrors.ErrAuthentication)
			},
			expectedStatusCode: 500,
			wantError:          true,
			expectedError:      merrors.ErrAuthentication,
		},
	}
	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			authRouter := mock_db.NewMockAuthorization(c)
			tc.mockBehavior(authRouter, tc.input.username, tc.input.password)

			authRouters := NewAuthRouter(authRouter)

			// Test Router
			router := http.NewServeMux()
			router.HandleFunc("POST /sign-in", authRouters.GenerateToken)

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-in", bytes.NewBufferString(tc.inputBody))

			// Perform Request
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tc.expectedStatusCode, w.Code)
			if !tc.wantError {
				token, _ := lib.GenerateToken(1)
				expectedResponseBody := fmt.Sprintf(`{"token":"%s"}`, token)
				assert.Equal(t, expectedResponseBody, w.Body.String())
			} else {
				assert.Equal(t, tc.expectedError.Error()+"\n", w.Body.String())
			}
		})
	}
}

func TestAuthRouter_CreateUser(t *testing.T) {
	type mockBehavior func(s *mock_db.MockAuthorization, name, username, password string)

	type input struct {
		name     string
		username string
		password string
	}

	testTable := []struct {
		name               string
		input              input
		inputBody          string
		mockBehavior       mockBehavior
		expectedRepsonse   string
		expectedStatusCode int
		expectedUserId     int
		expectedError      error
	}{
		{
			name: "OK",
			input: input{
				name:     "Kai'sa",
				username: "kaisa",
				password: "qwerty",
			},
			inputBody: `{"name":"Kai'sa","username":"kaisa","password":"qwerty"}`,
			mockBehavior: func(s *mock_db.MockAuthorization, name, username, password string) {
				passwordHash := lib.GeneratePasswordHash(password)
				s.EXPECT().CreateUser(name, username, passwordHash).Return(1, nil)
			},
			expectedRepsonse:   `{"id":1,"message":"User is created"}`,
			expectedStatusCode: 200,
			expectedError:      nil,
		},
		{
			name: "Sign up failed",
			input: input{
				name:     "Kai'sa",
				username: "kaisa",
				password: "qwerty",
			},
			inputBody: `{"name":"Kai'sa","username":"kaisa","password":"qwerty"}`,
			mockBehavior: func(s *mock_db.MockAuthorization, name, username, password string) {
				passwordHash := lib.GeneratePasswordHash(password)
				s.EXPECT().CreateUser(name, username, passwordHash).Return(-1, merrors.ErrUserCreation)
			},
			expectedStatusCode: 500,
			expectedError:      merrors.ErrUserCreation,
		},
	}
	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			authRouter := mock_db.NewMockAuthorization(c)
			tc.mockBehavior(authRouter, tc.input.name, tc.input.username, tc.input.password)

			authRouters := NewAuthRouter(authRouter)

			// Test Router
			router := http.NewServeMux()
			router.HandleFunc("POST /sign-up", authRouters.CreateUser)

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-up", bytes.NewBufferString(tc.inputBody))

			// Perform Request
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tc.expectedStatusCode, w.Code)
			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError.Error()+"\n", w.Body.String())
			} else {
				assert.Equal(t, tc.expectedRepsonse, w.Body.String())
			}
		})
	}
}
