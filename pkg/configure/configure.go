package configure

import (
	"github.com/yametech/yamecloud/common"
	"github.com/yametech/yamecloud/pkg/k8s"
	"github.com/yametech/yamecloud/pkg/k8s/client"
	dynclient "k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

// InstallConfigure
type InstallConfigure struct {
	// kubernetes reset config
	*rest.Config
	// k8s CacheInformerFactory
	*client.CacheInformerFactory
	// k8s client
	dynclient.Interface
	// ResourceLister resource lister
	k8s.ResourceLister
}

func NewInstallConfigure(resLister k8s.ResourceLister) (*InstallConfigure, error) {
	var (
		dynInterface dynclient.Interface
		resetConfig  *rest.Config
		err          error
	)

	if common.InCluster {
		_, resetConfig, err = client.CreateInClusterConfig()
		if err != nil {
			return nil, err
		}
	}

	if resetConfig == nil {
		dynInterface, resetConfig, err = client.BuildClientSet(*common.KubeConfig)
	}

	cacheInformerFactory, err := client.NewCacheInformerFactory(resLister, resetConfig)
	if err != nil {
		return nil, err
	}

	installConfigure := &InstallConfigure{
		CacheInformerFactory: cacheInformerFactory,

		ResourceLister: resLister,

		Interface: dynInterface,
		Config:    resetConfig,
	}

	return installConfigure, nil
}
