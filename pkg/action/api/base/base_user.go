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

// Get BaseUser
func (s *baseServer) GetBaseUser(g *gin.Context) {
	name := g.Param("name")
	if name == "" {
		common.RequestParametersError(g, fmt.Errorf("params not obtain name=%s", name))
		return
	}
	item, err := s.BaseUser.Get("", name)
	if err != nil {
		common.InternalServerError(g, err, err)
		return
	}
	g.JSON(http.StatusOK, item)
}

// Subscribe BaseUser
func (s *baseServer) ListBaseUser(g *gin.Context) {
	selector, err := s.generateSelector(g)
	if err != nil {
		common.InternalServerError(g, "", err)
		return
	}
	list, err := s.BaseUser.List("", selector)
	if err != nil {
		common.InternalServerError(g, "", err)
		return
	}
	g.JSON(http.StatusOK, list)
}

// ApplyBaseUser Update or Create BaseUser
func (s *baseServer) ApplyBaseUser(g *gin.Context) {
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

	password := utils.Get(_unstructured.Object, "spec.password")
	utils.Set(_unstructured.Object, "spec.password", utils.Sha1(password.(string)))

	name := _unstructured.GetName()
	newUnstructuredExtend, isUpdate, err := s.BaseUser.Apply("", name, &service.UnstructuredExtend{Unstructured: _unstructured})
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

// Update  BaseUser
func (s *baseServer) UpdateBaseUser(g *gin.Context) {
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
	password := utils.Get(_unstructured.Object, "spec.password")

	utils.Set(_unstructured.Object, "spec.password", utils.Sha1(password.(string)))

	newUnstructuredExtend, _, err := s.BaseUser.Apply("", _unstructured.GetName(), &service.UnstructuredExtend{Unstructured: _unstructured})
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

// Delete BaseUser
func (s *baseServer) DeleteBaseUser(g *gin.Context) {
	name := g.Param("name")
	if name == "" {
		common.RequestParametersError(g, fmt.Errorf("params not obtain name=%s", "", name))
		return
	}
	err := s.BaseUser.Delete("", name)
	if err != nil {
		common.InternalServerError(g, err, err)
		return
	}
	g.JSON(http.StatusOK, nil)
}

func (s *baseServer) ValidateBaseUserData(g *gin.Context, data *unstructured.Unstructured) error {
	specTenantId := utils.Get(data.Object, "spec.tenant_id")
	if specTenantId == nil || specTenantId == "" {
		return fmt.Errorf("spec.tenant_id is null")
	}

	tenantId := specTenantId.(string)
	err := s.CheckTenantId(g, tenantId)
	if err != nil {
		return fmt.Errorf("check tenantId error (%s)", err)
	}

	specDepartmentId := utils.Get(data.Object, "spec.department_id")

	selector := fmt.Sprintf("tenant.yamecloud.io=%s", tenantId)
	if specDepartmentId != nil && specDepartmentId != "" {
		departmentId := specDepartmentId.(string)
		departmentObjList, err := s.BaseDepartment.List("", selector)
		if err != nil {
			return fmt.Errorf("list department error (%s)", err)
		}

		flag := false
		for _, item := range departmentObjList.Items {
			if item.GetName() == departmentId {
				flag = true
			}
		}

		if !flag {
			return fmt.Errorf("illegal department (%s)", departmentId)
		}

		selector = fmt.Sprintf("department.yamecloud.io=%s", departmentId)

	}

	specRoles := utils.Get(data.Object, "spec.roles")
	roles := utils.ToStringArray(specRoles)

	if len(roles) > 0 {
		roleObjList, err := s.BaseRole.List("", selector)
		if err != nil {
			return fmt.Errorf("list role error (%s)", err)
		}

		roleMap := make(map[string]int, 0)
		for _, item := range roleObjList.Items {
			roleMap[item.GetName()]++
		}

		for _, role := range roles {
			if _, exist := roleMap[role]; !exist {
				return fmt.Errorf("illegal role (%s)", role)
			}
		}
	}
	return nil
}
