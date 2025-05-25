package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/mitchan/go_gym_bro/internal/app"
)

func SetupRoutes(a *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(a.Middleware.Authenticate)

		r.Get("/workouts/{id}", a.Middleware.RequireUser(a.WorkoutHandler.HandleGetWorkoutById))
		r.Post("/workouts", a.Middleware.RequireUser(a.WorkoutHandler.HandleCreate))
		r.Put("/workouts/{id}", a.Middleware.RequireUser(a.WorkoutHandler.HandleUpdateWorkoutById))
		r.Delete("/workouts/{id}", a.Middleware.RequireUser(a.WorkoutHandler.HandleDeleteWorkout))
	})

	r.Get("/health", a.HealthCheck)
	r.Post("/users", a.UserHandler.HandleRegisterUser)
	r.Post("/tokens/authentication", a.TokenHandler.HandleCreateToken)
	return r
}
