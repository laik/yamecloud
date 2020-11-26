module github.com/yametech/yamecloud

go 1.15

require (
	k8s.io/apimachinery v0.19.0
	k8s.io/client-go 0.0.0
)

replace (
	github.com/gogo/protobuf => github.com/gogo/protobuf v1.3.1
	//google.golang.org/grpc => google.golang.org/grpc v1.26.0
	k8s.io/api => k8s.io/api v0.19.0
	k8s.io/apimachinery => k8s.io/apimachinery v0.19.0
	k8s.io/client-go => k8s.io/client-go v0.19.0
)
