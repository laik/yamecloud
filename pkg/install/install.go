package install

import (
	"github.com/yametech/yamecloud/pkg/action/api"
	"github.com/yametech/yamecloud/pkg/service/microweb"
)

func WebServiceInstall(name, version, addr string) error {
	microWebService, err := microweb.NewMicroWebService(name, version)
	if err != nil {
		panic(err)
	}
	webService := microweb.NewService(microWebService, nil)

	httpServer := api.NewDefaultImplement()

	if err := httpServer.Run(addr); err != nil {
		panic(err)
	}

	webService.Handle("/", httpServer)
	if err := webService.Run(); err != nil {
		panic(err)
	}

	return nil
}
