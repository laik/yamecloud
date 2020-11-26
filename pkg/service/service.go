package service

type IService interface {
	Start(stop <-chan struct{}) error
	Stop() error
}
