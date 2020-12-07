package install

import (
	"github.com/micro/go-micro"
	"github.com/yametech/yamecloud/pkg/k8s"
	"github.com/yametech/yamecloud/pkg/micro/gateway"
	"github.com/yametech/yamecloud/pkg/micro/service"
	"github.com/yametech/yamecloud/pkg/micro/web"
	"net/http"
)

func GatewayInstall(datasource k8s.Interface) error {
	if err := gateway.NewMicroGateway(); err != nil {
		return err
	}
	return gateway.NewGateway(datasource).Run()
}

func WebInstall(name, version, addr string, datasource k8s.Interface, handler http.Handler) error {
	microWebService, err := web.NewMicroWebService(name, version)
	if err != nil {
		return err
	}
	return web.NewService(microWebService, datasource).
		Handle("/", handler).
		Run()
}

func ServiceInstall(name, version string, datasource k8s.Interface) error {
	opts := make([]micro.Option, 0)
	microService := service.NewMicroService(name, version, opts...)
	return service.NewService(microService, datasource).Run()
}
