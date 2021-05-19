package permission

import "github.com/yametech/yamecloud/pkg/k8s"

type OpName = string

const (
	View     OpName = "view"
	Apply    OpName = "apply"
	Delete   OpName = "delete"
	Log      OpName = "log"
	Attach   OpName = "attach"
	Annotate OpName = "annotate"
	Label    OpName = "label"
	Run      OpName = "run"
	Rerun    OpName = "rerun"
	Metrics  OpName = "metrics"
)

type TreeNode struct {
	Name     string      `json:"name"`
	Children []*TreeNode `json:"children"`
}

// Rename to Tree
var Tree = TreeNode{
	Name: "compass",
	Children: []*TreeNode{
		{
			Name: k8s.Event,
			Children: []*TreeNode{
				{Name: View},
				{Name: Apply},
				{Name: Delete},
			},
		},
		{
			Name: k8s.Node,
			Children: []*TreeNode{
				{Name: View},
				{Name: Apply},
				{Name: Delete},
				{Name: Annotate},
			},
		},
		{
			Name: "metrics",
			Children: []*TreeNode{
				{Name: Metrics},
			},
		},
		{
			Name: "workload",
			Children: []*TreeNode{
				{
					Name: k8s.WorkloadsTemplate,
					Children: []*TreeNode{
						{Name: View},
						{Name: Apply},
						{Name: Delete},
					},
				},
				{
					Name: k8s.Stone,
					Children: []*TreeNode{
						{Name: View},
						{Name: Apply},
						{Name: Delete},
					},
				},
				{
					Name: k8s.StatefulSet1,
					Children: []*TreeNode{
						{Name: View},
						{Name: Apply},
						{Name: Delete},
					},
				},
				{
					Name: k8s.Pod,
					Children: []*TreeNode{
						{Name: View},
						{Name: Apply},
						{Name: Delete},
						{Name: Log},
						{Name: Attach},
					},
				},
				{
					Name: k8s.Water,
					Children: []*TreeNode{
						{Name: View},
						{Name: Apply},
						{Name: Delete},
					},
				},
				{
					Name: k8s.Injector,
					Children: []*TreeNode{
						{Name: View},
						{Name: Apply},
						{Name: Delete},
					},
				},
				{
					Name: k8s.Deployment,
					Children: []*TreeNode{
						{Name: View},
						{Name: Apply},
						{Name: Delete},
					},
				},
				{
					Name: k8s.ReplicaSet,
					Children: []*TreeNode{
						{Name: View},
						{Name: Apply},
						{Name: Delete},
					},
				},
				{
					Name: k8s.DaemonSet,
					Children: []*TreeNode{
						{Name: View},
						{Name: Apply},
						{Name: Delete},
					},
				},
				{
					Name: k8s.StatefulSet,
					Children: []*TreeNode{
						{Name: View},
						{Name: Apply},
						{Name: Delete},
					},
				},
				{
					Name: k8s.Job,
					Children: []*TreeNode{
						{Name: View},
						{Name: Apply},
						{Name: Delete},
					},
				},
				{
					Name: k8s.CronJobs,
					Children: []*TreeNode{
						{Name: View},
						{Name: Apply},
						{Name: Delete},
					},
				},
			},
		},
		{
			Name: "tekton",
			Children: []*TreeNode{
				{
					Name: k8s.Pipeline,
					Children: []*TreeNode{
						{Name: View},
						{Name: Apply},
						{Name: Delete},
						{Name: Run},
					},
				},
				{
					Name: k8s.PipelineRun,
					Children: []*TreeNode{
						{Name: View},
						{Name: Apply},
						{Name: Delete},
						{Name: Rerun},
					},
				},
				{
					Name: k8s.PipelineResource,
					Children: []*TreeNode{
						{Name: View},
						{Name: Apply},
						{Name: Delete},
					},
				},
				{
					Name: k8s.Task,
					Children: []*TreeNode{
						{Name: View},
						{Name: Apply},
						{Name: Delete},
					},
				},
				{
					Name: k8s.TaskRun,
					Children: []*TreeNode{
						{Name: View},
						{Name: Apply},
						{Name: Delete},
					},
				},
				{
					Name: k8s.TektonGraph,
					Children: []*TreeNode{
						{Name: View},
						{Name: Apply},
						{Name: Delete},
					},
				},
				{
					Name: k8s.OpsSecret,
					Children: []*TreeNode{
						{Name: View},
						{Name: Apply},
						{Name: Delete},
					},
				},
				{
					Name: k8s.TektonStore,
					Children: []*TreeNode{
						{Name: View},
						{Name: Apply},
						{Name: Delete},
					},
				},
				{
					Name: k8s.TektonWebHook,
					Children: []*TreeNode{
						{Name: View},
						{Name: Apply},
						{Name: Delete},
					},
				},
			},
		},
		{
			Name: "servicemesh",
			Children: []*TreeNode{
				{
					Name: k8s.Gateway,
					Children: []*TreeNode{
						{Name: View},
						{Name: Apply},
						{Name: Delete},
					},
				},
				{
					Name: k8s.VirtualService,
					Children: []*TreeNode{
						{Name: View},
						{Name: Apply},
						{Name: Delete},
					},
				},
				{
					Name: k8s.DestinationRule,
					Children: []*TreeNode{
						{Name: View},
						{Name: Apply},
						{Name: Delete},
					},
				},
				{
					Name: k8s.ServiceEntry,
					Children: []*TreeNode{
						{Name: View},
						{Name: Apply},
						{Name: Delete},
					},
				},
				{
					Name: k8s.WorkloadEntry,
					Children: []*TreeNode{
						{Name: View},
						{Name: Apply},
						{Name: Delete},
					},
				},
			},
		},
		{
			Name: "configuration",
			Children: []*TreeNode{
				{
					Name: k8s.ConfigMaps,
					Children: []*TreeNode{
						{Name: View},
						{Name: Apply},
						{Name: Delete},
					},
				},
				{
					Name: k8s.Secrets,
					Children: []*TreeNode{
						{Name: View},
						{Name: Apply},
						{Name: Delete},
					},
				},
				{
					Name: k8s.HorizontalPodAutoscaler,
					Children: []*TreeNode{
						{Name: View},
						{Name: Apply},
						{Name: Delete},
					},
				},
				{
					Name: k8s.ResourceQuota,
					Children: []*TreeNode{
						{Name: View},
						{Name: Apply},
						{Name: Delete},
					},
				},
			},
		},
		{
			Name: "sdn",
			Children: []*TreeNode{
				{
					Name: k8s.Service,
					Children: []*TreeNode{
						{Name: View},
						{Name: Apply},
						{Name: Delete},
					},
				},
				{
					Name: k8s.Endpoint,
					Children: []*TreeNode{
						{Name: View},
						{Name: Apply},
						{Name: Delete},
					},
				},
				{
					Name: k8s.Ingress,
					Children: []*TreeNode{
						{Name: View},
						{Name: Apply},
						{Name: Delete},
					},
				},
				{
					Name: k8s.NetworkPolicy,
					Children: []*TreeNode{
						{Name: View},
						{Name: Apply},
						{Name: Delete},
					},
				},
			},
		},
		{
			Name: "storage",
			Children: []*TreeNode{
				{
					Name: k8s.PersistentVolume,
					Children: []*TreeNode{
						{Name: View},
						{Name: Apply},
						{Name: Delete},
					},
				},
				{
					Name: k8s.PersistentVolumeClaims,
					Children: []*TreeNode{
						{Name: View},
						{Name: Apply},
						{Name: Delete},
					},
				},
				{
					Name: k8s.StorageClass,
					Children: []*TreeNode{
						{Name: View},
						{Name: Apply},
						{Name: Delete},
					},
				},
			},
		},
		{
			Name: k8s.Namespace,
			Children: []*TreeNode{
				{Name: View},
				{Name: Apply},
				{Name: Delete},
				{Name: Annotate},
			},
		},
		{
			Name: "ovnconfig",
			Children: []*TreeNode{
				{
					Name: k8s.IP,
					Children: []*TreeNode{
						{Name: View},
						{Name: Apply},
						{Name: Delete},
					},
				},
				{
					Name: k8s.SubNet,
					Children: []*TreeNode{
						{Name: View},
						{Name: Apply},
						{Name: Delete},
					},
				},
				{
					Name: k8s.NetworkAttachmentDefinition,
					Children: []*TreeNode{
						{Name: View},
						{Name: Apply},
						{Name: Delete},
					},
				},
			},
		},
		{
			Name: "tenant",
			Children: []*TreeNode{
				{
					Name: k8s.BaseTenant,
					Children: []*TreeNode{
						{Name: View},
						{Name: Apply},
						{Name: Delete},
					},
				},
				{
					Name: k8s.BaseDepartment,
					Children: []*TreeNode{
						{Name: View},
						{Name: Apply},
						{Name: Delete},
					},
				},
				{
					Name: k8s.BaseRole,
					Children: []*TreeNode{
						{Name: View},
						{Name: Apply},
						{Name: Delete},
					},
				},
				{
					Name: k8s.BaseUser,
					Children: []*TreeNode{
						{Name: View},
						{Name: Apply},
						{Name: Delete},
					},
				},
			},
		},
		{
			Name: "accesscontrol",
			Children: []*TreeNode{
				{
					Name: k8s.ServiceAccount,
					Children: []*TreeNode{
						{Name: View},
						{Name: Apply},
						{Name: Delete},
					},
				},
				{
					Name: k8s.RoleBinding,
					Children: []*TreeNode{
						{Name: View},
						{Name: Apply},
						{Name: Delete},
					},
				},
				{
					Name: k8s.Role,
					Children: []*TreeNode{
						{Name: View},
						{Name: Apply},
						{Name: Delete},
					},
				},
				{
					Name: k8s.PodSecurityPolicie,
					Children: []*TreeNode{
						{Name: View},
						{Name: Apply},
						{Name: Delete},
					},
				},
			},
		},
	},
}
