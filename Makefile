# the makefile describe
REPO=harbor.ym/devops
VERSION=0.0.1

code-gen:
	@bash ./hack/code-generator/generate-groups.sh "deepcopy" \
      github.com/yametech/yamecloud/pkg/client github.com/yametech/yamecloud/pkg/apis \
      yamecloud:v1

dep:
	go mod vendor

build: dep
	go build ./cmd/...

build-image: servicemesh tekton service

servicemesh:
	docker build -t ${REPO}/service-mesh:${VERSION} -f images/Dockerfile.servicemesh .

tekton:
	docker build -t ${REPO}/tekton:${VERSION} -f images/Dockerfile.tekton .

service:
	docker build -t ${REPO}/service:${VERSION} -f images/Dockerfile.service .

gateway:
	go run cmd/gateway/*.go api --handler=http --address 0.0.0.0:8000
