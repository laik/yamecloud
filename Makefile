code-gen:
	@bash ./hack/code-generator/generate-groups.sh "deepcopy" \
      github.com/yametech/yamecloud/pkg/client github.com/yametech/yamecloud/pkg/apis \
      yamecloud:v1

gateway:
	go run cmd/gateway/*.go api --handler=http --address 0.0.0.0:8000