package app

import (
	"encoding/json"
	_ "github.com/ereminiu/filmoteka/docs"
	v1 "github.com/ereminiu/filmoteka/internal/controller"
	"github.com/ereminiu/filmoteka/internal/db"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

func NewRouter(repos *db.Repositories) *http.ServeMux {
	router := http.NewServeMux()

	movieRouter := v1.NewMovieRouter(repos.Movie)

	// User group
	router.HandleFunc("POST /actor-list", empty)
	router.HandleFunc("POST /search-movie", empty)
	router.HandleFunc("POST /movie-list", movieRouter.GetAllMovies)

	// Admin group
	router.HandleFunc("POST /add-actor", empty)
	router.HandleFunc("PUT /change-actor-field", empty)
	router.HandleFunc("DELETE /delete-actor", empty)
	router.HandleFunc("DELETE /delete-actor-field", empty)
	router.HandleFunc("POST /add-movie", movieRouter.AddMovie)
	router.HandleFunc("DELETE /delete-movie", empty)
	router.HandleFunc("DELETE /delete-movie-field", empty)
	router.HandleFunc("PUT /change-movie-field", movieRouter.ChangeField)

	router.HandleFunc("POST /hello", getReverseName)
	router.Handle("/swagger/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:3000/swagger/doc.json")))

	return router
}

type userInput struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// @Summary Get Reversed Name
// @Description reversing the given name
// @Accept json
// @Produce json
// @Param input body userInput true "username and age"
// @Success 200 {object} userInput
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /hello [post]
func getReverseName(w http.ResponseWriter, r *http.Request) {
	var input userInput
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		logrus.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	output := userInput{Name: v1.ReverseName(input.Name), Age: input.Age + 1}
	jsonResponse, err := json.Marshal(output)
	if err != nil {
		logrus.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func empty(w http.ResponseWriter, r *http.Request) {
}
