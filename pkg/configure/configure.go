package configure

import (
	"github.com/yametech/yamecloud/common"
	"github.com/yametech/yamecloud/pkg/k8s"
	"github.com/yametech/yamecloud/pkg/k8s/client"
	"k8s.io/client-go/discovery"
	dynclient "k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// InstallConfigure
type InstallConfigure struct {
	// kubernetes reset config
	*rest.Config
	// k8s CacheInformerFactory
	*client.CacheInformerFactory
	// k8s dyc client
	dynclient.Interface
	// Clientset
	*kubernetes.Clientset
	// ResourceLister resource lister
	k8s.ResourceLister
	// ..
	*discovery.DiscoveryClient
}

func NewInstallConfigure(resLister k8s.ResourceLister) (*InstallConfigure, error) {
	var (
		dynInterface dynclient.Interface
		resetConfig  *rest.Config
		clientSet    *kubernetes.Clientset
		err          error
	)

	if common.InCluster {
		clientSet, resetConfig, err = client.CreateInClusterConfig()
		if err != nil {
			return nil, err
		}
	}

	if resetConfig == nil {
		clientSet, dynInterface, resetConfig, err = client.BuildClientSet(*common.KubeConfig)
	}

	cacheInformerFactory, err := client.NewCacheInformerFactory(resLister, resetConfig, clientSet)
	if err != nil {
		return nil, err
	}

	if dynInterface == nil {
		dynInterface = cacheInformerFactory.Interface
	}

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(resetConfig)
	if err != nil {
		return nil, err
	}

	installConfigure := &InstallConfigure{
		Interface:            dynInterface,
		Config:               resetConfig,
		Clientset:            clientSet,
		ResourceLister:       resLister,
		CacheInformerFactory: cacheInformerFactory,
		DiscoveryClient:      discoveryClient,
	}

	return installConfigure, nil
}
