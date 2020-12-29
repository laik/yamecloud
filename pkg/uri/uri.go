package uri

import (
	"github.com/yametech/yamecloud/pkg/k8s"
	"github.com/yametech/yamecloud/pkg/permission"
)

type Op struct {
	// Service real microservice name
	Service string `json:"service"`
	// Type Uri operator type
	Type permission.Type `json:"type"`
	// Resource kubernetes and crd resource describe
	Resource k8s.ResourceType `json:"resource"`
	// Namespace operator namespace
	Namespace string `json:"namespace"`
}

// Parser yamecloud URI general interface analysis
type Parser interface {
	ParseOp(uri string) (*Op, error)
}

func NewUriParser() Parser {
	return nil
}

var _ Parser = (*parseImplement)(nil)

type parseImplement struct{}

func (p *parseImplement) ParseOp(uri string) (*Op, error) {
	panic("implement me")
}
