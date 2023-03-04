//go:build unit

package db

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateAccount(t *testing.T) {
	//arrangement
	ap := CreateAccountParams{
		Owner:    "jonathan",
		Balance:  100.00,
		Currency: "THB",
	}
	ti := time.Now()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	//-- name: CreateAccount :one INSERT INTO accounts ( owner, balance, currency ) VALUES ( $1, $2, $3 ) RETURNING id, owner, balance, currency, created_at
	mock.ExpectQuery(`-- name: CreateAccount :one INSERT INTO accounts ( owner, balance, currency ) VALUES ( $1, $2, $3 ) RETURNING id, owner, balance, currency, created_at`).
		WithArgs(ap.Owner, ap.Balance, ap.Currency).
		WillReturnRows(sqlmock.NewRows([]string{"id", "owner", "balance", "currency", "created_at"}).AddRow(1, ap.Owner, ap.Balance, ap.Currency, ti))

	//action
	q := New(db)
	a, err := q.CreateAccount(context.Background(), ap)

	//assertion
	assert.Nil(t, err)
	assert.Equal(t, int64(1), a.ID)
	assert.Equal(t, ap.Balance, a.Balance)
	assert.Equal(t, ap.Owner, a.Owner)
	assert.Equal(t, ap.Currency, a.Currency)
	assert.Equal(t, ti, a.CreatedAt)
}
func TestGetAllAccount(t *testing.T) {

}
func TestGetOneAccount(t *testing.T) {

}
func TestUpdateAccount(t *testing.T) {

}
func TestDeleteAccount(t *testing.T) {

}
