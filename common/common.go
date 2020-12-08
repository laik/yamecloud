package common

import (
	"flag"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

const (
	HttpRequestUserHeaderKey = `x-wrapper-username`
	MicroSaltUserHeader      = `go.micro.gateway.login`
	AuthorizationHeader      = "Authorization"
)

var (
	// InCluster Flag for the application runtime
	InCluster bool
	// DefaultConfigFile is the default bootstrap configuration
	KubeConfig *string
)

func init() {
	flag.BoolVar(&InCluster, "incluster", false, "-incluster true")
	if home := homedir.HomeDir(); home != "" {
		KubeConfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		KubeConfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
}
