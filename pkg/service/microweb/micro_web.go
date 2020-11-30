package microweb

import (
	"github.com/micro/cli"
	"github.com/micro/go-micro/web"
	"github.com/yametech/yamecloud/pkg/k8s"
	"github.com/yametech/yamecloud/pkg/service"
	"time"
)

var _ service.Interface = &Service{}

type Service struct {
	k8s.Interface
	web.Service
}

func NewService(service2 web.Service, datasource k8s.Interface) service.Interface {
	return &Service{
		datasource,
		service2,
	}
}

func (s *Service) DataSource() k8s.Interface {
	return s.Interface
}

func (s *Service) Name() string {
	return "micro-web-service"
}

func NewMicroWebService(name, version string, options ...web.Option) (web.Service, error) {
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
