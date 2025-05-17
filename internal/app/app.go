package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mitchan/go_gym_bro/internal/api"
	"github.com/mitchan/go_gym_bro/internal/store"
	"github.com/mitchan/go_gym_bro/migrations"
)

type Application struct {
	DB             *sql.DB
	Logger         *log.Logger
	TokenHandler   *api.TokenHandler
	UserHandler    *api.UserHandler
	WorkoutHandler *api.WorkoutHandler
}

func NewApplication() (*Application, error) {
	db, err := store.Open()
	if err != nil {
		return nil, err
	}

	err = store.MigrateFs(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// store
	tokenStore := store.NewTokenStore(db)
	userStore := store.NewPostgresUserStore(db)
	workoutStore := store.NewPostgresWorkoutStore(db)

	// handlers
	tokenHandler := api.NewTokenStoreHandler(tokenStore, userStore, logger)
	userHandler := api.NewUserHandler(userStore, logger)
	workoutHandler := api.NewWorkoutHandler(workoutStore, logger)

	app := Application{
		DB:             db,
		Logger:         logger,
		TokenHandler:   tokenHandler,
		UserHandler:    userHandler,
		WorkoutHandler: workoutHandler,
	}
	return &app, nil
}

func (a Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Status is available\n")
}
