package workloads

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yametech/yamecloud/pkg/action/api/common"
	"github.com/yametech/yamecloud/pkg/action/api/workloads/content"
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/utils"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"net/http"
	"strconv"
	"strings"
)

func (s *workloadServer) GetTemplate(g *gin.Context) {
	namespace := g.Param("namespace")
	name := g.Param("name")
	if namespace == "" || name == "" {
		common.RequestParametersError(g, fmt.Errorf("params not obtain namespace=%s name=%s", namespace, name))
		return
	}

	item, err := s.Template.Get("", name)
	if err != nil {
		common.InternalServerError(g, err, err)
		return
	}
	labels := item.GetLabels()
	if value, exist := labels["namespace"]; !exist {
		g.JSON(http.StatusOK, "")
	} else if value != namespace {
		g.JSON(http.StatusOK, "")
	}

	g.JSON(http.StatusOK, item)
}

func (s *workloadServer) ListTemplate(g *gin.Context) {
	namespace := g.Param("namespace")
	labelSelector := fmt.Sprintf("namespace=%s", namespace)
	list, err := s.Template.List("", labelSelector)
	if err != nil {
		common.InternalServerError(g, "", err)
		return
	}
	g.JSON(http.StatusOK, list)
}

func (s *workloadServer) CreateTemplate(g *gin.Context) {
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
	namespace := _unstructured.GetLabels()["namespace"]
	name := fmt.Sprintf("%s.%s", _unstructured.GetName(), namespace)
	utils.Set(_unstructured.Object, "metadata.name", name)

	_unstructuredExist, _ := s.Template.Get("", name)
	if _unstructuredExist.GetName() != "" {
		common.RequestParametersError(g, fmt.Errorf("object %s already exists", name))
		return
	}

	newUnstructuredExtend, isUpdate, err := s.Template.Apply("", name, &service.UnstructuredExtend{Unstructured: _unstructured})
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

func (s *workloadServer) UpdateTemplate(g *gin.Context) {
	name := g.Param("name")
	if name == "" {
		common.RequestParametersError(g, fmt.Errorf("params not obtain name=%s", name))
		return
	}
	raw, err := g.GetRawData()
	if err != nil {
		common.RequestParametersError(g, fmt.Errorf("get raw data error (%s)", err))
		return
	}

	_unstructured := &unstructured.Unstructured{}
	if err := json.Unmarshal(raw, _unstructured); err != nil {
		common.RequestParametersError(g, fmt.Errorf("unmarshal from json data error (%s)", err))
		return
	}

	newUnstructuredExtend, _, err := s.Template.Apply("", name, &service.UnstructuredExtend{Unstructured: _unstructured})
	if err != nil {
		common.InternalServerError(g, err, err)
		return
	}
	g.JSON(
		http.StatusOK,
		[]service.UnstructuredExtend{
			*newUnstructuredExtend,
		})
}

func (s *workloadServer) LabelTemplate(g *gin.Context) {
	name := g.Param("name")
	if name == "" {
		common.RequestParametersError(g, fmt.Errorf("params not obtain name=%s", name))
		return
	}
	raw, err := g.GetRawData()
	if err != nil {
		common.RequestParametersError(g, fmt.Errorf("get raw data error (%s)", err))
		return
	}

	updateNetWorkAttachmentData := make(map[string]string)
	if err := json.Unmarshal(raw, &updateNetWorkAttachmentData); err != nil {
		common.RequestParametersError(g, fmt.Errorf("unmarshal from json data error (%s)", err))
		return
	}

	result := make([]*service.UnstructuredExtend, 0)

	for k, v := range updateNetWorkAttachmentData {
		obj, err := s.Template.LabelKV("", name, k, v)
		if err != nil {
			common.InternalServerError(g, err, err)
			return
		}
		result = append(result, obj)
	}
	g.JSON(http.StatusOK, result)
}

func (s *workloadServer) DeleteTemplate(g *gin.Context) {
	name := g.Param("name")
	if name == "" {
		common.RequestParametersError(g, fmt.Errorf("params not obtain name=%s", name))
		return
	}
	err := s.Template.Delete("", name)
	if err != nil {
		common.InternalServerError(g, err, err)
		return
	}
	g.JSON(http.StatusOK, nil)
}

type templateParameter struct {
	Data struct {
		AppName      string `json:"appName"`
		Namespace    string `json:"namespace"`
		StorageClass string `json:"storageClass"`
		Replicas     string `json:"replicas"`
		TemplateName string `json:"templateName"`

		Annotations map[string]string `json:"annotations"`
	} `json: "data"`
}

func (s *workloadServer) DeployTemplate(g *gin.Context) {
	tpParams := &templateParameter{}
	if err := g.BindJSON(tpParams); err != nil {
		common.RequestParametersError(g, fmt.Errorf("obtain parameters invalid, error: %s", err))
		return
	}
	template, err := s.Template.Get("", tpParams.Data.TemplateName)
	if err != nil {
		common.RequestParametersError(g, fmt.Errorf("can not get template %s, error: %s", tpParams.Data.TemplateName, err))
		return
	}

	resourceTypeInterface, _ := template.Get("spec.resourceType")
	resourceType, ok := resourceTypeInterface.(string)
	if resourceType == "" || !ok {
		common.RequestParametersError(g, fmt.Errorf("can not get template resource type %s, error: %s", tpParams.Data.TemplateName, err))
		return
	}

	expected := content.NewTemplateModel()
	if err := renderBaseTemplate(template, expected, tpParams.Data.Namespace); err != nil {
		common.RequestParametersError(g, fmt.Errorf("can not convert template %s, error: %s", tpParams.Data.TemplateName, err))
		return
	}

	expected.AddMetadata(tpParams.Data.Namespace, tpParams.Data.AppName, fmt.Sprintf("%s.%s", tpParams.Data.AppName, tpParams.Data.Namespace))

	var unstructuredData *unstructured.Unstructured
	var newUnstructuredObj *service.UnstructuredExtend
	switch resourceType {
	case "Stone":
		if err := renderServicesTemplate(template, expected); err != nil {
			common.RequestParametersError(g, fmt.Errorf("can not convert service template %s, error: %s", tpParams.Data.TemplateName, err))
			return
		}

		namespaceUnstructuredExtend, err := s.Namespace.Get("", tpParams.Data.Namespace)
		if err != nil {
			common.RequestParametersError(g, fmt.Errorf("can not get namespace %s, error: %s", tpParams.Data.Namespace, err))
			return
		}

		coordinatesStr, ok := namespaceUnstructuredExtend.GetAnnotations()["nuwa.kubernetes.io/default_resource_limit"]
		if !ok {
			common.RequestParametersError(g, fmt.Errorf("can not get allowed node on namespace %s", tpParams.Data.Namespace))
			return
		}

		coordinates := make([]map[string]interface{}, 0)
		if err := json.Unmarshal([]byte(coordinatesStr), &coordinates); err != nil {
			common.RequestParametersError(g, fmt.Errorf("can not unmarshal get allowed node to namespace %s", tpParams.Data.Namespace))
			return
		}

		if err := renderCoordinatesTemplate(coordinates, expected, tpParams.Data.Replicas); err != nil {
			common.RequestParametersError(g, fmt.Errorf("can not convert coordinates template %s, error: %s", tpParams.Data.TemplateName, err))
			return
		}

		if err := renderVolumeClaimsTemplate(template, expected, tpParams.Data.StorageClass); err != nil {
			common.RequestParametersError(g, fmt.Errorf("can not convert volume claim template %s, error: %s", tpParams.Data.TemplateName, err))
			return
		}

		unstructuredData, err = content.Render(expected, content.StoneTpl)
		if err != nil {
			common.RequestParametersError(g, fmt.Errorf("render stone template %s, error: %s", tpParams.Data.TemplateName, err))
			return
		}

		newUnstructuredObj, _, err = s.Template.CreateStone(tpParams.Data.Namespace, unstructuredData.GetName(), &service.UnstructuredExtend{Unstructured: unstructuredData})
		if err != nil {
			common.InternalServerError(g, err, fmt.Errorf("create stone error: %v", err))
			return
		}
	case "Deployment":

	}

	g.JSON(http.StatusOK, newUnstructuredObj)

	return
}

func renderVolumeClaimsTemplate(extend *service.UnstructuredExtend, expected content.TemplateModel, storageClassName string) error {
	volumeClaimsStr, _ := extend.Get("spec.volumeClaims")
	volumeClaims := make([]map[string]interface{}, 0)

	err := json.Unmarshal([]byte(volumeClaimsStr.(string)), &volumeClaims)
	if err != nil {
		return fmt.Errorf("convert volume claims error, value: %v", volumeClaims)
	}

	if len(volumeClaims) > 0 && storageClassName == "" {
		return fmt.Errorf("the volume claim request is defined, but the storageClass is not statement")
	}

	for _, volumeClaim := range volumeClaims {
		metadataName := utils.Get(volumeClaim, "metadata.name").(string)
		size := utils.Get(volumeClaim, "spec.resources.requests.storage").(string)
		_size, err := strconv.ParseUint(size, 10, 64)
		if err != nil {
			return fmt.Errorf("parse volume claims size error %s, value: %v", err, volumeClaims)
		}
		expected.AddVolumeClaimTemplate(metadataName, storageClassName, int64(_size))
	}

	return nil
}

func renderCoordinatesTemplate(coordinates []map[string]interface{}, expected content.TemplateModel, replicas string) error {
	groupCoordinates := make(map[string][]map[string]string)
	for _, coordinate := range coordinates {
		group := utils.Get(coordinate, "zone").(string)
		rack := utils.Get(coordinate, "rack").(string)
		host := utils.Get(coordinate, "host").(string)
		if _, exist := groupCoordinates[group]; !exist {
			groupCoordinates[group] = make([]map[string]string, 0)
		}
		groupCoordinates[group] = append(groupCoordinates[group], map[string]string{"zone": group, "rack": rack, "host": host})
	}

	for group, gcs := range groupCoordinates {
		coordinateModel := expected.AddCoordinate(group)
		_replicas, err := strconv.ParseInt(replicas, 10, 64)
		if err != nil {
			return fmt.Errorf("parse replicas error: %v", err)
		}
		coordinateModel.AddReplicas(int(_replicas))
		for _, v := range gcs {
			coordinateModel.AddZoneSet(v["zone"], v["rack"], v["host"])
		}
	}

	return nil
}

func exclude(portName string) string {
	needChange := false
	for _, item := range []string{"port", "http", "https"} {
		if strings.ToLower(portName) == item {
			needChange = true
		}
	}
	if needChange {
		return ""
	}
	return portName
}

func renderServicesTemplate(extend *service.UnstructuredExtend, expected content.TemplateModel) error {
	serviceStr, _ := extend.Get("spec.service")

	services := make(map[string]interface{})

	err := json.Unmarshal([]byte(serviceStr.(string)), &services)
	if err != nil {
		return fmt.Errorf("convert services error, value: %v", services)
	}

	serviceType := utils.Get(services, "type").(string)
	servicePorts := utils.Get(services, "ports").([]interface{})

	for _, servicePortMap := range servicePorts {
		servicePort := servicePortMap.(map[string]interface{})

		name := utils.Get(servicePort, "name").(string)
		serviceModel := expected.AddService(exclude(name), serviceType)

		protocol := utils.Get(servicePort, "protocol").(string)
		port := utils.Get(servicePort, "port").(string)

		targetPort := utils.Get(servicePort, "targetPort").(string)
		serviceModel.AddServiceSpec2(protocol, port, targetPort)
	}

	return nil
}

func renderBaseTemplate(extend *service.UnstructuredExtend, expected content.TemplateModel, deployNamespace string) error {
	metadata, _ := extend.Get("spec.metadata")
	containers := make([]map[string]interface{}, 0)

	err := json.Unmarshal([]byte(metadata.(string)), &containers)
	if err != nil {
		return fmt.Errorf("convert containers error, value: %v", metadata)
	} else if len(containers) == 0 {
		return fmt.Errorf("containers not item, value: %v", containers)
	}

	for _, _container := range containers {

		name := utils.Get(_container, "base.name").(string)
		image := utils.Get(_container, "base.image").(string)
		if image == "" {
			return fmt.Errorf("container %s not define image", name)
		}

		imageFrom := utils.Get(_container, "base.image").(string)
		switch imageFrom {
		case "private":
			imagePullSecret := utils.Get(_container, "base.imagePullSecret").(string)
			imagePullSecretNamespace := utils.Get(_container, "base.imagePullSecretNamespace").(string)
			if deployNamespace != imagePullSecretNamespace {
				return fmt.Errorf("container pull use secret %s but deploy namespace %s not equal to defined namespace %s", imagePullSecret, deployNamespace, imagePullSecretNamespace)
			}
			expected.AddImagePullSecrets(imagePullSecret)
		case "public":
		}

		imagePullPolicy := utils.Get(_container, "base.imagePullPolicy").(string)
		cpuLimit := utils.Get(_container, "base.resource.limits.cpu").(string)
		memoryLimit := utils.Get(_container, "base.resource.limits.memory").(string)
		cpuRequest := utils.Get(_container, "base.resource.requests.cpu").(string)
		memoryRequest := utils.Get(_container, "base.resource.requests.memory").(string)

		containerModel := expected.
			AddContainer(name, image).
			SetImagePullPolicy(imagePullPolicy).
			AddResourceLimits2(cpuLimit, memoryLimit, cpuRequest, memoryRequest)

		iCommands := utils.Get(_container, "commands").([]interface{})
		for _, iCommand := range iCommands {
			containerModel.AddCommand(iCommand.(string))
		}

		iArgs := utils.Get(_container, "args").([]interface{})
		for _, iArg := range iArgs {
			containerModel.AddArgs(iArg.(string))
		}

		environments := utils.Get(_container, "environment").([]interface{})
		for _, iEnvironment := range environments {
			environment := iEnvironment.(map[string]interface{})
			envType := utils.Get(environment, "type").(string)
			switch envType {
			case "Normal":
				name := utils.Get(environment, "envConfig.name").(string)
				value := utils.Get(environment, "envConfig.value").(string)
				containerModel.AddEnvironment(name, value)
			default:
			}
		}

		volumeMounts := utils.Get(_container, "volumeMounts").(map[string]interface{})
		volumeMountsItems, ok := utils.Get(volumeMounts, "items").([]interface{})
		if !ok {
			return fmt.Errorf("containers volume mounts get items error, value: %v", volumeMountsItems)
		}

		for _, volumeMountsItem := range volumeMountsItems {
			volumeMount := volumeMountsItem.(map[string]interface{})
			mountType := utils.Get(volumeMount, "mountType").(string)

			switch mountType {
			case "VolumeClaim":
				name := utils.Get(volumeMount, "mountConfig.name").(string)
				mountPath := utils.Get(volumeMount, "mountConfig.mountPath").(string)

				containerModel.AddVolumeMounts(name, mountPath, "")

			case "ConfigMaps":
				mountPath := utils.Get(volumeMount, "mountConfig.mountPath").(string)
				subPath := utils.Get(volumeMount, "mountConfig.subPath").(string)
				namespace := utils.Get(volumeMount, "mountConfig.namespace").(string)
				if deployNamespace != namespace {
					return fmt.Errorf("deploy namespace %s not equal defined configmap namespace %s", deployNamespace, namespace)
				}

				containerModel.AddVolumeMounts(replaceName(mountPath), mountPath, subPath)

				configName := utils.Get(volumeMount, "mountConfig.configName").(string)
				configKey := utils.Get(volumeMount, "mountConfig.configKey").(string)

				expected.AddVolumes(replaceName(mountPath)).
					AddConfigMap(configName, configKey, configKey)

			case "Secrets":
				mountPath := utils.Get(volumeMount, "mountConfig.mountPath").(string)
				subPath := utils.Get(volumeMount, "mountConfig.subPath").(string)
				namespace := utils.Get(volumeMount, "mountConfig.namespace").(string)
				if deployNamespace != namespace {
					return fmt.Errorf("deploy namespace %s not equal defined secret namespace %s", deployNamespace, namespace)
				}

				containerModel.AddVolumeMounts(replaceName(mountPath), mountPath, subPath)

				secretName := utils.Get(volumeMount, "mountConfig.secretName").(string)
				secretKey := utils.Get(volumeMount, "mountConfig.secretKey").(string)

				expected.AddVolumes(replaceName(mountPath)).
					AddSecret(secretName, secretKey, secretKey)
			}
		}
	}

	return nil
}

func replaceName(path string) string {
	return strings.TrimPrefix(strings.Replace(strings.Replace(path, "/", "-", -1), ".", "-", -1), "-")
}

//
//workloadsTemplate := &workloadsTemplate{
//	Metadata:     make(metadataTemplate, 0),
//	Service:      serviceTemplate{},
//	VolumeClaims: make(volumeClaimsTemplate, 0),
//}
//if err := json.Unmarshal([]byte(*workloads.Spec.Metadata), &workloadsTemplate.Metadata); err != nil {
//	common.ToRequestParamsError(g, err)
//	return
//}
//
//if err := json.Unmarshal([]byte(*workloads.Spec.Service), &workloadsTemplate.Service); err != nil {
//	common.ToRequestParamsError(g, err)
//	return
//}
//if err := json.Unmarshal([]byte(workloads.Spec.VolumeClaims), &workloadsTemplate.VolumeClaims); err != nil {
//	common.ToRequestParamsError(g, err)
//	return
//}
//
//var runtimeObj runtime.Object

//
//	var runtimeClassGVR schema.GroupVersionResource
//	switch *workloads.Spec.ResourceType {
//	case "Stone":
//		obj, err := w.namespace.Get("", deployTemplate.Namespace)
//		if err != nil {
//			common.ToRequestParamsError(g, err)
//			return
//		}
//		namespaceUnstructured := obj.(*unstructured.Unstructured)
//		namespace := &corev1.Namespace{}
//		if err = runtime.DefaultUnstructuredConverter.FromUnstructured(namespaceUnstructured.Object, namespace); err != nil {
//			common.ToRequestParamsError(g, err)
//			return
//		}
//
//		notResourceAllocatedError := fmt.Errorf("node resources are not allocated in this namespace, please contact the administrator")
//		annotations := namespace.GetAnnotations()
//		if annotations == nil {
//			common.ToRequestParamsError(g,
//				notResourceAllocatedError,
//			)
//			return
//		}
//		limits, ok := annotations[constraint.NamespaceAnnotationForNodeResource]
//		if !ok {
//			common.ToRequestParamsError(g,
//				notResourceAllocatedError,
//			)
//			return
//		}
//
//		cds := make(nuwav1.Coordinates, 0)
//		err = json.Unmarshal([]byte(limits), &cds)
//		if err != nil {
//			common.ToRequestParamsError(g, err)
//			return
//		}
//		if len(cds) == 0 {
//			common.ToRequestParamsError(g,
//				notResourceAllocatedError,
//			)
//			return
//		}
//		replicas, err := strconv.ParseUint(deployTemplate.Replicas, 10, 32)
//		if err != nil {
//			common.ToRequestParamsError(g, err)
//			return
//		}
//		rs := int32(replicas)
//		_, cgs := groupBy(cds, rs)
//
//		serviceTemplate, err := workloadsTemplateToServiceSpec(workloadsTemplate)
//		if err != nil {
//			if err != nil {
//				common.ToRequestParamsError(g, err)
//				return
//			}
//		}
//
//		// PodSpec
//		podSpec := corev1.PodSpec{}
//		podSpec.Containers = workloadsTemplateToPodContainers(workloadsTemplate)
//		podSpec.ImagePullSecrets = workloadsTemplateImagePullSecrets(workloadsTemplate)
//
//		// Labels
//		labels := map[string]string{
//			"app":               deployTemplate.AppName,
//			"app-template-name": deployTemplate.TemplateName,
//		}
//		runtimeObj = &nuwav1.Stone{
//			TypeMeta: metav1.TypeMeta{
//				Kind:       "Stone",
//				APIVersion: "nuwa.nip.io/v1",
//			},
//			ObjectMeta: metav1.ObjectMeta{
//				Name:        deployTemplate.AppName,
//				Namespace:   deployTemplate.Namespace,
//				Labels:      labels,
//				Annotations: deployTemplate.Annotations,
//			},
//			Spec: nuwav1.StoneSpec{
//				Template: corev1.PodTemplateSpec{
//					ObjectMeta: metav1.ObjectMeta{
//						Name:   deployTemplate.AppName,
//						Labels: labels,
//					},
//					Spec: podSpec,
//				},
//				Strategy:             "Alpha", // TODO
//				Coordinates:          cgs,
//				Service:              *serviceTemplate,
//				VolumeClaimTemplates: workloadsTemplateToVolumeClaims(workloadsTemplate, deployTemplate),
//			},
//		}
//		runtimeClassGVR = types.ResourceStone
//	}
//
//	unstructuredObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(runtimeObj)
//	if err != nil {
//		common.ToRequestParamsError(g, err)
//		return
//	}
//	unstructuredData := &unstructured.Unstructured{Object: unstructuredObj}
//
//	w.generic.SetGroupVersionResource(runtimeClassGVR)
//	newObj, _, err := w.generic.Apply(deployTemplate.Namespace, deployTemplate.AppName, unstructuredData)
//	if err != nil {
//		common.ToInternalServerError(g, unstructuredData, err)
//		return
//	}
//	g.JSON(
//		http.StatusOK,
//		[]unstructured.Unstructured{
//			*newObj,
//		})
//
//	return
//}

//
//func workloadsTemplateToServiceSpec(wt *workloadsTemplate) (*corev1.ServiceSpec, error) {
//	serviceSpec := &corev1.ServiceSpec{
//		Type: corev1.ServiceType(wt.Service.Type),
//	}
//	for _, item := range wt.Service.Ports {
//		port, err := strconv.ParseInt(item.Port, 10, 32)
//		if err != nil {
//			return nil, err
//		}
//		newItem := corev1.ServicePort{
//			Name:       item.Name,
//			Protocol:   corev1.Protocol(item.Protocol),
//			Port:       int32(port),
//			TargetPort: intstr.Parse(item.TargetPort),
//		}
//		serviceSpec.Ports = append(serviceSpec.Ports, newItem)
//	}
//	return serviceSpec, nil
//}
//
//func string2int32(s string) int32 {
//	i, _ := strconv.ParseInt(s, 10, 32)
//	return int32(i)
//}
//
//func workloadsTemplateToPodContainers(wt *workloadsTemplate) []corev1.Container {
//	containers := make([]corev1.Container, 0)
//	for _, item := range wt.Metadata {
//		volumeMounts := make([]corev1.VolumeMount, 0)
//		for _, subItem := range item.VolumeMounts.Items {
//			volumeMounts = append(volumeMounts,
//				corev1.VolumeMount{
//					Name:      subItem.Name,
//					MountPath: subItem.MountPath,
//				})
//		}
//		container := corev1.Container{
//			Name:            item.Base.Name,
//			Image:           item.Base.Image,
//			ImagePullPolicy: corev1.PullPolicy(item.Base.ImagePullPolicy),
//			Resources: corev1.ResourceRequirements{
//				Limits: corev1.ResourceList{
//					corev1.ResourceName("cpu"): resource.MustParse(
//						fmt.Sprintf("%s", item.Base.Resource.Limits.CPU),
//					),
//					corev1.ResourceName("memory"): resource.MustParse(
//						fmt.Sprintf("%sM", item.Base.Resource.Limits.Memory),
//					),
//				},
//				Requests: corev1.ResourceList{
//					corev1.ResourceName("cpu"): resource.MustParse(
//						fmt.Sprintf("%s", item.Base.Resource.Requests.CPU)),
//					corev1.ResourceName("memory"): resource.MustParse(
//						fmt.Sprintf("%sM", item.Base.Resource.Requests.Memory),
//					),
//				},
//			},
//			VolumeMounts: volumeMounts,
//		}
//		// Command
//		for _, cmd := range item.Commands {
//			container.Command = append(container.Command, cmd)
//		}
//		// Args
//		for _, arg := range item.Args {
//			container.Args = append(container.Args, arg)
//		}
//		// Environment TODO
//		envs := make([]corev1.EnvVar, 0)
//		for _, evnConfig := range item.Environment {
//			switch evnConfig.Type {
//			case "ConfigMaps":
//				env := corev1.EnvVar{
//					Name: evnConfig.EnvConfig.Name,
//					ValueFrom: &corev1.EnvVarSource{
//						ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
//							LocalObjectReference: corev1.LocalObjectReference{Name: evnConfig.EnvConfig.ConfigName},
//							Key:                  evnConfig.EnvConfig.ConfigKey,
//						},
//					},
//				}
//				envs = append(envs, env)
//			case "Secrets":
//				env := corev1.EnvVar{
//					Name: evnConfig.EnvConfig.Name,
//					ValueFrom: &corev1.EnvVarSource{
//						SecretKeyRef: &corev1.SecretKeySelector{
//							LocalObjectReference: corev1.LocalObjectReference{Name: evnConfig.EnvConfig.ConfigName},
//							Key:                  evnConfig.EnvConfig.ConfigKey,
//						},
//					},
//				}
//				envs = append(envs, env)
//			//case "Other":
//			//case "CUSTOM":
//			case "Normal":
//				env := corev1.EnvVar{Name: evnConfig.EnvConfig.Name, Value: evnConfig.EnvConfig.Value}
//				envs = append(envs, env)
//			}
//		}
//		container.Env = envs
//
//		// readinessProbe
//		if item.ReadyProbe.Status {
//			container.ReadinessProbe = &corev1.Probe{
//				Handler:             corev1.Handler{},
//				InitialDelaySeconds: string2int32(item.ReadyProbe.Delay),
//				TimeoutSeconds:      string2int32(item.ReadyProbe.Timeout),
//				PeriodSeconds:       string2int32(item.ReadyProbe.Cycle),
//				FailureThreshold:    string2int32(item.ReadyProbe.RetryCount),
//			}
//			if item.ReadyProbe.Pattern.Command != "" {
//				container.ReadinessProbe.Handler.Exec = &corev1.ExecAction{Command: []string{item.ReadyProbe.Pattern.Command}}
//			}
//			if item.ReadyProbe.Pattern.URL != "" && item.ReadyProbe.Pattern.HTTPPort != "" {
//				container.ReadinessProbe.Handler.HTTPGet = &corev1.HTTPGetAction{
//					Path: item.ReadyProbe.Pattern.URL,
//					Port: intstr.Parse(item.ReadyProbe.Pattern.HTTPPort),
//				}
//			}
//			if item.ReadyProbe.Pattern.TCPPort != "" {
//				container.ReadinessProbe.Handler.TCPSocket = &corev1.TCPSocketAction{
//					Port: intstr.Parse(item.ReadyProbe.Pattern.TCPPort),
//				}
//			}
//		}
//
//		// livenessProbe
//		if item.LiveProbe.Status {
//			container.LivenessProbe = &corev1.Probe{
//				Handler:             corev1.Handler{},
//				InitialDelaySeconds: string2int32(item.LiveProbe.Delay),
//				TimeoutSeconds:      string2int32(item.LiveProbe.Timeout),
//				PeriodSeconds:       string2int32(item.LiveProbe.Cycle),
//				FailureThreshold:    string2int32(item.LiveProbe.RetryCount),
//			}
//
//			if item.LiveProbe.Pattern.Command != "" {
//				container.LivenessProbe.Handler.Exec = &corev1.ExecAction{Command: []string{item.LiveProbe.Pattern.Command}}
//			}
//			if item.LiveProbe.Pattern.URL != "" && item.LiveProbe.Pattern.HTTPPort != "" {
//				container.LivenessProbe.Handler.HTTPGet = &corev1.HTTPGetAction{
//					Path: item.LiveProbe.Pattern.URL,
//					Port: intstr.Parse(item.LiveProbe.Pattern.HTTPPort),
//				}
//			}
//			if item.LiveProbe.Pattern.TCPPort != "" {
//				container.LivenessProbe.Handler.TCPSocket = &corev1.TCPSocketAction{
//					Port: intstr.Parse(item.LiveProbe.Pattern.TCPPort),
//				}
//			}
//		}
//		// LifeCycle
//		if item.LifeCycle.Status {
//			container.Lifecycle = &corev1.Lifecycle{
//				PostStart: &corev1.Handler{},
//				PreStop:   &corev1.Handler{},
//			}
//			if item.LifeCycle.PostStart.Command != "" {
//				container.Lifecycle.PostStart.Exec = &corev1.ExecAction{
//					Command: []string{item.LifeCycle.PostStart.Command},
//				}
//			}
//			if item.LifeCycle.PostStart.URL != "" && item.LifeCycle.PostStart.TCPPort != "" {
//				container.Lifecycle.PostStart.HTTPGet = &corev1.HTTPGetAction{
//					Path: item.LifeCycle.PostStart.URL,
//					Port: intstr.Parse(item.LifeCycle.PostStart.HTTPPort),
//				}
//			}
//			if item.LifeCycle.PostStart.TCPPort != "" {
//				container.Lifecycle.PostStart.TCPSocket = &corev1.TCPSocketAction{
//					Port: intstr.Parse(item.LifeCycle.PostStart.TCPPort),
//				}
//			}
//
//			if item.LifeCycle.PreStop.Command != "" {
//				container.Lifecycle.PreStop.Exec = &corev1.ExecAction{
//					Command: []string{item.LifeCycle.PostStart.Command},
//				}
//			}
//			if item.LifeCycle.PreStop.URL != "" && item.LifeCycle.PreStop.TCPPort != "" {
//				container.Lifecycle.PreStop.HTTPGet = &corev1.HTTPGetAction{
//					Path: item.LifeCycle.PreStop.URL,
//					Port: intstr.Parse(item.LifeCycle.PostStart.HTTPPort),
//				}
//			}
//			if item.LifeCycle.PreStop.TCPPort != "" {
//				container.Lifecycle.PreStop.TCPSocket = &corev1.TCPSocketAction{
//					Port: intstr.Parse(item.LifeCycle.PreStop.TCPPort),
//				}
//			}
//
//		}
//		containers = append(containers, container)
//	}
//	return containers
//}
//
//func workloadsTemplateToVolumeClaims(wt *workloadsTemplate, dt *deployTemplate) []corev1.PersistentVolumeClaim {
//	pvcs := make([]corev1.PersistentVolumeClaim, 0)
//	for _, item := range wt.VolumeClaims {
//		pvc := corev1.PersistentVolumeClaim{
//			ObjectMeta: metav1.ObjectMeta{
//				Name: item.Metadata.Name,
//			},
//			Spec: corev1.PersistentVolumeClaimSpec{},
//		}
//		accessModes := make([]corev1.PersistentVolumeAccessMode, 0)
//		for _, subItem := range item.Spec.AccessModes {
//			accessModes = append(accessModes, corev1.PersistentVolumeAccessMode(subItem))
//		}
//		resourceRequire := corev1.ResourceRequirements{
//			Requests: corev1.ResourceList{
//				corev1.ResourceName("storage"): resource.MustParse(item.Spec.Resources.Requests.Storage + "Mi"),
//			},
//		}
//		//if item.Metadata.IsUseDefaultStorageClass {
//		//	pvc.ObjectMeta.Annotations = map[string]string{
//		//		"volume.alpha.kubernetes.io/storage-class": "default",
//		//	}
//		//	pvc.Spec = corev1.PersistentVolumeClaimSpec{
//		//		AccessModes: accessModes,
//		//		Resources:   resourceRequire,
//		//	}
//		//} else {
//		pvc.Spec = corev1.PersistentVolumeClaimSpec{
//			StorageClassName: &dt.StorageClass,
//			AccessModes:      accessModes,
//			Resources:        resourceRequire,
//		}
//		//}
//		pvcs = append(pvcs, pvc)
//	}
//	return pvcs
//}
//
//func groupBy(cds nuwav1.Coordinates, replicas int32) (group int, result []nuwav1.CoordinatesGroup) {
//	var temp = make(map[string]nuwav1.Coordinates)
//	for _, coordinate := range cds {
//		if _, ok := temp[coordinate.Zone]; !ok {
//			temp[coordinate.Zone] = make(nuwav1.Coordinates, 0)
//		}
//		temp[coordinate.Zone] = append(temp[coordinate.Zone], coordinate)
//	}
//	for k, v := range temp {
//		result = append(result, nuwav1.CoordinatesGroup{
//			Group:    k,
//			Zoneset:  v,
//			Replicas: &replicas,
//		})
//	}
//	return len(temp), result
//}
//
//func workloadsTemplateImagePullSecrets(w *workloadsTemplate) []corev1.LocalObjectReference {
//	result := make([]corev1.LocalObjectReference, 0)
//	for _, item := range w.Metadata {
//		result = append(result, corev1.LocalObjectReference{Name: item.Base.ImagePullSecret})
//	}
//	return result
//}
