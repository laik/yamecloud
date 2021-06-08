# the makefile describe
REPO=yametech
VERSION=0.2.0

CRD_OPTIONS ?= "crd:trivialVersions=true"

# find or download controller-gen
# download controller-gen if necessary
controller-gen:
ifeq (, $(shell which controller-gen))
	@{ \
	CONTROLLER_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$CONTROLLER_GEN_TMP_DIR ;\
	curl -OL https://github.com/yametech/controller-tools/archive/v0.4.1.tar.gz && tar -zxvf v0.4.1.tar.gz && cd controller-tools-0.4.1 ;\
	cd ./cmd/controller-gen && go install && cd ../helpgen && go install && cd ../type-scaffold && go install ;\
	rm -rf $$CONTROLLER_GEN_TMP_DIR ;\
	}
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif

# Just install controller-gen tools set
install-tools: controller-gen
	@echo "install controller-gen done"

# Generate code
generate: controller-gen
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

# Generate manifests e.g. CRD, RBAC etc.
manifests: controller-gen
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=deploy/crds > /dev/null 2>&1 &
	  
dep:
	go mod vendor

build: dep
	go build ./cmd/...

build-image: servicemesh tekton service base gateway shell watcher sdn workloadplus workloads push-image

push-image: 
	docker push ${REPO}/servicemesh:${VERSION}
	docker push ${REPO}/tekton:${VERSION}
	docker push ${REPO}/service:${VERSION}
	docker push ${REPO}/base:${VERSION}
	docker push ${REPO}/gateway:${VERSION}
	docker push ${REPO}/shell:${VERSION}
	docker push ${REPO}/watcher:${VERSION}
	docker push ${REPO}/sdn:${VERSION}
	docker push ${REPO}/workloads:${VERSION}
	docker push ${REPO}/workloadplus:${VERSION}

sdn:
	docker build -t ${REPO}/sdn:${VERSION} -f images/Dockerfile.sdn .

watcher:
	ame${REPO}/watcher:${VERSION} -f images/Dockerfile.watcher .


servicemesh:
	docker build -t ${REPO}/servicemesh:${VERSION} -f images/Dockerfile.servicemesh .

tekton:
	docker build -t ${REPO}/tekton:${VERSION} -f images/Dockerfile.tekton .

service:
	docker build -t ${REPO}/service:${VERSION} -f images/Dockerfile.service .

base:
	docker build -t ${REPO}/base:${VERSION} -f images/Dockerfile.base .

gateway:
	docker build -t ${REPO}/gateway:${VERSION} -f images/Dockerfile.gateway .

shell:
	docker build -t ${REPO}/shell:${VERSION} -f images/Dockerfile.shell .

workloads:
	docker build -t ${REPO}/workloads:${VERSION} -f images/Dockerfile.workloads .

workloadplus:
	docker build -t ${REPO}/workloadplus:${VERSION} -f images/Dockerfile.workloadplus .

gateway_run:
	go run cmd/gateway/*.go api --handler=http --address 0.0.0.0:8080

base_run:
	go run cmd/base/*.go

service_run:
	go run cmd/service/*.go

servicemesh_run:
	go run cmd/servicemesh/*.go

tekton_run:
	go run cmd/tekton/*.go

workloads_run:
	go run cmd/workloads/*.go

workloadplus_run:
	go run cmd/workloadplus/*.go

watcher_run:
	SUBLIST="" SUBTOPIC=tekton,ovn go run cmd/watcher/*.go

shell_run:
	go run cmd/shell/*.go