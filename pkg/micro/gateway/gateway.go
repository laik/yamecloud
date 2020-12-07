package gateway

import (
	"github.com/micro/micro/plugin"
	"github.com/yametech/yamecloud/pkg/k8s"
	self "github.com/yametech/yamecloud/pkg/micro"
	"net/http"
)

var _ self.Interface = &Gateway{}

type Gateway struct {
	k8s.Interface
}

func (s *Gateway) Name() string {
	return "gateway"
}

func NewGateway(datasource k8s.Interface) self.Interface {
	return &Gateway{
		datasource,
	}
}

func (s *Gateway) Run(string) error {
	return nil
}

func (s *Gateway) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) self.Interface {
	panic("don't not implement me")
}

func (s *Gateway) DataSource() k8s.Interface {
	return s.Interface
}

func (s *Gateway) Handle(pattern string, handler http.Handler) self.Interface {
	panic("don't not implement me")
}

func NewMicroGateway(plugins ...plugin.Plugin) error {
	for _, _plugin := range plugins {
		if err := plugin.Register(_plugin); err != nil {
			return err
		}
	}
	return nil
}
