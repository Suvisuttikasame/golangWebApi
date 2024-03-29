//go:build integration

package db

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestItStore(t *testing.T) {
	db, err := initTestDB()
	if err != nil {
		t.Fatal(err)
	}
	_, err = SeedUser(db, CreateUserParams{
		Username: "p1",
		Email:    "p1@mail.com",
		Password: "test",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = SeedUser(db, CreateUserParams{
		Username: "p2",
		Email:    "p2@mail.com",
		Password: "test",
	})
	if err != nil {
		t.Fatal(err)
	}

	a1, err := SeedAccount(db, CreateAccountParams{
		Owner:    "p1",
		Balance:  100,
		Currency: "THB",
	})
	if err != nil {
		t.Fatal(err)
	}

	a2, err := SeedAccount(db, CreateAccountParams{
		Owner:    "p2",
		Balance:  100,
		Currency: "THB",
	})
	if err != nil {
		t.Fatal(err)
	}

	s := NewStore(db)
	r, err := s.TransferTx(context.Background(), TransferParams{
		FromAccountID: a1.ID,
		ToAccountID:   a2.ID,
		Amount:        20,
	})

	assert.Nil(t, err)
	t.Log(r)

}

func initTestDB() (*sql.DB, error) {
	psqlInfo := "host=localhost port=5432 user=postgres password=postgres dbname=simple_bank sslmode=disable"
	db, err := sql.Open("postgres", psqlInfo)

	return db, err

}
