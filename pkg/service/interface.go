package service

import (
	"github.com/yametech/yamecloud/pkg/k8s"
	"net/http"
)

type Interface interface {
	Run() error
	Name() string
	DataSource() k8s.Interface
	Handle(pattern string, handler http.Handler)
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
}

var _ Interface = &DefaultService{}

type DefaultService struct{}

func (d DefaultService) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	panic("implement me")
}

func (d DefaultService) Run() error {
	panic("implement me")
}

func (d DefaultService) Name() string {
	panic("implement me")
}

func (d DefaultService) DataSource() k8s.Interface {
	panic("implement me")
}

func (d DefaultService) Handle(pattern string, handler http.Handler) {
	panic("implement me")
}
