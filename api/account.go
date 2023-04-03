package api

import (
	"database/sql"
	"fmt"
	db "goApp/db/sqlc"
	"goApp/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=THB USD"`
}

type GetOneAccountURIParam struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type GetAllAccountQyeryParam struct {
	PageId   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=1,max=5"`
}

func (sv *Server) CreateAccount(ctx *gin.Context) {
	var req CreateAccountRequest
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	cap := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}

	a, err := sv.store.CreateAccount(ctx, cap)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, a)

}

func (sv *Server) GetOneAccount(ctx *gin.Context) {
	var id GetOneAccountURIParam

	err := ctx.ShouldBindUri(&id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	a, err := sv.store.GetAccount(ctx, id.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, util.ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, a)
}

func (sv *Server) GetAllAccount(ctx *gin.Context) {
	var p GetAllAccountQyeryParam

	err := ctx.ShouldBindQuery(&p)
	data, _ := ctx.Get("authorization_key")
	fmt.Println("breakdown data", data)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	la := db.ListAccountParams{
		Limit:  p.PageSize,
		Offset: (p.PageId - 1) * p.PageSize,
	}

	l, err := sv.store.ListAccount(ctx, la)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, l)
}
