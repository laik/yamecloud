package network

import (
	"fmt"
	"github.com/yametech/yamecloud/pkg/action/api"
)

type storageServer struct {
	name string
	*api.Server
	// action services
}

func (s *storageServer) Name() string { return s.name }

func NewStorageServer(serviceName string, server *api.Server) *storageServer {
	storageServer := &storageServer{
		name:   serviceName,
		Server: server,
	}
	group := storageServer.Group(fmt.Sprintf("/%s", serviceName))

	_ = group
	return storageServer
}
