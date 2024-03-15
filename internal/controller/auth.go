package controller

import (
	"context"
	"encoding/json"
	"github.com/ereminiu/filmoteka/internal/controller/lib"
	"github.com/ereminiu/filmoteka/internal/db"
	merrors "github.com/ereminiu/filmoteka/internal/db/errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

const (
	userCtx = "userId"
)

type AuthRouter struct {
	r db.Authorization
}

func NewAuthRouter(r db.Authorization) *AuthRouter {
	return &AuthRouter{r: r}
}

type outputType map[string]interface{}

type createUserInput struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (ar *AuthRouter) GenerateToken(w http.ResponseWriter, r *http.Request) {
	var input createUserInput

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		logrus.Error(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	logrus.Println(input)

	input.Password = lib.GeneratePasswordHash(input.Password)
	userId, err := ar.r.GetUser(input.Username, input.Password)
	if err != nil {
		logrus.Error(err)
		http.Error(w, merrors.ErrAuthentication.Error(), http.StatusInternalServerError)
		return
	}

	token, err := lib.GenerateToken(userId)
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(outputType{
		"token": token,
	})
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (ar *AuthRouter) CreateUser(w http.ResponseWriter, r *http.Request) {
	var input createUserInput

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		logrus.Error(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	logrus.Println(input)

	input.Password = lib.GeneratePasswordHash(input.Password)
	id, err := ar.r.CreateUser(input.Name, input.Username, input.Password)
	if err != nil {
		logrus.Error(err)
		http.Error(w, merrors.ErrUserCreation.Error(), http.StatusInternalServerError)
		return
	}

	output := struct {
		Id      int    `json:"id"`
		Message string `json:"message"`
	}{
		Id:      id,
		Message: "User is created",
	}

	jsonResponse, err := json.Marshal(output)
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (ar *AuthRouter) UserIdentity(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		userId, err := lib.ParseToken(headerParts[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx := r.Context()
		req := r.WithContext(context.WithValue(ctx, userCtx, userId))
		*r = *req

		logrus.Printf("hello from middleware")
		next(w, r)
	}
}
