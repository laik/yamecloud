package install

import (
	"fmt"
	"github.com/micro/go-micro/v2"
	"github.com/yametech/yamecloud/pkg/action/api"
	"github.com/yametech/yamecloud/pkg/k8s"
	microservice "github.com/yametech/yamecloud/pkg/micro"
	"github.com/yametech/yamecloud/pkg/micro/gateway"
	"github.com/yametech/yamecloud/pkg/micro/service"
	"github.com/yametech/yamecloud/pkg/micro/web"
	"net/http"
)

const webNormalName = "go.micro.api.%s"

func GatewayInstall(datasource k8s.Interface, server *api.Server) (microservice.Interface, error) {
	authorization := gateway.NewAuthorization(server.Interface)
	if err := gateway.NewMicroGateway(server, authorization); err != nil {
		return nil, err
	}
	return gateway.NewGateway(datasource), nil
}

func WebServiceInstall(name, version string, datasource k8s.Interface, handler http.Handler) (microservice.Interface, error) {
	microWebService, err := web.NewMicroWebService(fmt.Sprintf(webNormalName, name), version)
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
