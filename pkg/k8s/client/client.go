package client

import (
	"github.com/yametech/yamecloud/pkg/k8s"
	"time"

	"k8s.io/client-go/dynamic"
	client "k8s.io/client-go/dynamic"
	informers "k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	// High enough QPS to fit all expected use cases.
	qps = 1e6
	// High enough Burst to fit all expected use cases.
	burst = 1e6
	// full sync cache resource time
	period = 30 * time.Second
)

var sharedCacheInformerFactory *CacheInformerFactory

func BuildClientSet(path string) (*kubernetes.Clientset, client.Interface, *rest.Config, error) {
	restConfig, err := clientcmd.BuildConfigFromFlags("", path)
	if err != nil {
		return nil, nil, nil, err
	}
	cli, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		return nil, nil, nil, err
	}
	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, nil, nil, err
	}
	return clientSet, cli, restConfig, nil
}

func buildDynamicClientFromRest(clientCfg *rest.Config) (client.Interface, error) {
	clientCfg.QPS = qps
	clientCfg.Burst = burst
	dynClient, err := client.NewForConfig(clientCfg)
	if err != nil {
		return nil, err
	}
	return dynClient, nil
}

type CacheInformerFactory struct {
	Interface client.Interface
	Informer  informers.DynamicSharedInformerFactory
	Clientset *kubernetes.Clientset
	stopChan  chan struct{}
}

func NewCacheInformerFactory(resLister k8s.ResourceLister, restConf *rest.Config, clientset *kubernetes.Clientset) (*CacheInformerFactory, error) {
	if sharedCacheInformerFactory != nil {
		return sharedCacheInformerFactory, nil
	}
	_client, err := buildDynamicClientFromRest(restConf)
	if err != nil {
		return nil, err
	}
	stop := make(chan struct{})
	sharedInformerFactory := informers.NewDynamicSharedInformerFactory(_client, period)
	resLister.Ranges(sharedInformerFactory, stop)
	sharedInformerFactory.Start(stop)

	sharedCacheInformerFactory =
		&CacheInformerFactory{
			Interface: _client,
			Informer:  sharedInformerFactory,
			Clientset: clientset,
			stopChan:  stop,
		}

	return sharedCacheInformerFactory, nil
}

func CreateInClusterConfig() (*kubernetes.Clientset, *rest.Config, error) {
	restConfig, err := rest.InClusterConfig()
	if err != nil {
		return nil, nil, err
	}
	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, nil, err
	}
	return clientSet, restConfig, nil
}
