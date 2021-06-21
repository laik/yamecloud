package accesscontrol

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yametech/yamecloud/pkg/action/api/common"
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/utils"
	v1 "k8s.io/api/core/v1"
	"net/http"
)

type patchServiceAccountSecret struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

// List ServiceAccount
func (s *accessControlServer) ListServiceAccount(g *gin.Context) {
	list, err := s.ServiceAccount.List("", "")
	if err != nil {
		common.InternalServerError(g, "", err)
		return
	}
	g.JSON(http.StatusOK, list)
}

func (s *accessControlServer) GetServiceAccount(g *gin.Context) {
	namespace := g.Param("namespace")
	name := g.Param("name")
	if name == "" || namespace == "" {
		common.RequestParametersError(g, fmt.Errorf("params not obtain namespace=%s name=%s", namespace, name))
		return
	}

	item, err := s.ServiceAccount.Get(namespace, name)
	if err != nil {
		common.InternalServerError(g, "", err)
		return
	}
	g.JSON(http.StatusOK, item)
}

func (s *accessControlServer) ApplyServiceAccount(g *gin.Context) {
	method := g.Param("method")
	rawData, err := g.GetRawData()
	if err != nil {
		common.RequestParametersError(g, fmt.Errorf("get raw data error (%s)", err))
		return
	}

	pad := patchServiceAccountSecret{}
	err = json.Unmarshal(rawData, &pad)
	if err != nil {
		common.InternalServerError(g, "", err)
		return
	}
	serviceAccountUnstructed, err := s.ServiceAccount.Get(pad.Namespace, "default")
	if err != nil {
		common.InternalServerError(g, "", err)
		return
	}

	serviceAccount := &v1.ServiceAccount{}
	if err := utils.UnmarshalObject(serviceAccountUnstructed, serviceAccount); err != nil {
		common.InternalServerError(g, "", err)
		return
	}
	if method == "add" {
		serviceAccount.Secrets = append(serviceAccount.Secrets, v1.ObjectReference{Name: pad.Name})
	} else {
		serviceAccount.Secrets = removeObjectReference(serviceAccount.Secrets, pad.Name)
	}

	unstructured, err := utils.ObjectToUnstructured(serviceAccount)
	if err != nil {
		common.InternalServerError(g, "", err)
		return
	}
	newUnstructuredExtend, isUpdate, err := s.ServiceAccount.Apply(pad.Namespace, "default", &service.UnstructuredExtend{Unstructured: unstructured})
	if err != nil {
		common.InternalServerError(g, newUnstructuredExtend, fmt.Errorf("apply object error (%s)", err))
		return
	}
	if isUpdate {
		g.JSON(
			http.StatusOK,
			[]service.UnstructuredExtend{
				*newUnstructuredExtend,
			})
	} else {
		g.JSON(http.StatusOK, newUnstructuredExtend)
	}
}

func removeObjectReference(slice []v1.ObjectReference, name string) []v1.ObjectReference {
	tmpMap := make(map[string]v1.ObjectReference)
	for i := range slice {
		tmpMap[slice[i].Name] = slice[i]
	}
	delete(tmpMap, name)
	result := make([]v1.ObjectReference, 0)
	for _, v := range tmpMap {
		result = append(result, v)
	}
	return result
}
