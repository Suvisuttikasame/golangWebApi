package api

import (
	"goApp/authentication"
	db "goApp/db/sqlc"
	"goApp/middleware"
	"goApp/util"

	"github.com/gin-gonic/gin"
)

type Server struct {
	store      *db.Store
	router     *gin.Engine
	config     *util.Config
	tokerMaker authentication.AuthenPaseto
}

func NewServer(s *db.Store, cfg *util.Config, mkr authentication.AuthenPaseto) *Server {
	server := &Server{
		store:      s,
		config:     cfg,
		tokerMaker: mkr,
	}
	r := gin.Default()

	r.POST("/users/register", server.Register)
	r.POST("/users/login", server.Login)

	autherize := r.Group("/")
	autherize.Use(middleware.Middleware(mkr))
	autherize.POST("/accounts", server.CreateAccount)
	autherize.GET("/accounts", server.GetAllAccount)
	autherize.GET("/accounts/:id", server.GetOneAccount)

	autherize.POST("/transfers", server.CreateTransfer)

	server.router = r
	return server

}

func (sv *Server) Start(addr string) error {
	return sv.router.Run(addr)
}
