package gateway

import (
	"github.com/micro/micro/plugin"
	"github.com/yametech/yamecloud/pkg/k8s"
	self "github.com/yametech/yamecloud/pkg/micro"
	"net/http"
)

func Wrapper(handler http.Handler) plugin.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(handler.ServeHTTP)
	}
}

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

func (s *Gateway) Run() error { return nil }

func (s *Gateway) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) self.Interface {
	panic("don't not implement me")
}

func (s *Gateway) DataSource() k8s.Interface { return s.Interface }

func (s *Gateway) Handle(pattern string, handler http.Handler) self.Interface {
	panic("don't not implement me")
}

func NewMicroGateway(handlers ...http.Handler) error {
	handlerWrappers := make([]plugin.Handler, 0)
	for _, handler := range handlers {
		handlerWrappers = append(handlerWrappers, Wrapper(handler))
	}
	if err := plugin.Register(plugin.NewPlugin(plugin.WithHandler(handlerWrappers...))); err != nil {
		return err
	}
	return nil
}
