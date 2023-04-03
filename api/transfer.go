package api

import (
	"errors"
	db "goApp/db/sqlc"
	"goApp/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateTransferRequest struct {
	FromAccountID int64 `json:"from_account_id" binding:"required,numeric,gt=0"`
	ToAccountID   int64 `json:"to_account_id" binding:"required,numeric,gt=0"`
	Amount        int64 `json:"amount" binding:"required,gt=0"`
}

func (sv *Server) CreateTransfer(ctx *gin.Context) {
	var tr CreateTransferRequest

	err := ctx.ShouldBind(&tr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	err = sv.validateTransferData(ctx, tr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	tp := db.TransferParams{
		FromAccountID: tr.FromAccountID,
		ToAccountID:   tr.ToAccountID,
		Amount:        tr.Amount,
	}

	r, err := sv.store.TransferTx(ctx, tp)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, r)

}

func (sv *Server) validateTransferData(ctx *gin.Context, arg CreateTransferRequest) error {
	if arg.FromAccountID == arg.ToAccountID {
		return errors.New("Duplicate from_account_id and to_account_id")
	}

	a1, err := sv.store.GetAccount(ctx, arg.FromAccountID)
	if err != nil {
		return err
	}

	a2, err := sv.store.GetAccount(ctx, arg.ToAccountID)
	if err != nil {
		return err
	}

	if a1.Currency != a2.Currency {
		return errors.New("Different currncy between two account")
	}

	if a1.Balance < arg.Amount {
		return errors.New("Not enough money to do transaction")
	}

	return nil

}
