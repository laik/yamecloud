package gateway

//
//import (
//	"github.com/yametech/yamecloud/pkg/action/service"
//	"github.com/yametech/yamecloud/pkg/k8s"
//	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
//	"testing"
//)
//
//func Test_auth(t *testing.T) {
//}
//
//type FakeService struct {
//	Data map[string]interface{}
//}
//
//func (f *FakeService) ForceUpdate(namespace string, resource k8s.ResourceType, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, error) {
//	panic("implement me")
//}
//
//func (f *FakeService) Delete(namespace string, resource k8s.ResourceType, name string) error {
//	panic("implement me")
//}
//
//func (f *FakeService) Apply(namespace string, resource k8s.ResourceType, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, error) {
//	panic("implement me")
//}
//
//func (f *FakeService) Install(resourceType k8s.ResourceType, i service.IResourceService) {
//	f.Data[resourceType] = i
//}
//
//func (f *FakeService) List(namespace string, resource k8s.ResourceType, selector string) (*service.UnstructuredListExtend, error) {
//	panic("implement me")
//}
//
//func (f FakeService) Get(namespace string, resource k8s.ResourceType, name string) (*service.UnstructuredExtend, error) {
//	ue := &service.UnstructuredExtend{
//		Unstructured: &unstructured.Unstructured{},
//	}
//	switch resource {
//	case k8s.BaseUser:
//		err := ue.UnmarshalJSON([]byte(user))
//		if err != nil {
//			return nil, err
//		}
//	case k8s.BaseDepartment:
//		err := ue.UnmarshalJSON([]byte(dept))
//		if err != nil {
//			return nil, err
//		}
//	}
//	return ue, nil
//}
//func NewFakeService() *FakeService {
//	return &FakeService{
//		Data: make(map[string]interface{}),
//	}
//}
//
//const (
//	user = "{\"apiVersion\":\"fuxi.nip.io/v1\",\"kind\":\"BaseUser\",\"metadata\":{\"annotations\":{\"kubectl.kubernetes.io/last-applied-configuration\":\"{\\\"apiVersion\\\":\\\"fuxi.nip.io/v1\\\",\\\"kind\\\":\\\"BaseUser\\\",\\\"metadata\\\":{\\\"annotations\\\":{},\\\"name\\\":\\\"admin\\\",\\\"namespace\\\":\\\"kube-system\\\"},\\\"spec\\\":{\\\"department_id\\\":\\\"kube-admin\\\",\\\"display\\\":\\\"admin\\\",\\\"email\\\":\\\"laik.lj@me.com\\\",\\\"name\\\":\\\"admin\\\",\\\"password\\\":\\\"admin\\\"}}\\n\"},\"creationTimestamp\":\"2020-12-27T08:04:30Z\",\"generation\":2,\"managedFields\":[{\"apiVersion\":\"fuxi.nip.io/v1\",\"fieldsType\":\"FieldsV1\",\"fieldsV1\":{\"f:metadata\":{\"f:annotations\":{\".\":{},\"f:kubectl.kubernetes.io/last-applied-configuration\":{}}},\"f:spec\":{\".\":{},\"f:department_id\":{},\"f:display\":{},\"f:email\":{},\"f:name\":{}}},\"manager\":\"kubectl-client-side-apply\",\"operation\":\"Update\",\"time\":\"2020-12-27T08:04:30Z\"},{\"apiVersion\":\"fuxi.nip.io/v1\",\"fieldsType\":\"FieldsV1\",\"fieldsV1\":{\"f:spec\":{\"f:password\":{}}},\"manager\":\"base\",\"operation\":\"Update\",\"time\":\"2020-12-28T02:12:17Z\"}],\"name\":\"admin\",\"namespace\":\"kube-system\",\"resourceVersion\":\"147841930\",\"selfLink\":\"/apis/fuxi.nip.io/v1/namespaces/kube-system/baseusers/admin\",\"uid\":\"aaa46013-0977-4de5-b2a3-a1a6da743a26\"},\"spec\":{\"department_id\":\"kube-admin\",\"display\":\"admin\",\"email\":\"laik.lj@me.com\",\"name\":\"admin\",\"password\":\"Abc12345\"}}\n"
//	dept = "{\"apiVersion\":\"fuxi.nip.io/v1\",\"kind\":\"BaseDepartment\",\"metadata\":{\"annotations\":{\"kubectl.kubernetes.io/last-applied-configuration\":\"{\\\"apiVersion\\\":\\\"fuxi.nip.io/v1\\\",\\\"kind\\\":\\\"BaseDepartment\\\",\\\"metadata\\\":{\\\"annotations\\\":{},\\\"name\\\":\\\"kube-admin\\\",\\\"namespace\\\":\\\"kube-system\\\"}}\\n\"},\"creationTimestamp\":\"2020-12-27T08:04:30Z\",\"generation\":1,\"managedFields\":[{\"apiVersion\":\"fuxi.nip.io/v1\",\"fieldsType\":\"FieldsV1\",\"fieldsV1\":{\"f:metadata\":{\"f:annotations\":{\".\":{},\"f:kubectl.kubernetes.io/last-applied-configuration\":{}}}},\"manager\":\"kubectl-client-side-apply\",\"operation\":\"Update\",\"time\":\"2020-12-27T08:04:30Z\"}],\"name\":\"kube-admin\",\"namespace\":\"kube-system\",\"resourceVersion\":\"147263280\",\"selfLink\":\"/apis/fuxi.nip.io/v1/namespaces/kube-system/basedepartments/kube-admin\",\"uid\":\"25fe4f2b-af22-4462-94aa-4725c12c62a1\"}}\n"
//)
