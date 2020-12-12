package install

import (
	"fmt"
	"github.com/micro/go-micro"
	"github.com/yametech/yamecloud/pkg/k8s"
	microservice "github.com/yametech/yamecloud/pkg/micro"
	"github.com/yametech/yamecloud/pkg/micro/gateway"
	"github.com/yametech/yamecloud/pkg/micro/service"
	"github.com/yametech/yamecloud/pkg/micro/web"
	"net/http"
)

const nameNormal = "go.micro.api.%s"

func GatewayInstall(datasource k8s.Interface, handler http.Handler) (microservice.Interface, error) {
	if err := gateway.NewMicroGateway(handler); err != nil {
		return nil, err
	}
	return gateway.NewGateway(datasource), nil
}

func WebServiceInstall(name, version string, datasource k8s.Interface, handler http.Handler) (microservice.Interface, error) {
	name = fmt.Sprintf(nameNormal, name)
	microWebService, err := web.NewMicroWebService(name, version)
	if err != nil {
		return nil, err
	}
	return web.NewService(microWebService, datasource).Handle("/", handler), nil
}

func ServiceInstall(name, version string, datasource k8s.Interface) (microservice.Interface, error) {
	opts := make([]micro.Option, 0)
	microService := service.NewMicroService(name, version, opts...)
	return service.NewService(microService, datasource), nil
}
