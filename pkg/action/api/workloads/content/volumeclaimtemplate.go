package content

import "fmt"

type VolumeClaimTemplateModels []VolumeClaimTemplateModel

func (s VolumeClaimTemplateModels) Get(name string) VolumeClaimTemplateModel {
	for _, item := range s {
		if item.GetName() == name {
			return item
		}
	}
	return nil
}

func (s VolumeClaimTemplateModels) Index(c VolumeClaimTemplateModel) int {
	for index, item := range s {
		if item.GetName() == c.GetName() {
			return index
		}
	}
	return -1
}

func (s *VolumeClaimTemplateModels) Set(c VolumeClaimTemplateModel) {
	_index := s.Index(c)
	if _index == -1 {
		*s = append(*s, c)
	} else {
		(*s)[_index] = c
	}
}

type VolumeClaimTemplateModel interface {
	AddVolumeClaimTemplate(metadataName, storageClass string, size int64)
	GetName() string
}

var _ VolumeClaimTemplateModel = &VolumeClaimTemplateModelImpl{}

type VolumeClaimTemplateModelImpl map[string]string

func (v VolumeClaimTemplateModelImpl) GetName() string {
	return v["metadata_name"]
}

func NewVolumeClaimTemplateModel() VolumeClaimTemplateModel {
	return make(VolumeClaimTemplateModelImpl)
}

func (v VolumeClaimTemplateModelImpl) AddVolumeClaimTemplate(metadataName, storageClass string, size int64) {
	v["metadata_name"] = metadataName
	v["size"] = fmt.Sprintf("%d", size)
	v["storage_class_name"] = storageClass
}
