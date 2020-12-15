package service

import (
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/yametech/yamecloud/pkg/k8s"
	self "github.com/yametech/yamecloud/pkg/micro"
	"net/http"
	"time"
)

var _ self.Interface = &Service{}

type Service struct {
	k8s.Interface
	micro.Service
}

func NewService(service micro.Service, datasource k8s.Interface) self.Interface {
	return &Service{
		datasource,
		service,
	}
}

func (s *Service) Run() error {
	panic("implement me")
}

func (s *Service) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) self.Interface {
	panic("implement me")
}

func (s *Service) DataSource() k8s.Interface {
	return s.Interface
}

func (s *Service) Handle(pattern string, handler http.Handler) self.Interface {
	panic("implement me")
}

func NewMicroService(name, version string, options ...micro.Option) micro.Service {
	_service := micro.NewService(
		micro.Name(name),
		micro.Version(version),
		micro.RegisterTTL(time.Second*15),
		micro.RegisterInterval(time.Second*10),
		micro.Action(func(ctx *cli.Context) {}),
	)
	_service.Init(options...)
	return _service
}
