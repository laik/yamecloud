module github.com/yametech/yamecloud

go 1.14

require (
	github.com/gin-gonic/gin v1.6.3
	github.com/micro/cli v0.2.0
	github.com/micro/go-micro v1.16.0
	github.com/micro/micro v1.16.0 // indirect
	k8s.io/apimachinery v0.19.3
	k8s.io/client-go v1.5.1
)

replace (
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
	github.com/gogo/protobuf => github.com/gogo/protobuf v1.3.1
	github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.5.3
	k8s.io/api => k8s.io/api v0.19.3
	k8s.io/apimachinery => k8s.io/apimachinery v0.19.3
	k8s.io/client-go => k8s.io/client-go v0.19.3
)
