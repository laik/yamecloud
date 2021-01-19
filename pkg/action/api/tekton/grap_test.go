package tekton

import (
	"encoding/json"
	"testing"
)

func TestGetDefaultInitData(t *testing.T) {
	data := "{\"nodes\":[{\"anchorPoints\":[[0,0.5],[1,0.5]],\"id\":\"1-1\",\"role\":0,\"taskName\":\"node-1-1\",\"x\":20,\"y\":20}]}"
	bytes, _ := json.Marshal(GetDefaultInitData())
	if data != string(bytes) {
		t.Fatal("expect not equal")
	}

}
