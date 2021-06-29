package base

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yametech/yamecloud/pkg/action/api/common"
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/utils"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"net/http"
)

// Get BaseDepartment
func (s *baseServer) GetBaseDepartment(g *gin.Context) {
	name := g.Param("name")
	if name == "" {
		common.RequestParametersError(g, fmt.Errorf("params not obtain name=%s", name))
		return
	}
	item, err := s.BaseDepartment.Get("", name)
	if err != nil {
		common.InternalServerError(g, err, err)
		return
	}
	g.JSON(http.StatusOK, item)
}

// Subscribe BaseDepartment
func (s *baseServer) ListBaseDepartment(g *gin.Context) {
	selector, err := s.generateSelector(g)
	if err != nil {
		common.InternalServerError(g, "", err)
		return
	}
	list, err := s.BaseDepartment.List("", selector)
	if err != nil {
		common.InternalServerError(g, "", err)
		return
	}
	g.JSON(http.StatusOK, list)
}

// ApplyBaseDepartment Update or Create BaseDepartment
func (s *baseServer) ApplyBaseDepartment(g *gin.Context) {
	raw, err := g.GetRawData()
	if err != nil {
		common.RequestParametersError(g, fmt.Errorf("get raw data error (%s)", err))
		return
	}

	_unstructured := &unstructured.Unstructured{}
	if err := _unstructured.UnmarshalJSON(raw); err != nil {
		common.RequestParametersError(g, fmt.Errorf("unmarshal from json data error (%s)", err))
		return
	}
	err = s.ValidateBaseUserData(g, _unstructured)
	if err != nil {
		common.RequestParametersError(g, err)
		return
	}

	name := _unstructured.GetName()
	newUnstructuredExtend, isUpdate, err := s.BaseDepartment.Apply("", name, &service.UnstructuredExtend{Unstructured: _unstructured})
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

// Delete BaseDepartment
func (s *baseServer) DeleteBaseDepartment(g *gin.Context) {
	name := g.Param("name")
	if name == "" {
		common.RequestParametersError(g, fmt.Errorf("params not obtain name=%s", name))
		return
	}
	err := s.BaseDepartment.Delete("", name)
	if err != nil {
		common.InternalServerError(g, err, err)
		return
	}
	g.JSON(http.StatusOK, nil)
}

func (s *baseServer) ValidateDepartmentData(g *gin.Context, data *unstructured.Unstructured) error {
	specTenantId := utils.Get(data.Object, "spec.tenant_id")
	if specTenantId == "" {
		return fmt.Errorf("spec.tenant_id is null")
	}
	tenantId := specTenantId.(string)
	err := s.CheckTenantId(g, tenantId)
	if err != nil {
		return fmt.Errorf("check tenantId error (%s)", err)
	}
	deptSpecNamespaces := utils.Get(data.Object, "spec.namespaces")
	deptNamespaces := utils.ToStringArray(deptSpecNamespaces)
	if len(deptNamespaces) > 0 {

		tenantObj, err := s.BaseTenant.Get("", tenantId)
		if err != nil {
			return fmt.Errorf(tenantId + " do not exist")
		}
		specNamespaces, err := tenantObj.Get("spec.namespaces")
		if err != nil {
			return fmt.Errorf("get tenant namespaces error (%s)", err)
		}
		tenantNamespaces := utils.ToStringArray(specNamespaces)
		namespaceMap := make(map[string]int, 0)
		for _, v := range tenantNamespaces {
			namespaceMap[v]++
		}
		for _, deptNamespace := range deptNamespaces {
			_, exits := namespaceMap[deptNamespace]
			if !exits {
				return fmt.Errorf("illegal namespace (%s)", deptNamespace)
			}
		}
	}
	specOwner := utils.Get(data.Object, "spec.owner")
	if specOwner != "" {
		owner := specOwner.(string)
		selector := fmt.Sprintf("tenant.yamecloud.io=%s", tenantId)
		userObjList, err := s.BaseUser.List("", selector)
		if err != nil {
			return fmt.Errorf("get tenant user error (%s)", err)
		}
		exist := false
		for _, userObj := range userObjList.Items {
			if userObj.GetName() == owner {
				exist = true
				break
			}

		}
		if !exist {
			return fmt.Errorf("illegal owner (%s)", owner)
		}
	}
	return nil
}
