module github.com/yametech/yamecloud

go 1.14

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32 // indirect
	github.com/gin-gonic/gin v1.6.3
	github.com/micro/go-micro/v2 v2.9.1
	github.com/micro/micro/v2 v2.9.3
	github.com/spf13/pflag v1.0.5
	k8s.io/apimachinery v0.20.1
	k8s.io/client-go v0.0.0-00010101000000-000000000000
)

replace (
	github.com/gogo/protobuf => github.com/gogo/protobuf v1.3.1
	github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.5.3
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
	k8s.io/api => k8s.io/api v0.19.3
	k8s.io/apimachinery => k8s.io/apimachinery v0.19.3
	k8s.io/client-go => k8s.io/client-go v0.19.3
)
