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
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatal(err)
	}
	ti := time.Now()

	mock.ExpectQuery(`-- name: ListAccount :many
	SELECT id, owner, balance, currency, created_at FROM accounts
	ORDER BY id
	LIMIT $1
	OFFSET $2`).
		WithArgs(2, 0).
		WillReturnRows(sqlmock.NewRows([]string{"id", "owner", "balance", "currency", "created_at"}).
			AddRow(1, "ronaldo", 100.00, "THB", ti).
			AddRow(2, "messi", 200.00, "THB", ti))
	lp := ListAccountParams{
		Limit:  2,
		Offset: 0,
	}

	q := New(db)
	items, err := q.ListAccount(context.Background(), lp)

	assert.Nil(t, err)
	assert.Greater(t, len(items), 0)

}
func TestGetOneAccount(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatal(err)
	}
	ti := time.Now()

	mock.ExpectQuery(`-- name: GetAccount :one
	SELECT id, owner, balance, currency, created_at FROM accounts
	WHERE id = $1 LIMIT 1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "owner", "balance", "currency", "created_at"}).
			AddRow(1, "salah", 100.00, "THB", ti))

	q := New(db)
	a, err := q.GetAccount(context.Background(), 1)

	assert.Nil(t, err)
	assert.Equal(t, a.ID, int64(1))
	assert.Equal(t, a.Owner, "salah")
	assert.Equal(t, a.Balance, int64(100))
	assert.Equal(t, a.Currency, "THB")
	assert.Equal(t, a.CreatedAt, ti)

}
func TestUpdateAccount(t *testing.T) {
	up := UpdateAccountParams{
		ID:       1,
		Owner:    "hello",
		Balance:  30,
		Currency: "THB",
	}
	ti := time.Now()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatal(err)
	}

	mock.ExpectQuery(`-- name: UpdateAccount :one
				UPDATE accounts
				set owner = COALESCE(NULLIF($1, ''), owner) ,
				balance = $2,
				currency = COALESCE(NULLIF($3, ''), currency)
				WHERE id = $4
				RETURNING id, owner, balance, currency, created_at`).
		WithArgs(up.Owner, up.Balance, up.Currency, up.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "owner", "balance", "currency", "created_at"}).
			AddRow(up.ID, up.Owner, up.Balance, up.Currency, ti))

	q := New(db)
	i, err := q.UpdateAccount(context.Background(), up)

	assert.Nil(t, err)
	assert.Equal(t, i.ID, up.ID)
	assert.Equal(t, i.Owner, up.Owner)
	assert.Equal(t, i.Balance, up.Balance)
	assert.Equal(t, i.Currency, up.Currency)
	assert.Equal(t, i.CreatedAt, ti)

}
func TestDeleteAccount(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatal(err)
	}
	ti := time.Now()

	mock.ExpectQuery(`-- name: DeleteAccount :one
	DELETE FROM accounts
	WHERE id = $1
	RETURNING id, owner, balance, currency, created_at`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "owner", "balance", "currency", "created_at"}).
			AddRow(1, "salah", 100.00, "THB", ti))

	q := New(db)
	a, err := q.DeleteAccount(context.Background(), 1)

	assert.Nil(t, err)
	assert.Equal(t, a.ID, int64(1))
	assert.Equal(t, a.Owner, "salah")
	assert.Equal(t, a.Balance, int64(100))
	assert.Equal(t, a.Currency, "THB")
	assert.Equal(t, a.CreatedAt, ti)

}
