package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/minhdang2803/simple_bank/db/sqlc"
)

// HTTP request
type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// custom validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validatorCurrency)
	}

	router.POST("/accounts", server.CreateAccount)
	router.GET("/accounts/:id", server.GetAccount)
	router.GET("/accounts", server.ListAccount)
	router.POST("/transfers", server.CreateTransfer)
	router.POST("/users", server.CreateUser)
	router.GET("/users", server.GetUser)
	router.POST("/users/update", server.UpdateUser)
	server.router = router
	return server

}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
