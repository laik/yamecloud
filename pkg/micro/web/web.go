package web

import (
	"github.com/micro/cli"
	"github.com/micro/go-micro/web"
	"github.com/micro/go-plugins/registry/kubernetes"
	"github.com/yametech/yamecloud/common"
	"github.com/yametech/yamecloud/pkg/k8s"
	self "github.com/yametech/yamecloud/pkg/micro"
	"net/http"
	"time"
)

var _ self.Interface = &Service{}

type Service struct {
	k8s.Interface
	web.Service
}

func NewService(service2 web.Service, datasource k8s.Interface) self.Interface {
	return &Service{
		datasource,
		service2,
	}
}

func (s *Service) Run() error { return s.Service.Run() }

func (s *Service) Handle(pattern string, handler http.Handler) self.Interface {
	s.Service.Handle(pattern, handler)
	return s
}

func (s *Service) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) self.Interface {
	s.Service.HandleFunc(pattern, handler)
	return s
}

func (s *Service) DataSource() k8s.Interface { return s.Interface }

func (s *Service) Name() string { return "micro-web-micro" }

func NewMicroWebService(name, version string) (web.Service, error) {
	options := make([]web.Option, 0)
	if common.InCluster {
		options = append(options, web.Registry(kubernetes.NewRegistry()))
	}
	_service := web.NewService(
		web.Name(name),
		web.Version(version),
		web.RegisterTTL(time.Second*15),
		web.RegisterInterval(time.Second*10),
		web.Action(func(ctx *cli.Context) {}),
	)
	if err := _service.Init(options...); err != nil {
		return nil, err
	}
	return _service, nil
}
