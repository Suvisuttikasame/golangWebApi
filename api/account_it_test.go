//go:build integration

package api

import (
	"database/sql"
	"encoding/json"
	db "goApp/db/sqlc"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestITCreateAccountAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)

	req := httptest.NewRequest(http.MethodPost, "/accounts", strings.NewReader(`{"owner":"tommy","currency":"THB"}`))
	req.Header.Set("Content-type", "application/json")

	rec := httptest.NewRecorder()

	psql := "host=localhost port=5432 user=postgres password=postgres dbname=simple_bank sslmode=disable"
	conn, err := sql.Open("postgres", psql)

	assert.Nil(t, err)

	var a db.Account

	s := db.NewStore(conn)
	server := &Server{store: s}

	r := gin.New()
	r.POST("/accounts", server.CreateAccount)

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	err = json.Unmarshal(rec.Body.Bytes(), &a)
	assert.Nil(t, err)
	assert.Equal(t, "tommy", a.Owner)
	assert.Equal(t, "THB", a.Currency)
}
