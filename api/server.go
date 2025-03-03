package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/gurukanth/simplebank/db/sqlc"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) (server *Server) {
	server = &Server{store: store}
	router := gin.Default()

	//Handle router
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts/", server.listAccounts)

	router.POST("/transfers", server.createTransfer)

	server.router = router
	return
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
