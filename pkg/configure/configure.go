package configure

import (
	"fmt"
	"github.com/yametech/yamecloud/common"
	"github.com/yametech/yamecloud/pkg/k8s/client"
	"github.com/yametech/yamecloud/pkg/k8s/types"
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

// InstallConfigure ...
type InstallConfigure struct {
	// kubernetes reset config
	*rest.Config
	// k8s CacheInformerFactory
	*client.CacheInformerFactory
	// k8s client
	dynclient.Interface
	// ResourceLister resource lister
	types.ResourceLister
}

func NewInstallConfigure(resLister types.ResourceLister) (*InstallConfigure, error) {
	var (
		cli         dynclient.Interface
		resetConfig *rest.Config
		err         error
	)

	switch RuntimeMode {
	case Default:
		cli, resetConfig, err = client.BuildClientSet(*common.KubeConfig)
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

	return &InstallConfigure{
		CacheInformerFactory: cacheInformerFactory,
		Interface:            cli,
		Config:               resetConfig,
		ResourceLister:       resLister,
	}, nil
}
