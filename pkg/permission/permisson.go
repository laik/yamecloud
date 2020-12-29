package permission

type Type = uint16

const (
	//
	VIEW Type = 0x00
	//
	APPLY Type = 0x01
	//
	DELETE Type = 0x03
	// LOG pod log op
	LOG Type = 0x04
	// LOG pod attach shell op
	ATTACH Type = 0x05
	// ANNOTATE metadata annotation
	ANNOTATE Type = 0x06
	// LABEL metadata labels
	LABEL Type = 0x07
)

type OpTypeName = string

const (
	Log      OpTypeName = "log"
	Attach   OpTypeName = "attach"
	Annotate OpTypeName = "annotate"
	Label    OpTypeName = "label"
)
