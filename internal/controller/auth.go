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

// @Summary SignIn
// @Tags auth
// @Description login
// @ID login
// @Accept  json
// @Produce  json
// @Param input body createUserInput true "credentials"
// @Success 200 {string} string "token"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal error"
// @Failure default {string} string "error"
// @Router /sign-in [post]
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

	userRole := lib.ParseUserRole(input.Username)
	token, err := lib.GenerateToken(userId, userRole)
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

	logrus.Println("Success")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// @Summary SignUp
// @Tags auth
// @Description create account
// @Id create-account
// @Accept json
// @Produce json
// @Param input body createUserInput true "name username password"
// @Success 200 {object} outputWithId
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /sign-up [post]
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

	output := outputWithId{
		Id:      id,
		Message: "User is created",
	}

	logrus.Println("Success")
	logrus.Println(output)

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
		logrus.Printf("check authorization")
		header := r.Header.Get("Authorization")
		if header == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 {
			logrus.Error("invalid auth-header")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		userId, userRole, err := lib.ParseToken(headerParts[1])
		if err != nil {
			logrus.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if userRole != "admin" {
			logrus.Error("not enough rights")
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		ctx := r.Context()
		req := r.WithContext(context.WithValue(ctx, userCtx, userId))
		*r = *req

		logrus.Printf("user is authorized as admin")
		next(w, r)
	}
}
