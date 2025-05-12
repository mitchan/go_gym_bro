package store

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres port=5433 sslmode=disable")
	if err != nil {
		t.Fatalf("opening test DB: %v", err)
	}

	err = Migrate(db, "../../migrations/")
	if err != nil {
		t.Fatalf("migrating test DB: %v", err)
	}

	_, err = db.Exec("truncate workouts, workout_entries cascade")
	if err != nil {
		t.Fatalf("truncating tables: %v", err)
	}

	return db
}

func TestCreateWorkout(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewPostgresWorkoutStore(db)

	tests := []struct {
		name      string
		workout   *Workout
		wantError bool
	}{
		{
			name: "valid workout",
			workout: &Workout{
				Title:           "push day",
				Description:     "Upper body day",
				DurationMinutes: 60,
				CaloriesBurned:  200,
				Entries: []WorkoutEntry{
					{
						ExerciseName: "Bench press",
						Sets:         3,
						Reps:         intPtr(10),
						Weight:       floatPtr(135.5),
						Notes:        "warm up properly",
						OrderIndex:   1,
					},
				},
			},
			wantError: false,
		},
		{
			name: "workout with invalid entries",
			workout: &Workout{
				Title:           "full body",
				Description:     "complete workout",
				DurationMinutes: 90,
				CaloriesBurned:  500,
				Entries: []WorkoutEntry{
					{
						ExerciseName: "plank",
						Sets:         3,
						Reps:         intPtr(60),
						Notes:        "keep form",
						OrderIndex:   1,
					},
					{
						ExerciseName:    "squats",
						Sets:            4,
						Reps:            intPtr(12),
						DurationSeconds: intPtr(60),
						Weight:          floatPtr(185.0),
						Notes:           "keep form",
						OrderIndex:      2,
					},
				},
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createdWorkout, err := store.CreateWorkout(tt.workout)
			if tt.wantError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.workout.Title, createdWorkout.Title)
			assert.Equal(t, tt.workout.Description, createdWorkout.Description)
			assert.Equal(t, tt.workout.DurationMinutes, createdWorkout.DurationMinutes)

			retrieved, err := store.GetWorkoutById(int64(createdWorkout.ID))
			require.NoError(t, err)

			assert.Equal(t, createdWorkout.ID, retrieved.ID)
			assert.Equal(t, len(tt.workout.Entries), len(retrieved.Entries))

			for i, entry := range retrieved.Entries {
				assert.Equal(t, tt.workout.Entries[i].ExerciseName, entry.ExerciseName)
				assert.Equal(t, tt.workout.Entries[i].Sets, entry.Sets)
				assert.Equal(t, tt.workout.Entries[i].OrderIndex, entry.OrderIndex)
			}
		})
	}
}

func intPtr(i int) *int {
	return &i
}

func floatPtr(f float64) *float64 {
	return &f
}
