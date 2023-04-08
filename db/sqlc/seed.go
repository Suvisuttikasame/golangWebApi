package db

import (
	"context"
	"database/sql"
)

func SeedAccount(db *sql.DB, arg CreateAccountParams) (Account, error) {
	q := New(db)
	a, err := q.CreateAccount(context.Background(), arg)
	if err != nil {
		return a, err
	}
	return a, nil
}

func SeedUser(db *sql.DB, arg CreateUserParams) (User, error) {
	q := New(db)
	u, err := q.CreateUser(context.Background(), arg)
	if err != nil {
		return u, err
	}
	return u, nil
}
