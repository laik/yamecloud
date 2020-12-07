package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Interface interface {
	http.Handler
	Run(addr ...string) error
}

type HandleFunc = gin.HandlerFunc

var _ Interface = NewServer()

type Server struct {
	*gin.Engine
}

func NewServer() *Server {
	return &Server{
		Engine: gin.New(),
	}
}
