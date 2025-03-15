package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/minhdang2803/simple_bank/db/sqlc"
)

// HTTP request
type Server struct{
	store *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server{
	server := &Server{store: store}
	router := gin.Default()
	router.POST("/accounts", server.CreateAccount)
	router.GET("/accounts/:id", server.GetAccount)
	router.GET("/accounts", server.ListAccount)
	server.router = router
	return server
	
}

func (server *Server) Start(address string) error{
	return server.router.Run(address)
}

func errorResponse (err error) gin.H{
	return gin.H{"error":err.Error()}
}