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

const ipListData = `
{
    "apiVersion": "kubeovn.io/v1",
    "kind": "IPList",
    "metadata": {
        "continue": "",
        "resourceVersion": "148628130",
        "selfLink": "/apis/kubeovn.io/v1/ips"
    },
    "items": [
        {
            "apiVersion": "kubeovn.io/v1",
            "kind": "IP",
            "metadata": {
                "creationTimestamp": "2020-11-27T01:52:40Z",
                "generation": 1,
                "labels": {
                    "ovn-default": "",
                    "ovn.kubernetes.io/subnet": "ovn-default"
                },
                "managedFields": [
                    {
                        "apiVersion": "kubeovn.io/v1",
                        "fieldsType": "FieldsV1",
                        "fieldsV1": {
                            "f:metadata": {
                                "f:labels": {
                                    ".": {},
                                    "f:ovn-default": {},
                                    "f:ovn.kubernetes.io/subnet": {}
                                }
                            },
                            "f:spec": {
                                ".": {},
                                "f:attachIps": {},
                                "f:attachMacs": {},
                                "f:attachSubnets": {},
                                "f:containerID": {},
                                "f:ipAddress": {},
                                "f:macAddress": {},
                                "f:namespace": {},
                                "f:nodeName": {},
                                "f:podName": {},
                                "f:subnet": {}
                            }
                        },
                        "manager": "kube-ovn-daemon",
                        "operation": "Update",
                        "time": "2020-11-27T01:52:40Z"
                    }
                ],
                "name": "apollo-portal-5bd55c697d-fqckg.apollo",
                "resourceVersion": "125531996",
                "selfLink": "/apis/kubeovn.io/v1/ips/apollo-portal-5bd55c697d-fqckg.apollo",
                "uid": "755b9127-a0cf-4657-ae28-dd665dbb70c6"
            },
            "spec": {
                "attachIps": [],
                "attachMacs": [],
                "attachSubnets": [],
                "containerID": "4ae8dda3dcf4e2136a764122c81abbf2b885ba5e2edf2c2d163b2047092ff612",
                "ipAddress": "10.16.15.151",
                "macAddress": "00:00:00:3A:42:C2",
                "namespace": "apollo",
                "nodeName": "node2",
                "podName": "apollo-portal-5bd55c697d-fqckg",
                "subnet": "ovn-default"
            }
        }
    ]

}`
