//go:build unit

package api

import (
	"encoding/json"
	db "goApp/db/sqlc"
	"goApp/util"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateTransferAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)

	req := httptest.NewRequest(http.MethodPost, "/transfers", strings.NewReader(`{"from_account_id":1,"to_account_id":2,"amount":30}`))
	req.Header.Set("Content-type", "application/json")

	rec := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rec)
	c.Request = req

	ti := util.MockTime()

	conn, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)

	//mock get acc 1
	mock.ExpectQuery(`-- name: GetAccount :one
	SELECT id, owner, balance, currency, created_at FROM accounts
	WHERE id = $1 LIMIT 1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "owner", "balance", "currency", "created_at"}).AddRow(1, "john", 100, "THB", ti))

	//mock get acc 2
	mock.ExpectQuery(`-- name: GetAccount :one
	SELECT id, owner, balance, currency, created_at FROM accounts
	WHERE id = $1 LIMIT 1`).
		WithArgs(2).
		WillReturnRows(sqlmock.NewRows([]string{"id", "owner", "balance", "currency", "created_at"}).AddRow(2, "david", 120, "THB", ti))

	//start tx
	mock.ExpectBegin()

	mock.ExpectQuery(`-- name: CreateTransfer :one
	INSERT INTO transfers (
	from_account_id, to_account_id, amount
	) VALUES (
	$1, $2, $3
	)
	RETURNING id, from_account_id, to_account_id, amount, created_at`).
		WithArgs(1, 2, 30).
		WillReturnRows(sqlmock.NewRows([]string{"id", "from_account_id", "to_account_id", "amount", "created_at"}).AddRow(1, 1, 2, 30, ti))

	mock.ExpectQuery(`-- name: CreateEntry :one
	INSERT INTO entries (
	account_id, amount
	) VALUES (
	$1, $2
	)
	RETURNING id, account_id, amount, created_at`).
		WithArgs(2, 30).
		WillReturnRows(sqlmock.NewRows([]string{"id", "account_id", "amount", "created_at"}).AddRow(1, 2, 30, ti))

	mock.ExpectQuery(`-- name: CreateEntry :one
	INSERT INTO entries (
	account_id, amount
	) VALUES (
	$1, $2
	)
	RETURNING id, account_id, amount, created_at`).
		WithArgs(1, -30).
		WillReturnRows(sqlmock.NewRows([]string{"id", "account_id", "amount", "created_at"}).AddRow(2, 1, -30, ti))

	mock.ExpectQuery(`-- name: GetAccountForUpdate :one
	SELECT id, owner, balance, currency, created_at FROM accounts
	WHERE id = $1 LIMIT 1
	FOR NO KEY UPDATE`).
		WithArgs(2).
		WillReturnRows(sqlmock.NewRows([]string{"id", "owner", "balance", "currency", "created_at"}).
			AddRow(2, "david", 120, "THB", ti))

	mock.ExpectQuery(`-- name: GetAccountForUpdate :one
	SELECT id, owner, balance, currency, created_at FROM accounts
	WHERE id = $1 LIMIT 1
	FOR NO KEY UPDATE`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "owner", "balance", "currency", "created_at"}).
			AddRow(1, "john", 100, "THB", ti))

	mock.ExpectQuery(`-- name: UpdateAccount :one
	UPDATE accounts
	  set owner = COALESCE(NULLIF($1, ''), owner) ,
	  balance = $2,
	  currency = COALESCE(NULLIF($3, ''), currency)
	WHERE id = $4
	RETURNING id, owner, balance, currency, created_at`).
		WithArgs(sqlmock.AnyArg(), 100-30, sqlmock.AnyArg(), 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "owner", "balance", "currency", "created_at"}).
			AddRow(1, "john", 100-30, "THB", ti))

	mock.ExpectQuery(`-- name: UpdateAccount :one
	UPDATE accounts
	  set owner = COALESCE(NULLIF($1, ''), owner) ,
	  balance = $2,
	  currency = COALESCE(NULLIF($3, ''), currency)
	WHERE id = $4
	RETURNING id, owner, balance, currency, created_at`).
		WithArgs(sqlmock.AnyArg(), 120+30, sqlmock.AnyArg(), 2).
		WillReturnRows(sqlmock.NewRows([]string{"id", "owner", "balance", "currency", "created_at"}).
			AddRow(2, "david", 120+30, "THB", ti))

	mock.ExpectCommit()
	s := db.NewStore(conn)

	server := &Server{store: s}

	server.CreateTransfer(c)

	var r db.TransferResult

	err = json.Unmarshal(rec.Body.Bytes(), &r)

	assert.Nil(t, err)
	assert.Equal(t, int64(70), r.FromAccount.Balance)
	assert.Equal(t, int64(150), r.ToAccount.Balance)

}
