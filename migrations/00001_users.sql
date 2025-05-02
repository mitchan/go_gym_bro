-- +goose Up
-- +goose StatementBegin
create table if not exists users (
  id bigserial primary key,
  username varchar(50) unique not null,
  email varchar(255) unique not null,
  password_hash varchar(255) not null,
  bio text,
  created_at timestamp with time zone default current_timestamp,
  updated_at timestamp with time zone default current_timestamp
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users;
-- +goose StatementEnd
