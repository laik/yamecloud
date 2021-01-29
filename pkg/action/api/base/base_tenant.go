package base

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yametech/yamecloud/pkg/action/api/common"
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/micro/gateway"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"net/http"
)

// Get BaseTenant
func (s *baseServer) GetBaseTenant(g *gin.Context) {
	name := g.Param("name")
	if name == "" {
		common.RequestParametersError(g, fmt.Errorf("params not obtain name=%s", name))
		return
	}
	item, err := s.BaseTenant.Get("", name)
	if err != nil {
		common.InternalServerError(g, err, err)
		return
	}
	g.JSON(http.StatusOK, item)
}

// Subscribe BaseTenant
func (s *baseServer) ListBaseTenant(g *gin.Context) {
	selector, err := s.generateSelector(g)
	if err != nil {
		common.InternalServerError(g, "", err)
		return
	}
	list, err := s.BaseTenant.List(g.Param("namespace"), selector)
	if err != nil {
		common.InternalServerError(g, "", err)
		return
	}
	g.JSON(http.StatusOK, list)
}

// Update or Create BaseTenant
func (s *baseServer) ApplyBaseTenant(g *gin.Context) {
	namespace := g.Param("namespace")

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
	name := _unstructured.GetName()
	err = s.CheckTenantId(g, name)
	if err != nil {
		common.RequestParametersError(g, fmt.Errorf("check tenantId error (%s)", err))
		return
	}
	newUnstructuredExtend, isUpdate, err := s.BaseTenant.Apply(namespace, name, &service.UnstructuredExtend{Unstructured: _unstructured})
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

// Delete BaseTenant
func (s *baseServer) DeleteBaseTenant(g *gin.Context) {
	namespace := g.Param("namespace")
	name := g.Param("name")
	if namespace == "" || name == "" {
		common.RequestParametersError(g, fmt.Errorf("params not obtain namespace=%s name=%s", namespace, name))
		return
	}
	err := s.BaseTenant.Delete(namespace, name)
	if err != nil {
		common.InternalServerError(g, err, err)
		return
	}
	g.JSON(http.StatusOK, nil)
}

func (s *baseServer) generateSelector(g *gin.Context) (string, error) {
	selector := ""
	userIdentification := g.Request.Header.Get(gateway.UserIdentification)
	if userIdentification == "" || userIdentification != string(gateway.Admin) {
		username := g.Request.Header.Get(gateway.AuthorizationUserName)
		userObj, err := s.BaseUser.Get("kube-system", username)
		if err != nil {
			return selector, err
		}
		specTenantId, err := userObj.Get("spec.tenant_id")
		if err != nil {
			return selector, err
		}
		if specTenantId == nil {
			return selector, err
		}
		tenantId := specTenantId.(string)
		selector = fmt.Sprintf("tenant.yamecloud.io=%s", tenantId)
		return selector, nil
	}
	return selector, nil
}

func (s *baseServer) CheckTenantId(g *gin.Context, tenantId string) error {
	userIdentification := g.Request.Header.Get(gateway.UserIdentification)
	if userIdentification == string(gateway.Admin) {
		return nil
	}

	username := g.Request.Header.Get(gateway.AuthorizationUserName)
	userObj, err := s.BaseUser.Get("kube-system", username)
	if err != nil {
		return fmt.Errorf("get user error (%s)", username)
	}
	specTenantId, err := userObj.Get("spec.tenant_id")
	if err != nil {
		return fmt.Errorf("get spec.tenant_id error (%s)", username)
	}
	if tenantId != specTenantId.(string) {
		return fmt.Errorf("illegal tenantId (%s)", tenantId)
	}
	return nil
}
