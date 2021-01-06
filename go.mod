module github.com/yametech/yamecloud

go 1.14

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.6.3
	github.com/micro/cli v0.2.0
	github.com/micro/go-micro v1.18.0
	github.com/micro/go-plugins/registry/etcd v0.0.0-20200119172437-4fe21aa238fd
	github.com/micro/go-plugins/registry/kubernetes v0.0.0-20200119172437-4fe21aa238fd
	github.com/micro/micro v1.16.0
	github.com/spf13/pflag v1.0.5
	k8s.io/apimachinery v0.20.1
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
)

replace (
	github.com/gogo/protobuf => github.com/gogo/protobuf v1.3.1
	github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.5.3
	github.com/micro/go-micro => github.com/micro/go-micro v1.16.0
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
	k8s.io/api => k8s.io/api v0.20.1
	k8s.io/apimachinery => k8s.io/apimachinery v0.20.1
	k8s.io/client-go => k8s.io/client-go v0.20.1
)
