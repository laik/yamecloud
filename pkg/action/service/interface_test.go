package service

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"testing"
)

func Test_SetUnstructuredExtend(t *testing.T) {
	ue := &UnstructuredExtend{
		&unstructured.Unstructured{},
	}
	if err := ue.UnmarshalJSON([]byte(stoneData)); err != nil {
		t.Fatal("problem with simulation parameters")
	}
	ue.Set("status.replicas", 2)
	_ = ue
	v, _ := ue.Get("status.replicas")
	if v.(int) != 2 {
		t.Fatal("expected not equal")
	}

}

const stoneData = `{
  "kind": "Stone",
  "apiVersion": "nuwa.nip.io/v1",
  "metadata": {
    "creationTimestamp": "2020-12-28T10:47:35Z",
    "generation": 1,
    "labels": {
      "app": "nezha-new-ui",
      "app-template-name": "nezha-new-ui-1609152449"
    },
    "managedFields": [
      {
        "apiVersion": "nuwa.nip.io/v1",
        "fieldsType": "FieldsV1",
        "fieldsV1": {
          "f:metadata": {
            "f:labels": {
              ".": {},
              "f:app": {},
              "f:app-template-name": {}
            }
          },
          "f:spec": {
            ".": {},
            "f:coordinates": {},
            "f:service": {
              ".": {},
              "f:ports": {},
              "f:type": {}
            },
            "f:strategy": {},
            "f:template": {
              ".": {},
              "f:metadata": {
                ".": {},
                "f:creationTimestamp": {},
                "f:labels": {
                  ".": {},
                  "f:app": {},
                  "f:app-template-name": {}
                },
                "f:name": {}
              },
              "f:spec": {
                ".": {},
                "f:containers": {},
                "f:imagePullSecrets": {}
              }
            }
          }
        },
        "manager": "workload",
        "operation": "Update",
        "time": "2020-12-28T10:47:35Z"
      },
      {
        "apiVersion": "nuwa.nip.io/v1",
        "fieldsType": "FieldsV1",
        "fieldsV1": {
          "f:metadata": {
            "f:annotations": {
              ".": {},
              "f:spec": {}
            }
          },
          "f:status": {
            ".": {},
            "f:replicas": {},
            "f:statefulset": {}
          }
        },
        "manager": "manager",
        "operation": "Update",
        "time": "2020-12-28T10:47:37Z"
      }
    ],
    "name": "nezha-new-ui",
    "namespace": "devops",
    "resourceVersion": "148112758",
    "selfLink": "/apis/nuwa.nip.io/v1/namespaces/devops/stones/nezha-new-ui",
    "uid": "756c3ad8-6c53-4bd0-8b0f-c2480e21ff47"
  },
  "spec": {
    "coordinates": [
      {
        "group": "B",
        "replicas": 1,
        "zoneset": [
          {
            "host": "node3",
            "rack": "W-01",
            "zone": "B"
          },
          {
            "host": "node2",
            "rack": "S-05",
            "zone": "B"
          },
          {
            "host": "node4",
            "rack": "S-02",
            "zone": "B"
          },
          {
            "host": "node3",
            "rack": "W-01",
            "zone": "B"
          }
        ]
      }
    ],
    "service": {
      "ports": [
        {
          "name": "nezha-new-ui",
          "port": 80,
          "protocol": "TCP",
          "targetPort": 80
        }
      ],
      "type": "ClusterIP"
    },
    "strategy": "Alpha",
    "template": {
      "metadata": {
        "creationTimestamp": null,
        "labels": {
          "app": "nezha-new-ui",
          "app-template-name": "nezha-new-ui-1609152449"
        },
        "name": "nezha-new-ui"
      },
      "spec": {
        "containers": [
          {
            "image": "harbor.ym/devops/devops-nezha-new-ui:v0.0.1",
            "imagePullPolicy": "Always",
            "name": "nezha-new-ui",
            "resources": {
              "limits": {
                "cpu": "300m",
                "memory": "170M"
              },
              "requests": {
                "cpu": "100m",
                "memory": "30M"
              }
            }
          }
        ],
        "imagePullSecrets": [
          {}
        ]
      }
    }
  },
  "status": {
    "replicas": 1,
    "statefulset": 1
  }
}`
