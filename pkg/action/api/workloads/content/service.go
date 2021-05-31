package content

import "fmt"

type ServiceModels []ServiceModel

func (s ServiceModels) Get(name string) ServiceModel {
	for _, item := range s {
		if item.GetName() == name {
			return item
		}
	}
	return nil
}

func (s ServiceModels) Index(name string) int {
	for index, item := range s {
		if item.GetName() == name {
			return index
		}
	}
	return -1
}

func (s *ServiceModels) Set(name string, c ServiceModel) {
	_index := s.Index(name)
	if _index == -1 {
		*s = append(*s, c)
	} else {
		(*s)[_index] = c
	}
}

type ServiceModel interface {
	AddServiceSpec(protocol Protocol, port int, targetPort int)
	AddServiceType(serviceType ServiceType)
	GetName() string
}

var _ ServiceModel = &ServiceModelImpl{}

type ServiceModelImpl map[string]string

func NewServiceModel() ServiceModel {
	return make(ServiceModelImpl)
}

func (s ServiceModelImpl) GetName() string {
	return s["name"]
}

func (s ServiceModelImpl) AddServiceSpec(protocol Protocol, port int, targetPort int) {
	s["protocol"] = protocol
	s["port"] = fmt.Sprintf("%d", port)
	s["target_port"] = fmt.Sprintf("%d", targetPort)
}

func (s ServiceModelImpl) AddServiceType(serviceType ServiceType) {
	if serviceType == "" {
		return
	}
	s["service_type"] = serviceType
}
