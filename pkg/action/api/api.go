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

var _ Interface = NewAPIServer()

type APIServer struct {
	*gin.Engine
}

func NewAPIServer() *APIServer {
	return &APIServer{
		Engine: gin.New(),
	}
}
