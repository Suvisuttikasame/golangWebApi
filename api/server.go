package api

import (
	db "goApp/db/sqlc"
	"goApp/util"

	"github.com/gin-gonic/gin"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
	config *util.Config
}

func NewServer(s *db.Store, cfg *util.Config) *Server {
	server := &Server{
		store:  s,
		config: cfg,
	}
	r := gin.Default()

	r.POST("/users/register", server.Register)
	r.POST("/users/login", server.Login)

	r.POST("/accounts", server.CreateAccount)
	r.GET("/accounts", server.GetAllAccount)
	r.GET("/accounts/:id", server.GetOneAccount)

	r.POST("/transfers", server.CreateTransfer)

	server.router = r
	return server

}

func (sv *Server) Start(addr string) error {
	return sv.router.Run(addr)
}

func ErrorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
