package workloads

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yametech/yamecloud/pkg/action/api/common"
	"github.com/yametech/yamecloud/pkg/utils"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"net/http"
	"net/url"
	"strings"
)

// ListCustomResourceDefinition List CustomResourceDefinition
func (w *workloadServer) ListCustomResourceDefinition(g *gin.Context) {
	list, err := w.CRD.List("", "")
	if err != nil {
		common.RequestParametersError(g, fmt.Errorf("params not obtain or params parse error: %s", err))
		return
	}
	g.JSON(http.StatusOK, list)
}

func (w *workloadServer) ListCustomResourceRouter(gvrString []string) ([]string, error) {
	list, err := w.CRD.List("", "")
	if err != nil {
		return nil, err
	}

	customResourceDefinitionList := &v1beta1.CustomResourceDefinitionList{}
	err = utils.UnmarshalObject(list, customResourceDefinitionList)
	if err != nil {
		return nil, err
	}

	results := make([]string, 0)
	for _, item := range customResourceDefinitionList.Items {
		var apiVersionUrl string
		for _, version := range item.Spec.Versions {
			apiVersionUrl = fmt.Sprintf("%s/%s/%s", item.Spec.Group, version.Name, item.Spec.Names.Plural)
			if utils.Contains(gvrString, apiVersionUrl) {
				continue
			}
			results = append(results, apiVersionUrl)
		}
	}
	return results, nil
}

func trimSpace(slice []string) []string {
	for i, s := range slice {
		if s == "" {
			slice = slice[i+1:]
		}
	}
	return slice
}

// ListGeneralCustomResourceDefinition List General CustomResourceDefinition
func (w *workloadServer) ListGeneralCustomResourceDefinition(g *gin.Context) {
	u, err := url.Parse(g.Request.RequestURI)
	if err != nil {
		common.RequestParametersError(g, fmt.Errorf("params not obtain or params parse url error: %s", err))
		return
	}
	paths := trimSpace(strings.Split(u.Path, "/"))
	if len(paths) < 5 {
		common.RequestParametersError(g, fmt.Errorf("request url error %s", g.Request.RequestURI))
		return
	}
	group, version, resource := paths[2], paths[3], paths[4]

	//w.generic.SetGroupVersionResource(groupVersionResource)
	list, err := w.CRD.ListGVR("", schema.GroupVersionResource{Group: group, Version: version, Resource: resource}, "")
	if err != nil {
		common.InternalServerError(g, err, err)
		return
	}

	g.JSON(http.StatusOK, list)
}
