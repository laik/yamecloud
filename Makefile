code-gen:
	@bash ./hack/code-generator/generate-groups.sh "deepcopy" \
      github.com/yametech/yamecloud/pkg/client github.com/yametech/yamecloud/pkg/apis \
      yamecloud:v1