module github.com/yametech/yamecloud

go 1.16

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32 // indirect
	github.com/gin-gonic/gin v1.6.3
	github.com/go-resty/resty/v2 v2.6.0
	github.com/igm/sockjs-go v3.0.0+incompatible
	github.com/micro/go-micro/v2 v2.9.1
	github.com/micro/micro/v2 v2.9.3
	github.com/spf13/pflag v1.0.5
	golang.org/x/mod v0.3.1-0.20200828183125-ce943fd02449 // indirect
	golang.org/x/tools v0.1.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/api v0.21.1
	k8s.io/apimachinery v0.21.1
	k8s.io/client-go v0.21.1
	k8s.io/klog/v2 v2.8.0 // indirect
	k8s.io/kube-openapi v0.0.0-20210305001622-591a79e4bda7 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.1.0 // indirect
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
