package events

import (
	"fmt"
	"github.com/yametech/yamecloud/pkg/action/api"
)

type eventsServer struct {
	name string
	*api.Server
	// action types
}

func (s *eventsServer) Name() string { return s.name }

func NewEventsServer(serviceName string, server *api.Server) *eventsServer {
	eventsServer := &eventsServer{
		name:   serviceName,
		Server: server,
	}
	group := eventsServer.Group(fmt.Sprintf("/%s", serviceName))

	_ = group
	return eventsServer
}
