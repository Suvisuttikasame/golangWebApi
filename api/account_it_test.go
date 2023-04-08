//go:build integration

package api

import (
	"database/sql"
	"encoding/json"
	"goApp/authentication"
	db "goApp/db/sqlc"
	"goApp/middleware"
	"goApp/util"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestITCreateAccountAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := authentication.Body{
		Id:       uuid.New(),
		Username: "tommy",
		Email:    "test@mail.com",
	}
	config, err := util.NewConfig("..")
	if err != nil {
		t.Fatal(err)
	}
	p, err := authentication.NewPasetoToken([]byte(config.SecretKey))
	if err != nil {
		t.Fatal(err)
	}
	token, err := p.CreateToken(body)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/accounts", strings.NewReader(`{"currency":"THB"}`))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()

	psql := "host=localhost port=5432 user=postgres password=postgres dbname=simple_bank sslmode=disable"
	conn, err := sql.Open("postgres", psql)

	assert.Nil(t, err)
	_, err = db.SeedUser(conn, db.CreateUserParams{
		Username: "tommy",
		Email:    "tommy@mail.com",
		Password: "test",
	})
	assert.Nil(t, err)

	var a db.Account

	s := db.NewStore(conn)
	server := &Server{store: s}

	r := gin.New()
	r.Use(middleware.Middleware(p))
	r.POST("/accounts", server.CreateAccount)

	r.ServeHTTP(rec, req)
	t.Log(rec)
	assert.Equal(t, http.StatusOK, rec.Code)

	err = json.Unmarshal(rec.Body.Bytes(), &a)
	assert.Nil(t, err)
	assert.Equal(t, "tommy", a.Owner)
	assert.Equal(t, "THB", a.Currency)
}
