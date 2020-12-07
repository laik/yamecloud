package configure

import (
	"flag"
	"fmt"
	"github.com/yametech/yamecloud/common"
	"github.com/yametech/yamecloud/pkg/k8s"
	"github.com/yametech/yamecloud/pkg/k8s/client"
	dynclient "k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

const (
	// InCluster when deploying in k8s, use this option
	InCluster = "InCluster"
	// Default when deploying in non k8s, use this option and the is default option
	Default = "Default"
)

var RuntimeMode = Default

var inCluster bool = false

func init() {
	flag.BoolVar(&inCluster, "in_cluster", false, "-in_cluster true")
}

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

	switch RuntimeMode {
	case Default:
		dynInterface, resetConfig, err = client.BuildClientSet(*common.KubeConfig)
	case InCluster:
		_, resetConfig, err = client.CreateInClusterConfig()
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("not define the runtime mode")
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
