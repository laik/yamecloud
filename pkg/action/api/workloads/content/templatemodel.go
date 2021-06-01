package content

type (
	Protocol        = string
	ImagePullPolicy = string
	ServiceType     = string
)

const (
	TCP Protocol = "tcp"
	UDP Protocol = "udp"

	IfNotPresent ImagePullPolicy = "IfNotPresent"
	Always       ImagePullPolicy = "Always"

	ClusterIP ServiceType = "ClusterIP"
	NodePort  ServiceType = "NodePort"
)

type TemplateModel interface {
	AddMetadata(namespace, name, uuid string) TemplateModel
	AddContainer(name, image string) ContainerModel
	AddVolumes(name string) VolumeModel
	AddService(name string, serviceType ServiceType) ServiceModel
	AddCoordinate(group string) CoordinateModel
	AddVolumeClaimTemplate(metadataName, storageClass string, size int64)
}

var _ TemplateModel = &TemplateModelImpl{}

type TemplateModelImpl map[string]interface{}

func NewTemplateModel() TemplateModel {
	return make(TemplateModelImpl)
}

func (d TemplateModelImpl) AddMetadata(namespace, name, uuid string) TemplateModel {
	d["namespace"], d["name"], d["uuid"] = namespace, name, uuid
	return d
}

func (d TemplateModelImpl) AddContainer(name, image string) ContainerModel {
	var containerModels ContainerModels

	if _containerModels, exist := d["containers"]; !exist {
		containerModels = make(ContainerModels, 0)
	} else {
		containerModels = _containerModels.(ContainerModels)
	}

	containerModel := containerModels.Get(name)
	if containerModel == nil {
		containerModel = NewContainerModelImpl()
	}

	containerModel.Set(name, image)
	containerModels.Set(containerModel)

	d["containers"] = containerModels

	return containerModel
}

func (d TemplateModelImpl) AddVolumes(name string) VolumeModel {
	var volumeModels VolumeModels

	if _volumeModels, exist := d["volumes"]; !exist {
		volumeModels = make(VolumeModels, 0)
	} else {
		volumeModels = _volumeModels.(VolumeModels)
	}

	volumeModel := volumeModels.Get(name)
	if volumeModel == nil {
		volumeModel = NewVolumeModel()
	}

	volumeModels.Set(volumeModel.SetName(name))
	d["volumes"] = volumeModels

	return volumeModel
}

func (d TemplateModelImpl) AddService(name string, serviceType ServiceType) ServiceModel {
	var serviceModels ServiceModels

	if _serviceModels, exist := d["service_ports"]; !exist {
		serviceModels = make(ServiceModels, 0)
	} else {
		serviceModels = _serviceModels.(ServiceModels)
	}

	serviceModel := serviceModels.Get(name)
	if serviceModel == nil {
		serviceModel = NewServiceModel()
	}
	serviceModel.AddServiceType(serviceType)

	serviceModels.Set(name, serviceModel)

	d["service_ports"] = serviceModels

	return serviceModel
}

func (d TemplateModelImpl) AddVolumeClaimTemplate(metadataName, storageClass string, size int64) {
	var volumeClaimTemplatesModels VolumeClaimTemplateModels

	if _volumeClaimTemplateModels, exist := d["volume_claim_templates"]; !exist {
		volumeClaimTemplatesModels = make(VolumeClaimTemplateModels, 0)
	} else {
		volumeClaimTemplatesModels = _volumeClaimTemplateModels.(VolumeClaimTemplateModels)
	}

	volumeClaimTemplatesModel := volumeClaimTemplatesModels.Get(metadataName)
	if volumeClaimTemplatesModel == nil {
		volumeClaimTemplatesModel = NewVolumeClaimTemplateModel()
	}
	volumeClaimTemplatesModel.AddVolumeClaimTemplate(metadataName, storageClass, size)

	volumeClaimTemplatesModels.Set(volumeClaimTemplatesModel)

	d["volume_claim_templates"] = volumeClaimTemplatesModels

	return
}

func (d TemplateModelImpl) AddCoordinate(group string) CoordinateModel {
	var coordinateModels CoordinateModels

	if _coordinateModels, exist := d["coordinates"]; !exist {
		coordinateModels = make(CoordinateModels, 0)
	} else {
		coordinateModels = _coordinateModels.(CoordinateModels)
	}

	coordinateModel := coordinateModels.Get(group)
	if coordinateModel == nil {
		coordinateModel = NewCoordinateModel()
	}
	coordinateModel.Set(group)
	coordinateModels.Set(coordinateModel)

	d["coordinates"] = coordinateModels

	return coordinateModel
}
