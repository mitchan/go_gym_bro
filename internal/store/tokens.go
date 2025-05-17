package store

import (
	"database/sql"
	"time"

	"github.com/mitchan/go_gym_bro/internal/tokens"
)

type PostgresTokenStore struct {
	db *sql.DB
}

type TokenStore interface {
	Insert(token *tokens.Token) error
	CreateNewToken(userID int, ttl time.Duration, scope string) (*tokens.Token, error)
	DeleteAllTokensForUser(userID int, scope string) error
}

func NewTokenStore(db *sql.DB) *PostgresTokenStore {
	return &PostgresTokenStore{
		db,
	}
}

func (t *PostgresTokenStore) CreateNewToken(userID int, ttl time.Duration, scope string) (*tokens.Token, error) {
	token, err := tokens.GenerateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = t.Insert(token)
	return token, err
}

func (t *PostgresTokenStore) Insert(token *tokens.Token) error {
	query := `
  insert into tokens (hash, user_id, expiry, scope)
	values ($1, $2, $3, $4)
	`

	_, err := t.db.Exec(query, token.Hash, token.UserID, token.Expiry, token.Scope)
	return err
}

func (t *PostgresTokenStore) DeleteAllTokensForUser(userID int, scope string) error {
	query := `
	delete from tokens
	where user_id = $1 and scope = $2
	`

	_, err := t.db.Exec(query, userID, scope)
	return err
}
