package content

import "fmt"

type ContainerModels []ContainerModel

func (cs ContainerModels) Get(name string) ContainerModel {
	for _, item := range cs {
		if item.GetName() == name {
			return item
		}
	}
	return nil
}

func (cs ContainerModels) Index(c ContainerModel) int {
	for index, item := range cs {
		if item.GetName() == c.GetName() {
			return index
		}
	}
	return -1
}

func (cs *ContainerModels) Set(c ContainerModel) {
	_index := cs.Index(c)
	if _index == -1 {
		*cs = append(*cs, c)
	} else {
		(*cs)[_index] = c
	}
}

type ContainerModel interface {
	AddCommand(commands ...string) ContainerModel
	AddArgs(args ...string) ContainerModel
	AddEnvironment(name, value string) ContainerModel
	AddResourceLimits(cpuLimit int64, memoryLimit int64, cpuRequest int64, memoryRequest int64) ContainerModel
	AddResourceLimits2(cpuLimit, memoryLimit, cpuRequest, memoryRequest string) ContainerModel
	AddVolumeMounts(name, mountPath, subPath string) ContainerModel
	SetImagePullPolicy(policy ImagePullPolicy) ContainerModel
	Set(name, image string)
	GetName() string
}

var _ ContainerModel = &ContainerModelImpl{}

type ContainerModelImpl map[string]interface{}

func (c ContainerModelImpl) Set(name, image string) {
	c["name"], c["image"] = name, image
}

func NewContainerModelImpl() ContainerModelImpl {
	return make(map[string]interface{})
}

func (c ContainerModelImpl) GetName() string {
	return c["name"].(string)
}

func (c ContainerModelImpl) AddCommand(commands ...string) ContainerModel {
	var _commands []string

	if __commands, exist := c["commands"]; !exist {
		_commands = make([]string, 0)
	} else {
		_commands = __commands.([]string)
	}
	_commands = append(_commands, commands...)
	c["commands"] = _commands

	return c
}

func (c ContainerModelImpl) AddArgs(args ...string) ContainerModel {
	var _args []string

	if __args, exist := c["args"]; !exist {
		_args = make([]string, 0)
	} else {
		_args = __args.([]string)
	}

	_args = append(_args, args...)
	c["args"] = _args

	return c
}

type Environments []map[string]string

func (e Environments) index(name string) int {
	for index, item := range e {
		_name, exist := item["name"]
		if !exist {
			continue
		}
		if _name == name {
			return index
		}
	}
	return -1
}

func (e *Environments) Set(name, value string) {
	_index := e.index(name)
	if _index == -1 {
		*e = append(*e, map[string]string{"name": name, "value": value})
	}
}

func (c ContainerModelImpl) AddEnvironment(name, value string) ContainerModel {
	var _environments Environments

	if __environments, exist := c["environments"]; !exist {
		_environments = make(Environments, 0)
	} else {
		_environments = __environments.(Environments)
	}

	_environments.Set(name, value)

	c["environments"] = _environments

	return c
}

func (c ContainerModelImpl) AddResourceLimits2(cpuLimit string, memoryLimit string, cpuRequest string, memoryRequest string) ContainerModel {
	var _resourceLimits map[string]string

	if __resourceLimits, exist := c["resource_limits"]; !exist {
		_resourceLimits = make(map[string]string)
	} else {
		_resourceLimits = __resourceLimits.(map[string]string)
	}

	_resourceLimits["cpu_limit"] = cpuLimit
	_resourceLimits["memory_limit"] = memoryLimit
	_resourceLimits["cpu_request"] = cpuRequest
	_resourceLimits["memory_request"] = memoryRequest

	c["resource_limits"] = _resourceLimits

	return c
}

func (c ContainerModelImpl) AddResourceLimits(cpuLimit int64, memoryLimit int64, cpuRequest int64, memoryRequest int64) ContainerModel {
	var _resourceLimits map[string]string

	if __resourceLimits, exist := c["resource_limits"]; !exist {
		_resourceLimits = make(map[string]string)
	} else {
		_resourceLimits = __resourceLimits.(map[string]string)
	}

	_resourceLimits["cpu_limit"] = fmt.Sprintf("%d", cpuLimit)
	_resourceLimits["memory_limit"] = fmt.Sprintf("%d", memoryLimit)
	_resourceLimits["cpu_request"] = fmt.Sprintf("%d", cpuRequest)
	_resourceLimits["memory_request"] = fmt.Sprintf("%d", memoryRequest)

	c["resource_limits"] = _resourceLimits

	return c
}

func (c ContainerModelImpl) SetImagePullPolicy(policy ImagePullPolicy) ContainerModel {
	c["image_pull_policy"] = policy
	return c
}

type VolumeMounts []map[string]string

func (e VolumeMounts) index(name string) int {
	for index, item := range e {
		_name, exist := item["name"]
		if !exist {
			continue
		}
		if _name == name {
			return index
		}
	}
	return -1
}

func (e *VolumeMounts) Set(name, mountPath, subPath string) {
	_index := e.index(name)
	if _index == -1 {
		*e = append(*e, map[string]string{"name": name, "mount_path": mountPath, "sub_path": subPath})
	}
}

func (c ContainerModelImpl) AddVolumeMounts(name, mountPath, subPath string) ContainerModel {
	var _volumeMounts VolumeMounts

	if __volumeMounts, exist := c["volume_mounts"]; !exist {
		_volumeMounts = make(VolumeMounts, 0)
	} else {
		_volumeMounts = __volumeMounts.(VolumeMounts)
	}

	_volumeMounts.Set(name, mountPath, subPath)

	c["volume_mounts"] = _volumeMounts

	return c
}
