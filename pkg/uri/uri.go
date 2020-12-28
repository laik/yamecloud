package uri

import "github.com/yametech/yamecloud/pkg/k8s"

type Type = uint64

const (
	// LOG pod log op
	LOG Type = 0x00
	// LOG pod attach shell op
	ATTACH Type = 0x01
	// ANNOTATE metadata annotation
	ANNOTATE Type = 0x02
	// LABEL metadata labels
	LABEL Type = 0x03
)

type OpTypeName = string

const (
	Log      OpTypeName = "log"
	Attach   OpTypeName = "attach"
	Annotate OpTypeName = "annotate"
	Label    OpTypeName = "label"
)

type Op struct {
	// Service real microservice name
	Service string `json:"service"`
	// Type Uri operator type
	Type Type `json:"type"`
	// Resource kubernetes and crd resource describe
	Resource k8s.ResourceType `json:"resource"`
	// Namespace operator namespace
	Namespace string `json:"namespace"`
}

// Parser yamecloud URI general interface analysis
type Parser interface {
	ParseOp(uri string) (*Op, error)
}

var _ Parser = (*parseImplement)(nil)

type parseImplement struct{}

func (p *parseImplement) ParseOp(uri string) (*Op, error) {
	panic("implement me")
}
