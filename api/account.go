package api

import (
	db "goApp/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=THB"`
}

func (sv *Server) CreateAccount(ctx *gin.Context) {
	var req CreateAccountRequest
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}
	cap := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}

	a, err := sv.store.CreateAccount(ctx, cap)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, a)

}
