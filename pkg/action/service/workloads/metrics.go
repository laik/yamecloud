package workloads

import (
	"context"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/yametech/yamecloud/common"
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

type Metrics struct {
	service.Interface
	client    *resty.Client
	InCluster bool
}

func NewMetrics(svcInterface service.Interface) *Metrics {
	_metrics := &Metrics{Interface: svcInterface, InCluster: common.InCluster}
	svcInterface.Install(k8s.Metrics, _metrics)
	return _metrics
}

const PrometheusAddress = "http://prometheus.kube-system.svc.cluster.local/api/v1/query_range"

func (m *Metrics) ProxyToPrometheus(params map[string]string, body []byte) (map[string]interface{}, error) {
	var bodyMap map[string]string
	var resultMap = make(map[string]interface{})
	err := json.Unmarshal(body, &bodyMap)
	if err != nil {
		return nil, err
	}

	if m.InCluster {
		for bodyKey, bodyValue := range bodyMap {
			resp, err := m.client.R().SetQueryParams(params).SetQueryParam("query", bodyValue).Get(PrometheusAddress)
			if err != nil {
				return nil, err
			}

			var metricsContextMap map[string]interface{}
			err = json.Unmarshal([]byte(resp.String()), &metricsContextMap)
			if err != nil {
				return nil, err
			}
			resultMap[bodyKey] = metricsContextMap
		}
		return resultMap, nil
	}

	for bodyKey, bodyValue := range bodyMap {
		req := m.ClientSet().CoreV1().RESTClient().
			Get().
			Namespace("kube-system").
			Resource("services").
			Name("prometheus:80").
			SubResource("proxy").
			Suffix("api/v1/query_range")

		for k, v := range params {
			req.Param(k, v)
		}

		req.Param("query", bodyValue)

		raw, err := req.DoRaw(context.Background())
		if err != nil {
			return nil, err
		}

		var metricsContextMap map[string]interface{}
		err = json.Unmarshal(raw, &metricsContextMap)
		if err != nil {
			return nil, err
		}
		resultMap[bodyKey] = metricsContextMap
	}

	return resultMap, nil
}

func (m Metrics) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	panic("implement me")
}

func (m Metrics) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	panic("implement me")
}

func (m Metrics) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	panic("implement me")
}
