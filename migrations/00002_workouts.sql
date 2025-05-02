-- +goose Up
-- +goose StatementBegin
create table if not exists workouts (
  id bigserial primary key,
  -- user_id
  title varchar(255) not null,
  description text,
  duration_minutes integer not null,
  calories_burned integer,
  created_at timestamp with time zone default current_timestamp,
  updated_at timestamp with time zone default current_timestamp
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table workouts;
-- +goose StatementEnd
