package api

import (
	db "goApp/db/sqlc"

	"github.com/gin-gonic/gin"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(s *db.Store) *Server {
	server := &Server{store: s}
	r := gin.Default()

	r.POST("/accounts", server.CreateAccount)

	server.router = r
	return server

}

func (sv *Server) Start(addr string) error {
	return sv.router.Run(addr)
}

func ErrorResponse(err error) gin.H {
	return gin.H{"error": err}
}
