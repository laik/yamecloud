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

// Get BaseRole
func (s *baseServer) GetBaseRole(g *gin.Context) {
	name := g.Param("name")
	if name == "" {
		common.RequestParametersError(g, fmt.Errorf("params not obtain name=%s", name))
		return
	}
	item, err := s.BaseRole.Get("", name)
	if err != nil {
		common.InternalServerError(g, err, err)
		return
	}
	g.JSON(http.StatusOK, item)
}

// Subscribe BaseRole
func (s *baseServer) ListBaseRole(g *gin.Context) {
	selector, err := s.generateSelector(g)
	if err != nil {
		common.InternalServerError(g, "", err)
		return
	}
	list, err := s.BaseRole.List("", selector)
	if err != nil {
		common.InternalServerError(g, "", err)
		return
	}
	g.JSON(http.StatusOK, list)
}

// Update or Create BaseRole
func (s *baseServer) ApplyBaseRole(g *gin.Context) {
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
	err = s.ValidateRoleData(g, _unstructured)
	if err != nil {
		common.RequestParametersError(g, err)
		return
	}

	name := _unstructured.GetName()
	newUnstructuredExtend, isUpdate, err := s.BaseRole.Apply("", name, &service.UnstructuredExtend{Unstructured: _unstructured})
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

// Update BaseRole
func (s *baseServer) UpdateBaseRole(g *gin.Context) {
	name := g.Param("name")
	if name == "" {
		common.RequestParametersError(g, fmt.Errorf("params not obtain name=%s", name))
		return
	}

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
	err = s.ValidateRoleData(g, _unstructured)
	if err != nil {
		common.RequestParametersError(g, err)
		return
	}

	newUnstructuredExtend, _, err := s.BaseRole.Apply("", name, &service.UnstructuredExtend{Unstructured: _unstructured})
	if err != nil {
		common.InternalServerError(g, newUnstructuredExtend, fmt.Errorf("apply object error (%s)", err))
		return
	}

	g.JSON(
		http.StatusOK,
		[]service.UnstructuredExtend{
			*newUnstructuredExtend,
		})
}

// DeleteBaseRole none
func (s *baseServer) DeleteBaseRole(g *gin.Context) {
	name := g.Param("name")
	if name == "" {
		common.RequestParametersError(g, fmt.Errorf("params not obtain name=%s", name))
		return
	}
	err := s.BaseRole.Delete("", name)
	if err != nil {
		common.InternalServerError(g, err, err)
		return
	}
	g.JSON(http.StatusOK, nil)
}

func (s *baseServer) ValidateRoleData(g *gin.Context, data *unstructured.Unstructured) error {
	specTenantId := utils.Get(data.Object, "spec.tenant_id")
	if specTenantId == "" {
		return fmt.Errorf("spec.tenant_id is null")
	}

	tenantId := specTenantId.(string)
	err := s.CheckTenantId(g, tenantId)
	if err != nil {
		return fmt.Errorf("check tenantId error (%s)", err)
	}

	specDepartmentId := utils.Get(data.Object, "spec.department_id")
	if specDepartmentId != nil && specDepartmentId != "" {
		departmentId := specDepartmentId.(string)
		departmentObjList, err := s.BaseDepartment.List("", fmt.Sprintf("tenant.yamecloud.io=%s", tenantId))
		if err != nil {
			return fmt.Errorf("list tenant department error (%s)", err)
		}

		flag := false
		for _, item := range departmentObjList.Items {
			if item.GetName() == departmentId {
				flag = true
			}
		}

		if !flag {
			return fmt.Errorf("illegal departmentId (%s)", departmentId)
		}

		roleSpecNamespaces := utils.Get(data.Object, "spec.namespaces")
		roleNamespaces := utils.ToStringArray(roleSpecNamespaces)

		if len(roleNamespaces) > 0 {
			departmentObj, err := s.BaseDepartment.Get("", departmentId)
			if err != nil {
				return fmt.Errorf("get department error (%s)", err)
			}

			deptSpecNamespaces := utils.Get(departmentObj.Object, "spec.namespaces")
			deptNamespaces := utils.ToStringArray(deptSpecNamespaces)
			deptNamespaceMap := make(map[string]int, 0)

			for _, deptNamespace := range deptNamespaces {
				deptNamespaceMap[deptNamespace]++
			}

			for _, roleNamespace := range roleNamespaces {
				if _, exist := deptNamespaceMap[roleNamespace]; !exist {
					return fmt.Errorf("illegal namespace (%s)", roleNamespace)
				}
			}

		}
	}
	return nil
}
