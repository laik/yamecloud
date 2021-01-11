package web

import (
	"github.com/micro/go-micro/v2/web"
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
	webService := web.NewService(
		web.Name(name),
		web.Version(version),
		web.RegisterTTL(time.Second*15),
		web.RegisterInterval(time.Second*10),
	)
	if err := webService.Init(); err != nil {
		return nil, err
	}
	return webService, nil
}
