package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

type IResourceServiceMaps map[k8s.ResourceType]service.IResourceService

func (i IResourceServiceMaps) ResourceService(r k8s.ResourceType) service.IResourceService {
	item, exist := i[r]
	if !exist {
		panic(fmt.Errorf("resource type service (%s) not exist", r))
	}
	return item
}

type Extends interface {
	Name() string
}

type Server struct {
	IResourceServiceMaps
	Extends
	*gin.Engine
}

func (s *Server) SetIResourceServiceMaps(i IResourceServiceMaps) {
	s.IResourceServiceMaps = i
}

func (s *Server) SetExtends(e Extends) {
	s.Extends = e
}

func NewServer() *Server {
	return &Server{
		Engine: gin.New(),
	}
}
