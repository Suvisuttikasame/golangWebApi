package api

import (
	"database/sql"
	"errors"
	"goApp/authentication"
	db "goApp/db/sqlc"
	"goApp/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,alpha"`
	Password string `json:"password" binding:"required,gte=8,alphanum"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required,alpha"`
	Password string `json:"password" binding:"required,gte=8,alphanum"`
}

type LoginResponse struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	AccessToken string `json:"accesstoken"`
}

type UserResponse struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

func (sv *Server) Register(ctx *gin.Context) {
	var req RegisterRequest

	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}
	hp, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	up := db.CreateUserParams{
		Email:    req.Email,
		Username: req.Username,
		Password: hp,
	}

	u, err := sv.store.CreateUser(ctx, up)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	res := UserResponse{
		Username: u.Username,
		Email:    u.Email,
	}

	ctx.JSON(http.StatusOK, res)

}

func (sv *Server) Login(ctx *gin.Context) {
	var req LoginRequest

	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	u, err := sv.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	r := util.CheckPasswordHash(req.Password, u.Password)
	if !r {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(errors.New("incorrect password!")))
		return
	}
	key := sv.config.SecretKey
	p, err := authentication.NewPasetoToken([]byte(key))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}
	b := authentication.Body{
		Id:       uuid.New(),
		Username: u.Username,
		Email:    u.Email,
	}
	token, err := p.CreateToken(b)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	res := LoginResponse{
		Username:    u.Username,
		Email:       u.Email,
		AccessToken: token,
	}
	ctx.JSON(http.StatusOK, res)

}
