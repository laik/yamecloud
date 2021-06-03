package content

type VolumeModels []VolumeModel

func (v VolumeModels) Get(name string) VolumeModel {
	for _, item := range v {
		if item.GetName() == name {
			return item
		}
	}
	return nil
}

func (v VolumeModels) Index(c VolumeModel) int {
	for index, item := range v {
		if item.GetName() == c.GetName() {
			return index
		}
	}
	return -1
}

func (v *VolumeModels) Set(c VolumeModel) {
	_index := v.Index(c)
	if _index == -1 {
		*v = append(*v, c)
	} else {
		(*v)[_index] = c
	}
}

type VolumeModel interface {
	AddConfigMap(itemName, itemKey, itemPath string)
	AddSecret(itemName, itemKey, itemPath string)
	SetName(string) VolumeModel
	GetName() string
}

var _ VolumeModel = &VolumeModelImpl{}

type VolumeModelImpl map[string]interface{}

func (v VolumeModelImpl) AddSecret(itemName, itemKey, itemPath string) {
	v["secret"] = map[string]interface{}{
		"name": itemName,
		"items": map[string]string{
			"key":  itemKey,
			"path": itemPath,
		},
	}
}

func NewVolumeModel() VolumeModel {
	return make(VolumeModelImpl)
}

func (v VolumeModelImpl) SetName(name string) VolumeModel {
	v["name"] = name
	return v
}

func (v VolumeModelImpl) AddConfigMap(itemName, itemKey, itemPath string) {
	v["configmap"] = map[string]interface{}{
		"name": itemName,
		"items": map[string]string{
			"key":  itemKey,
			"path": itemPath,
		},
	}
}

func (v VolumeModelImpl) GetName() string {
	return v["name"].(string)
}
