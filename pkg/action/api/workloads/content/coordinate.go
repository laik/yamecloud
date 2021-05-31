package content

type CoordinateModels []CoordinateModel

func (s CoordinateModels) Get(name string) CoordinateModel {
	for _, item := range s {
		if item.GetName() == name {
			return item
		}
	}
	return nil
}

func (s CoordinateModels) Index(c CoordinateModel) int {
	for index, item := range s {
		if item.GetName() == c.GetName() {
			return index
		}
	}
	return -1
}

func (s *CoordinateModels) Set(c CoordinateModel) {
	_index := s.Index(c)
	if _index == -1 {
		*s = append(*s, c)
	} else {
		(*s)[_index] = c
	}
}

type CoordinateModel interface {
	AddReplicas(n int) CoordinateModel
	AddZoneSet(zone, rack, host string) CoordinateModel
	GetName() string
	Set(name string)
}

var _ CoordinateModel = &CoordinateModelImpl{}

type CoordinateModelImpl map[string]interface{}

func (c CoordinateModelImpl) Set(name string) {
	c["group"] = name
}

func NewCoordinateModel() CoordinateModel {
	return make(CoordinateModelImpl)
}

func (c CoordinateModelImpl) GetName() string {
	return c["group"].(string)
}

func (c CoordinateModelImpl) AddReplicas(n int) CoordinateModel {
	c["replicas"] = n
	return c
}

func (c CoordinateModelImpl) AddZoneSet(zone, rack, host string) CoordinateModel {
	var zoneset []map[string]string
	zoneseter, exist := c["zoneset"]
	if !exist {
		zoneset = make([]map[string]string, 0)
	} else {
		zoneset = zoneseter.([]map[string]string)
	}

	zoneset = append(zoneset, map[string]string{"zone": zone, "rack": rack, "host": host})

	c["zoneset"] = zoneset
	return c
}
