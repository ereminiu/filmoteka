package controller

import (
	"encoding/json"
	"github.com/ereminiu/filmoteka/internal/controller/lib"
	"github.com/ereminiu/filmoteka/internal/db"
	"github.com/sirupsen/logrus"
	"net/http"
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
		// TODO: add verbose error
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
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