package micro

import (
	"github.com/yametech/yamecloud/pkg/k8s"
	"net/http"
)

type Interface interface {
	Run() error
	Name() string
	DataSource() k8s.Interface
	Handle(pattern string, handler http.Handler) Interface
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) Interface
}

var _ Interface = &FakeService{}

type FakeService struct{}

func (d FakeService) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) Interface {
	panic("implement me")
}

func (d FakeService) Run() error {
	panic("implement me")
}

func (d FakeService) Name() string {
	panic("implement me")
}

func (d FakeService) DataSource() k8s.Interface {
	panic("implement me")
}

func (d FakeService) Handle(pattern string, handler http.Handler) Interface {
	panic("implement me")
}
