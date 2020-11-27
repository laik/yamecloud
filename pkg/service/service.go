package service

type IService interface {
	Start(stop <-chan struct{}) error
	Stop() error
}

var _ IService = &DefaultService{}

type DefaultService struct {
}

func (d DefaultService) Start(stop <-chan struct{}) error {
	panic("implement me")
}

func (d DefaultService) Stop() error {
	panic("implement me")
}
