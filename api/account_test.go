//go:build unit

package api

import (
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

func TestCreateAccountAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)
	req := httptest.NewRequest(http.MethodPost, "/accounts", strings.NewReader(`{
		"owner":"john",
		"currency":"THB"
	}`))
	req.Header.Set("Content-type", "application/json")
	rec := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rec)

	c.Request = req

	conn, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatal(err)
	}
	ti := util.MockTime()

	mock.ExpectQuery(`-- name: CreateAccount :one
	INSERT INTO accounts (
	  owner, balance, currency
	) VALUES (
	  $1, $2, $3
	)
	RETURNING id, owner, balance, currency, created_at`).
		WithArgs("john", 0, "THB").
		WillReturnRows(sqlmock.NewRows([]string{"id", "owner", "balance", "currency", "created_at"}).
			AddRow(1, "john", 0, "THB", ti))

	s := db.NewStore(conn)

	server := &Server{store: s}

	server.CreateAccount(c)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, `{"id":1,"owner":"john","balance":0,"currency":"THB","created_at":"2023-01-01T12:45:40.000000003Z"}`, rec.Body.String())

}
