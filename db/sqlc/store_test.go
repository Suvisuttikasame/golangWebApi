package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestStore(t *testing.T) {
	md := TransferParams{
		ToAccountID:   2,
		FromAccountID: 1,
		Amount:        10,
	}
	ti := time.Now()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatal(err)
	}

	mock.ExpectBegin()

	mock.ExpectQuery(`-- name: CreateTransfer :one
	INSERT INTO transfers (
	  from_account_id, to_account_id, amount
	) VALUES (
	  $1, $2, $3
	)
	RETURNING id, from_account_id, to_account_id, amount, created_at`).
		WithArgs(md.FromAccountID, md.ToAccountID, md.Amount).
		WillReturnRows(sqlmock.NewRows([]string{"id", "from_account_id", "to_account_id", "amount", "created_at"}).AddRow(1, md.FromAccountID, md.ToAccountID, md.Amount, ti))

	mock.ExpectQuery(`-- name: CreateEntry :one
	INSERT INTO entries (
	  account_id, amount
	) VALUES (
	  $1, $2
	)
	RETURNING id, account_id, amount, created_at`).
		WithArgs(md.ToAccountID, md.Amount).
		WillReturnRows(sqlmock.NewRows([]string{"id", "account_id", "amount", "created_at"}).AddRow(1, md.ToAccountID, md.Amount, ti))

	mock.ExpectQuery(`-- name: CreateEntry :one
	INSERT INTO entries (
	  account_id, amount
	) VALUES (
	  $1, $2
	)
	RETURNING id, account_id, amount, created_at`).
		WithArgs(md.FromAccountID, -md.Amount).
		WillReturnRows(sqlmock.NewRows([]string{"id", "account_id", "amount", "created_at"}).AddRow(2, md.FromAccountID, -md.Amount, ti))

	mock.ExpectQuery(`-- name: GetAccountForUpdate :one
	SELECT id, owner, balance, currency, created_at FROM accounts
	WHERE id = $1 LIMIT 1
	FOR NO KEY UPDATE`).
		WithArgs(md.ToAccountID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "owner", "balance", "currency", "created_at"}).
			AddRow(md.ToAccountID, "logan", 100.00, "THB", ti))

	mock.ExpectQuery(`-- name: GetAccountForUpdate :one
	SELECT id, owner, balance, currency, created_at FROM accounts
	WHERE id = $1 LIMIT 1
	FOR NO KEY UPDATE`).
		WithArgs(md.FromAccountID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "owner", "balance", "currency", "created_at"}).
			AddRow(md.FromAccountID, "salah", 100.00, "THB", ti))

	mock.ExpectQuery(`-- name: UpdateAccount :one
	UPDATE accounts
	set owner = COALESCE(NULLIF($1, ''), owner) ,
	balance = $2,
	currency = COALESCE(NULLIF($3, ''), currency)
	WHERE id = $4
	RETURNING id, owner, balance, currency, created_at`).
		WithArgs(sqlmock.AnyArg(), 100-md.Amount, sqlmock.AnyArg(), md.FromAccountID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "owner", "balance", "currency", "created_at"}).
			AddRow(md.FromAccountID, "salah", 100-md.Amount, "THB", ti))

	mock.ExpectQuery(`-- name: UpdateAccount :one
	UPDATE accounts
	set owner = COALESCE(NULLIF($1, ''), owner) ,
	balance = $2,
	currency = COALESCE(NULLIF($3, ''), currency)
	WHERE id = $4
	RETURNING id, owner, balance, currency, created_at`).
		WithArgs(sqlmock.AnyArg(), 100+md.Amount, sqlmock.AnyArg(), md.ToAccountID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "owner", "balance", "currency", "created_at"}).
			AddRow(md.ToAccountID, "logan", 100+md.Amount, "THB", ti))

	//unit test update account
	mock.ExpectCommit()

	s := NewStore(db)

	result, err := s.TransferTx(context.Background(), md)

	assert.Nil(t, err)
	assert.Equal(t, result.FromEntry.AccountID, md.FromAccountID)
	assert.Equal(t, result.ToEntry.AccountID, md.ToAccountID)
	assert.Equal(t, result.Transfer.Amount, md.Amount)

}
