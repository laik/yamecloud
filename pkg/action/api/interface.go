package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yametech/yamecloud/pkg/action/service"
)

type Extends interface {
	Name() string
}

type Server struct {
	Extends
	*gin.Engine
	service.Interface
}

func (s *Server) SetExtends(e Extends) { s.Extends = e }

func NewServer(p service.Interface) *Server {
	engine := gin.New()
	engine.Use([]gin.HandlerFunc{gin.Logger(), gin.Recovery()}...)

	return &Server{
		Engine:    engine,
		Interface: p,
	}
}
