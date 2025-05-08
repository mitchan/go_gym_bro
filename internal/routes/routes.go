package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/mitchan/go_gym_bro/internal/app"
)

func SetupRoutes(a *app.Application) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/health", a.HealthCheck)
	r.Get("/workouts/{id}", a.WorkoutHandler.HandleGetWorkoutById)

	r.Post("/workouts", a.WorkoutHandler.HandleCreate)
	r.Put("/workouts/{id}", a.WorkoutHandler.HandleUpdateWorkoutById)
	r.Delete("/workouts/{id}", a.WorkoutHandler.HandleDeleteWorkout)
	return r
}
