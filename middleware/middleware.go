package middleware

import (
	"errors"
	"goApp/authentication"
	"goApp/util"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Middleware(tk authentication.AuthenPaseto) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorization := ctx.Request.Header.Get("Authorization")
		if authorization == "" {
			ctx.JSON(http.StatusUnauthorized, util.ErrorResponse(errors.New("invalid authorize header")))
			ctx.Abort()
			return
		}
		split := strings.SplitN(authorization, " ", 2)
		if split[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, util.ErrorResponse(errors.New("invalid authorize header")))
			ctx.Abort()
			return
		}
		bd, err := tk.Verification(split[1])
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, util.ErrorResponse(err))
			ctx.Abort()
			return
		}
		ctx.Set("authorization_key", bd)
		ctx.Next()
	}
}
