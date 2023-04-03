package middleware

import (
	"goApp/authentication"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestAuthMiddleware(t *testing.T) {
	t.Run("return invalid authorize header", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		req := httptest.NewRequest(http.MethodGet, "/accounts", nil)
		req.Header.Set("Content-type", "application/json")

		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Request = req

		// body := authentication.Body{
		// 	Id:       uuid.New(),
		// 	Username: "test",
		// 	Email:    "test@mail.com",
		// }
		key := "asdfgjhtuasdfgjhtuasdfgjhtuasdfg"
		p, err := authentication.NewPasetoToken([]byte(key))
		if err != nil {
			t.Fatal(err)
		}

		r := gin.New()
		r.Use(Middleware(p))
		r.GET("/accounts", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, "hello world")
			return
		})

		r.ServeHTTP(rec, req)

		t.Log(rec.Body)

	})

	t.Run("return code 200 success", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		body := authentication.Body{
			Id:       uuid.New(),
			Username: "test",
			Email:    "test@mail.com",
		}
		key := "asdfgjhtuasdfgjhtuasdfgjhtuasdfg"
		p, err := authentication.NewPasetoToken([]byte(key))
		if err != nil {
			t.Fatal(err)
		}
		tok, err := p.CreateToken(body)
		if err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest(http.MethodGet, "/accounts", nil)
		req.Header.Set("Content-type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tok)

		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Request = req

		r := gin.New()
		r.Use(Middleware(p))
		r.GET("/accounts", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, "hello world")
			return
		})

		r.ServeHTTP(rec, req)

		t.Log(rec.Body)

	})

}
