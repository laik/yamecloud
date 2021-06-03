package datasource

import (
	"context"
	"fmt"
	"github.com/yametech/yamecloud/pkg/configure"
	"github.com/yametech/yamecloud/pkg/k8s"
	"github.com/yametech/yamecloud/pkg/utils"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/retry"
	"reflect"
	"time"
)

var _ k8s.Interface = &dataSource{}

func NewInterface(configure *configure.InstallConfigure) k8s.Interface {
	return &dataSource{configure}
}

type dataSource struct {
	*configure.InstallConfigure
}

func (d *dataSource) ApplyGVR(namespace, name string, gvr *schema.GroupVersionResource, unstructured *unstructured.Unstructured) (newUnstructured *unstructured.Unstructured, isUpdate bool, err error) {
	err = retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		ctx := context.TODO()
		getObj, getErr := d.
			Interface.
			Resource(*gvr).
			Namespace(namespace).
			Get(ctx, name, metav1.GetOptions{})

		if errors.IsNotFound(getErr) {
			newObj, createErr := d.
				Interface.
				Resource(*gvr).
				Namespace(namespace).
				Create(ctx, unstructured, metav1.CreateOptions{})
			newUnstructured = newObj
			return createErr
		}
		if getErr != nil {
			return getErr
		}

		d.compareObject(getObj, unstructured, false)

		newObj, updateErr := d.
			Interface.
			Resource(*gvr).
			Namespace(namespace).
			Update(ctx, getObj, metav1.UpdateOptions{})

		newUnstructured = newObj
		isUpdate = true
		return updateErr
	})

	return
}

func (d *dataSource) DiscoveryClient() *discovery.DiscoveryClient {
	return d.InstallConfigure.DiscoveryClient
}

func (d *dataSource) ListGVR(namespace string, gvr schema.GroupVersionResource, selector string) (*unstructured.UnstructuredList, error) {
	return d.Interface.
		Resource(gvr).
		Namespace(namespace).
		List(
			context.TODO(),
			metav1.ListOptions{LabelSelector: selector},
		)
}

func (d *dataSource) RESTClient() rest.Interface {
	return d.InstallConfigure.RESTClient()
}

func (d *dataSource) ClientSet() *kubernetes.Clientset {
	return d.InstallConfigure.Clientset
}

func (d *dataSource) XGet(namespace string, resourceType k8s.ResourceType, name string) (*unstructured.Unstructured, error) {
	gvr, err := d.ResourceLister.GroupVersionResource(resourceType)
	if err != nil {
		return nil, err
	}
	rObject, err := d.Informer.
		ForResource(gvr).
		Lister().
		ByNamespace(namespace).
		Get(name)
	if err != nil {
		return nil, err
	}
	return utils.ObjectToUnstructured(rObject)
}

func (d *dataSource) Watch(namespace string, resourceType k8s.ResourceType, resourceVersion string, selector string) (<-chan watch.Event, error) {
	gvr, err := d.ResourceLister.GroupVersionResource(resourceType)
	if err != nil {
		return nil, err
	}
	timeoutSeconds := int64(0)
	watchInterface, err := d.Interface.
		Resource(gvr).
		Namespace(namespace).
		Watch(context.TODO(), metav1.ListOptions{
			ResourceVersion: resourceVersion,
			LabelSelector:   selector,
			TimeoutSeconds:  &timeoutSeconds,
		})
	if err != nil {
		return nil, err
	}
	return watchInterface.ResultChan(), nil
}

func (d *dataSource) Apply(namespace string, resourceType k8s.ResourceType, name string, unstructured *unstructured.Unstructured, forceUpdate bool) (newUnstructured *unstructured.Unstructured, isUpdate bool, err error) {
	err = retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		gvr, err := d.ResourceLister.GroupVersionResource(resourceType)
		if err != nil {
			return err
		}

		ctx := context.TODO()
		getObj, getErr := d.
			Interface.
			Resource(gvr).
			Namespace(namespace).
			Get(ctx, name, metav1.GetOptions{})

		if errors.IsNotFound(getErr) {
			newObj, createErr := d.
				Interface.
				Resource(gvr).
				Namespace(namespace).
				Create(ctx, unstructured, metav1.CreateOptions{})
			newUnstructured = newObj
			return createErr
		}
		if getErr != nil {
			return getErr
		}

		d.compareObject(getObj, unstructured, forceUpdate)

		newObj, updateErr := d.
			Interface.
			Resource(gvr).
			Namespace(namespace).
			Update(ctx, getObj, metav1.UpdateOptions{})

		newUnstructured = newObj
		isUpdate = true
		return updateErr
	})

	return
}

func (d *dataSource) Delete(namespace string, resourceType k8s.ResourceType, name string) error {
	gvr, err := d.ResourceLister.GroupVersionResource(resourceType)
	if err != nil {
		return err
	}
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		return d.Interface.
			Resource(gvr).
			Namespace(namespace).
			Delete(context.TODO(), name, metav1.DeleteOptions{})
	})
}

func (d *dataSource) Patch(namespace string, resourceType k8s.ResourceType, name string, data []byte) (result *unstructured.Unstructured, err error) {
	err = retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		gvr, err := d.ResourceLister.GroupVersionResource(resourceType)
		if err != nil {
			return err
		}
		result, err = d.
			Interface.
			Resource(gvr).
			Namespace(namespace).
			Patch(context.TODO(), name, types.MergePatchType, data, metav1.PatchOptions{})
		return err
	})
	return
}

func (d *dataSource) ListLimit(namespace string, resourceType k8s.ResourceType, flag string, pos, size int64, selector string) (*unstructured.UnstructuredList, error) {
	gvr, err := d.ResourceLister.GroupVersionResource(resourceType)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{LabelSelector: selector}
	if flag != "" {
		//opts.Continue = flag
	}
	if size > 0 {
		opts.Limit = size + pos
	}
	return d.Interface.
		Resource(gvr).
		Namespace(namespace).
		List(context.TODO(), opts)
}

func (d *dataSource) List(namespace string, resourceType k8s.ResourceType, selector string) (*unstructured.UnstructuredList, error) {
	gvr, err := d.ResourceLister.GroupVersionResource(resourceType)
	if err != nil {
		return nil, err
	}
	return d.Interface.
		Resource(gvr).
		Namespace(namespace).
		List(context.TODO(),
			metav1.ListOptions{LabelSelector: selector},
		)
}

func (d *dataSource) Get(namespace string, resourceType k8s.ResourceType, name string, subresources ...string) (*unstructured.Unstructured, error) {
	gvr, err := d.ResourceLister.GroupVersionResource(resourceType)
	if err != nil {
		return nil, err
	}
	return d.Interface.
		Resource(gvr).
		Namespace(namespace).
		Get(context.TODO(), name, metav1.GetOptions{}, subresources...)
}

func (d *dataSource) Cache() k8s.ICache { return d }

func (d *dataSource) compareObject(getObj, obj *unstructured.Unstructured, forceUpdate bool) {
	if !reflect.DeepEqual(getObj.Object["metadata"], obj.Object["metadata"]) {
		getObj.Object["metadata"] = d.compareMetadataLabelsOrAnnotation(
			getObj.Object["metadata"].(map[string]interface{}),
			obj.Object["metadata"].(map[string]interface{}),
		)
	}

	if forceUpdate {
		metadata := getObj.Object["metadata"].(map[string]interface{})
		if metadata == nil {
			goto NEXT0
		}

		annotations, exist := metadata["annotations"]
		if !exist {
			annotations = make(map[string]interface{})
		}

		annotationsMap := annotations.(map[string]interface{})
		annotationsMap["forceUpdate"] = fmt.Sprintf("%d", time.Now().Unix())
		metadata["annotations"] = annotationsMap
		getObj.Object["metadata"] = metadata
	}

NEXT0:
	if !reflect.DeepEqual(getObj.Object["spec"], obj.Object["spec"]) {
		getObj.Object["spec"] = obj.Object["spec"]
	}

	// configMap
	if !reflect.DeepEqual(getObj.Object["data"], obj.Object["data"]) {
		getObj.Object["data"] = obj.Object["data"]
	}

	if !reflect.DeepEqual(getObj.Object["binaryData"], obj.Object["binaryData"]) {
		getObj.Object["binaryData"] = obj.Object["binaryData"]
	}

	if !reflect.DeepEqual(getObj.Object["stringData"], obj.Object["stringData"]) {
		getObj.Object["stringData"] = obj.Object["stringData"]
	}

	if !reflect.DeepEqual(getObj.Object["type"], obj.Object["type"]) {
		getObj.Object["type"] = obj.Object["type"]
	}

	if !reflect.DeepEqual(getObj.Object["secrets"], obj.Object["secrets"]) {
		getObj.Object["secrets"] = obj.Object["secrets"]
	}

	if !reflect.DeepEqual(getObj.Object["imagePullSecrets"], obj.Object["imagePullSecrets"]) {
		getObj.Object["imagePullSecrets"] = obj.Object["imagePullSecrets"]
	}
	// storageClass field
	if !reflect.DeepEqual(getObj.Object["provisioner"], obj.Object["provisioner"]) {
		getObj.Object["provisioner"] = obj.Object["provisioner"]
	}

	if !reflect.DeepEqual(getObj.Object["parameters"], obj.Object["parameters"]) {
		getObj.Object["parameters"] = obj.Object["parameters"]
	}

	if !reflect.DeepEqual(getObj.Object["reclaimPolicy"], obj.Object["reclaimPolicy"]) {
		getObj.Object["reclaimPolicy"] = obj.Object["reclaimPolicy"]
	}

	if !reflect.DeepEqual(getObj.Object["volumeBindingMode"], obj.Object["volumeBindingMode"]) {
		getObj.Object["volumeBindingMode"] = obj.Object["volumeBindingMode"]
	}
}

func (d *dataSource) compareMetadataLabelsOrAnnotation(old, new map[string]interface{}) map[string]interface{} {
	newLabels, exist := new["labels"]
	if exist {
		old["labels"] = newLabels
	}
	newAnnotations, exist := new["annotations"]
	if exist {
		old["annotations"] = newAnnotations
	}

	newOwnerReferences, exist := new["ownerReferences"]
	if exist {
		old["ownerReferences"] = newOwnerReferences
	}
	return old
}
