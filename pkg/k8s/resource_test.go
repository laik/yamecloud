package k8s

import (
	"testing"
)

func Test_Subscribe(t *testing.T) {
	//resources := []string{Service, Endpoint, Ingress}
	resources := []string{Gateway, ServiceEntry, Sidecar, ServiceAccount, DestinationRule, WorkloadEntry, Sidecar}
	subscribe := GVRMaps.Subscribe(resources...)
	if len(subscribe) != len(resources) {
		t.Fatal("expect not equal")
	}
	existUnsubscribe := false
	for _, resource := range subscribe {
		for _, item := range resources {
			if resource.Name == item {
				existUnsubscribe = false
				break
			}
			existUnsubscribe = true
		}
	}
	if existUnsubscribe {
		t.Fatal("expect not equal")
	}

}
