package servicemesh

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

var _ service.IResourceService = &WorkloadEntry{}

type WorkloadEntry struct {
	service.Interface
}

func NewWorkloadEntry(svcInterface service.Interface) *WorkloadEntry {
	workloadEntry := &WorkloadEntry{Interface: svcInterface}
	svcInterface.Install(k8s.WorkloadEntry, workloadEntry)
	return workloadEntry
}

func (w *WorkloadEntry) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := w.Interface.Get(namespace, k8s.WorkloadEntry, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (w *WorkloadEntry) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := w.Interface.List(namespace, k8s.WorkloadEntry, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (w *WorkloadEntry) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	item, isUpdate, err := w.Interface.Apply(namespace, k8s.WorkloadEntry, name, unstructuredExtend)
	if err != nil {
		return nil, isUpdate, err
	}
	return item, isUpdate, nil
}

func (w *WorkloadEntry) Delete(namespace, name string) error {
	err := w.Interface.Delete(namespace, k8s.WorkloadEntry, name)
	if err != nil {
		return err
	}
	return nil
}
