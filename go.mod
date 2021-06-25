module github.com/yametech/yamecloud

go 1.16

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/disintegration/imaging v1.6.2
	github.com/docker/distribution v2.7.1+incompatible
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/gin-gonic/gin v1.6.3
	github.com/go-resty/resty/v2 v2.6.0
	github.com/igm/sockjs-go v3.0.0+incompatible
	github.com/jinzhu/copier v0.3.2
	github.com/kubeapps/kubeapps v1.11.3
	github.com/micro/go-micro/v2 v2.9.1
	github.com/micro/micro/v2 v2.9.3
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/pflag v1.0.5
	github.com/srwiley/oksvg v0.0.0-20210519022825-9fc0c575d5fe
	github.com/srwiley/rasterx v0.0.0-20210519020934-456a8d69b780
	gopkg.in/yaml.v2 v2.4.0
	helm.sh/helm/v3 v3.6.1
	k8s.io/api v0.21.2
	k8s.io/apiextensions-apiserver v0.21.1
	k8s.io/apimachinery v0.21.2
	k8s.io/cli-runtime v0.21.2
	k8s.io/client-go v0.21.2
	k8s.io/helm v2.17.0+incompatible
	sigs.k8s.io/yaml v1.2.0
)

replace (
	github.com/coreos/etcd => github.com/yametech/etcd v3.3.20-grpc1.27-origmodule+incompatible
	github.com/micro/go-micro/v2 => github.com/yametech/go-micro-2.9.1/v2 v2.0.0-20210515101553-bf309f80bb28
	github.com/micro/micro/v2 => github.com/yametech/micro-2.9.3/v2 v2.0.0-20210515104251-6a140d397ddb
	google.golang.org/grpc => google.golang.org/grpc v1.27.0
	k8s.io/api => k8s.io/api v0.21.1
	k8s.io/apimachinery => k8s.io/apimachinery v0.21.1
	k8s.io/client-go => k8s.io/client-go v0.21.1
)
