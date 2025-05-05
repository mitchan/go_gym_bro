package store

import "database/sql"

type Workout struct {
	ID              int            `json:"id"`
	Title           string         `json:"title"`
	Description     string         `json:"description"`
	DurationMinutes int            `json:"duration_minutes"`
	CaloriesBurned  int            `json:"calories_burned"`
	Entries         []WorkoutEntry `json:"entries"`
}

type WorkoutEntry struct {
	ID              int      `json:"id"`
	ExerciseName    string   `json:"exercise_name"`
	Sets            int      `json:"sets"`
	Reps            *int     `json:"reps"`
	DurationSeconds *int     `json:"duration_seconds"`
	Weight          *float64 `json:"weight"`
	Notes           string   `json:"notes"`
	OrderIndex      int      `json:"order_index"`
}

type PostgresWorkoutStore struct {
	db *sql.DB
}

func NewPostgresWorkoutStore(db *sql.DB) *PostgresWorkoutStore {
	return &PostgresWorkoutStore{
		db: db,
	}
}

type WorkoutStore interface {
	CreateWorkout(*Workout) (*Workout, error)
	GetWorkoutById(id int64) (*Workout, error)
	UpdateWorkout(*Workout) error
}

func (pg *PostgresWorkoutStore) CreateWorkout(w *Workout) (*Workout, error) {
	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `insert into workouts (title, description, duration_minutes, calories_burned)
	values ($1, $2, $3, $4)
	returning id`

	err = tx.QueryRow(query, w.Title, w.Description, w.DurationMinutes, w.CaloriesBurned).Scan(&w.ID)
	if err != nil {
		return nil, err
	}

	for _, entry := range w.Entries {
		query := `insert into workout_entries (workout_id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index)
	values ($1, $2, $3, $4, $5, $6, $7, $8)
	returning id`
		err = tx.QueryRow(query, w.ID, entry.ExerciseName, entry.Sets, entry.Reps, entry.DurationSeconds, entry.Weight, entry.Notes, entry.OrderIndex).Scan(&entry.ID)
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (pg *PostgresWorkoutStore) GetWorkoutById(id int64) (*Workout, error) {
	w := &Workout{}
	query := `select id, title, description, duration_minutes, calories_burned from workouts
	where id = $1`
	err := pg.db.QueryRow(query, id).Scan(&w.ID, &w.Title, &w.Description, &w.DurationMinutes, &w.CaloriesBurned)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	entryQuery := `select id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index
	from workout_entries
	where workout_id = $1
	order by order_index`

	rows, err := pg.db.Query(entryQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var entry WorkoutEntry
		err := rows.Scan(
			&entry.ID,
			&entry.ExerciseName,
			&entry.Sets,
			&entry.Reps,
			&entry.DurationSeconds,
			&entry.Weight,
			&entry.Notes,
			&entry.OrderIndex,
		)
		if err != nil {
			return nil, err
		}
		w.Entries = append(w.Entries, entry)
	}

	return w, nil
}

func (pg *PostgresWorkoutStore) UpdateWorkout(w *Workout) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
	update workouts
	set title = $1, description = $2, duration_minutes = $3, calories_burned = $4
	where id = $5
	`
	result, err := tx.Exec(query, w.Title, w.Description, w.DurationMinutes, w.CaloriesBurned, w.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	_, err = tx.Exec(`delete from workout_entries where workout_id = $1`, w.ID)
	if err != nil {
		return err
	}
	for _, entry := range w.Entries {
		query := `
    INSERT INTO workout_entries (workout_id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
  `

		_, err := tx.Exec(query,
			w.ID,
			entry.ExerciseName,
			entry.Sets,
			entry.Reps,
			entry.DurationSeconds,
			entry.Weight,
			entry.Notes,
			entry.OrderIndex,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
