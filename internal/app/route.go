package app

import (
	_ "github.com/ereminiu/filmoteka/docs"
	v1 "github.com/ereminiu/filmoteka/internal/controller"
	"github.com/ereminiu/filmoteka/internal/db"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

func NewRouter(repos *db.Repositories, mig *db.Migrator) *http.ServeMux {
	router := http.NewServeMux()

	migrateRouter := v1.NewMigrateRouter(mig)
	authRouter := v1.NewAuthRouter(repos.Authorization)
	movieRouter := v1.NewMovieRouter(repos.Movie)
	actorRouter := v1.NewActorRouter(repos.Actor)

	// Migrations
	router.HandleFunc("POST /migrate-up", migrateRouter.MigrateUp)
	router.HandleFunc("POST /migrate-down", migrateRouter.MigrateDown)
	router.HandleFunc("POST /migrate-force", migrateRouter.MigrateForce)

	// Authorization
	router.HandleFunc("POST /sign-up", authRouter.CreateUser)
	router.HandleFunc("POST /sign-in", authRouter.GenerateToken)

	// User group
	router.HandleFunc("GET /actor-list", actorRouter.GetAllActors)
	router.HandleFunc("POST /search-movies", movieRouter.SearchMovie)
	router.HandleFunc("POST /movie-list", movieRouter.GetAllMovies)

	// Admin group
	router.HandleFunc("POST /add-actor", authRouter.UserIdentity(actorRouter.AddActor))
	router.HandleFunc("POST /add-actor-to-movie", authRouter.UserIdentity(movieRouter.AddActorToMovie))
	router.HandleFunc("PUT /change-actor-field", authRouter.UserIdentity(actorRouter.ChangeField))
	router.HandleFunc("DELETE /delete-actor", authRouter.UserIdentity(actorRouter.DeleteActor))
	router.HandleFunc("POST /add-movie", authRouter.UserIdentity(movieRouter.AddMovie))
	router.HandleFunc("DELETE /delete-movie", authRouter.UserIdentity(movieRouter.DeleteMovie))
	router.HandleFunc("DELETE /delete-movie-field", authRouter.UserIdentity(movieRouter.DeleteField))
	router.HandleFunc("PUT /change-movie-field", authRouter.UserIdentity(movieRouter.ChangeField))

	router.Handle("/swagger/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:3000/swagger/doc.json")))

	return router
}
