package gateway

import (
	"github.com/micro/go-micro/v2/config/cmd"
	"github.com/micro/micro/v2/plugin"
	"github.com/yametech/yamecloud/pkg/k8s"
	self "github.com/yametech/yamecloud/pkg/micro"
	"net/http"
)

var _ self.Interface = &Gateway{}

type Gateway struct {
	k8s.Interface
}

func NewGateway(datasource k8s.Interface) self.Interface {
	return &Gateway{
		datasource,
	}
}

func (s *Gateway) Name() string { return "gateway" }

func (s *Gateway) Run() error {
	cmd.Init()
	return nil
}

func (s *Gateway) DataSource() k8s.Interface { return s.Interface }

func (s *Gateway) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) self.Interface {
	panic("don't not implement me")
}
func (s *Gateway) Handle(pattern string, handler http.Handler) self.Interface {
	panic("don't not implement me")
}

func NewMicroGateway(handler http.Handler) error {
	handlerWrappers := []plugin.Handler{
		filter(handler),
	}
	if err := plugin.Register(plugin.NewPlugin(plugin.WithHandler(handlerWrappers...))); err != nil {
		return err
	}
	return nil
}
