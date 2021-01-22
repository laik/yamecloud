package permission

type OpPosition = uint16

const (
	//
	ViewPosition OpPosition = 0x00
	//
	ApplyPosition OpPosition = 0x01
	//
	DeletePosition OpPosition = 0x03
	// LOG pod log op
	LogPosition OpPosition = 0x04
	// LOG pod attach shell op
	AttachPosition OpPosition = 0x05
	// ANNOTATE metadata annotation
	AnnotationPosition OpPosition = 0x06
	// LABEL metadata labels
	LabelPosition OpPosition = 0x07
)

type OpName = string

const (
	View     OpName = "view"
	Apply    OpName = "apply"
	Delete   OpName = "delete"
	Log      OpName = "log"
	Attach   OpName = "attach"
	Annotate OpName = "annotate"
	Label    OpName = "label"
)

var OpMap = map[OpName]OpPosition{
	View:     ViewPosition,
	Apply:    ApplyPosition,
	Delete:   DeletePosition,
	Log:      LogPosition,
	Attach:   AttachPosition,
	Annotate: AnnotationPosition,
	Label:    LabelPosition,
}
