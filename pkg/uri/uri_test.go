package uri

import (
	"bytes"
	"reflect"
	"testing"
)

func Test_readItem(t *testing.T) {
	const testURL = `/workload/api/v1/namespaces/im-ops`
	item, offset, err := readItem(testURL, 0)

	if err != nil || !bytes.Equal([]byte("workload"), item) || offset != 9 {
		t.Fatal("0 not expected returns")
	}

	item, offset, err = readItem(testURL, offset)

	if err != nil || !bytes.Equal([]byte("api"), item) {
		t.Fatal("1 not expected returns")
	}

	item, offset, err = readItem(testURL, offset)

	if err != nil || !bytes.Equal([]byte("v1"), item) {
		t.Fatal("2 not expected returns")
	}

	item, offset, err = readItem(testURL, offset)

	if err != nil || !bytes.Equal([]byte("namespaces"), item) {
		t.Fatal("3 not expected returns")
	}

	item, offset, err = readItem(testURL, offset)

	if err != nil || !bytes.Equal([]byte("im-ops"), item) {
		t.Fatal("4 not expected returns")
	}
}

func Test_readItem2(t *testing.T) {
	const testURL = `/workload/api/v1/pods`
	item, offset, err := readItem(testURL, 0)

	if err != nil || !bytes.Equal([]byte("workload"), item) {
		t.Fatal("0 not expected returns")
	}

	item, offset, err = readItem(testURL, offset)

	if err != nil || !bytes.Equal([]byte("api"), item) {
		t.Fatal("1 not expected returns")
	}

	item, offset, err = readItem(testURL, offset)

	if err != nil || !bytes.Equal([]byte("v1"), item) {
		t.Fatal("2 not expected returns")
	}

	item, offset, err = readItem(testURL, offset)

	if err != nil || !bytes.Equal([]byte("pods"), item) {
		t.Fatal("3 not expected returns")
	}
}

func Test_fiveURLParse(t *testing.T) {
	const testURL = `/workload/api/v1/namespaces/im-ops`
	uri, err := parse(testURL)
	if err != nil {
		t.Fatal("0 not expected returns")
	}
	if uri.Namespace != "im-ops" || uri.Service != "workload" {
		t.Fatal("1 not expected returns")
	}
}

/*
/workload/api/v1/pods
/workload/api/v1/namespaces
*/

func Test_APIFourURLParse1(t *testing.T) {
	const testURL = `/workload/api/v1/pods`
	uri, err := parse(testURL)
	if err != nil {
		t.Fatal("0 not expected returns")
	}
	if uri.Namespace != "" || uri.Service != "workload" || uri.Resource != "pods" {
		t.Fatal("1 not expected returns")
	}
}

func Test_APIFourURLParse2(t *testing.T) {
	const testURL = `/workload/api/v1/namespaces`
	uri, err := parse(testURL)
	if err != nil {
		t.Fatal("0 not expected returns")
	}
	if uri.Namespace != "" || uri.Service != "workload" || uri.Resource != "namespaces" {
		t.Fatal("1 not expected returns")
	}
}

/*
/workload/api/v1/namespaces/im
*/
func Test_APIFiveURLParse1(t *testing.T) {
	const testURL = `/workload/api/v1/namespaces/im`
	uri, err := parse(testURL)
	if err != nil {
		t.Fatal("0 not expected returns")
	}
	if uri.Namespace != "im" || uri.Service != "workload" || uri.Resource != "namespaces" {
		t.Fatal("1 not expected returns")
	}
}

/*
/workload/api/v1/namespaces/im/pods
*/
func Test_APISixURLParse(t *testing.T) {
	const testURL = `/workload/api/v1/namespaces/im/pods`
	uri, err := parse(testURL)
	if err != nil {
		t.Fatal("0 not expected returns")
	}
	if uri.Namespace != "im" || uri.Service != "workload" || uri.Resource != "pods" {
		t.Fatal("1 not expected returns")
	}
}

func Test_APISevenURLParse(t *testing.T) {
	const testURL = `/workload/api/v1/namespaces/im/pods/mypod`
	uri, err := parse(testURL)
	if err != nil {
		t.Fatal("0 not expected returns")
	}
	if uri.Namespace != "im" || uri.Service != "workload" || uri.Resource != "pods" || uri.ResourceName != "mypod" {
		t.Fatal("1 not expected returns")
	}
}

func Test_APIEightURLParse(t *testing.T) {
	const testURL = `/workload/api/v1/namespaces/im/pods/mypod/log`
	uri, err := parse(testURL)
	if err != nil {
		t.Fatal("0 not expected returns")
	}
	if uri.Namespace != "im" || uri.Service != "workload" || uri.Resource != "pods" || uri.ResourceName != "mypod" || uri.Op != "log" {
		t.Fatal("1 not expected returns")
	}
}

// /workload/apis/nuwa.nip.io/v1/stones
func Test_APISFiveURLParse1(t *testing.T) {
	const testURL = `/workload/apis/nuwa.nip.io/v1/stones`
	uri, err := parse(testURL)
	if err != nil {
		t.Fatal("0 not expected returns")
	}
	if uri.Namespace != "" || uri.Service != "workload" || uri.Resource != "stones" || uri.Version != "v1" {
		t.Fatal("1 not expected returns")
	}
}

// /workload/apis/nuwa.nip.io/v1/namespaces/im/stones
func Test_APISSixURLParse(t *testing.T) {
	const testURL = `/workload/apis/nuwa.nip.io/v1/namespaces/im/stones`
	uri, err := parse(testURL)
	if err != nil {
		t.Fatal("0 not expected returns")
	}
	if uri.Namespace != "im" || uri.Service != "workload" || uri.Resource != "stones" || uri.Version != "v1" {
		t.Fatal("1 not expected returns")
	}
}

// /workload/apis/nuwa.nip.io/v1/namespaces/im/stones/mystones
func Test_APISSevenURLParse(t *testing.T) {
	const testURL = `/workload/apis/nuwa.nip.io/v1/namespaces/im/stones/mystones`
	uri, err := parse(testURL)
	if err != nil {
		t.Fatal("0 not expected returns")
	}
	if uri.Namespace != "im" || uri.Service != "workload" || uri.Resource != "stones" || uri.Version != "v1" || uri.ResourceName != "mystones" {
		t.Fatal("1 not expected returns")
	}
}

// /workload/apis/nuwa.nip.io/v1/namespaces/im/stones/mystones/label
func Test_APISEightURLParse(t *testing.T) {
	const testURL = `/workload/apis/nuwa.nip.io/v1/namespaces/im/stones/mystones/label`
	uri, err := parse(testURL)
	if err != nil {
		t.Fatal("0 not expected returns")
	}
	if uri.Namespace != "im" || uri.Service != "workload" || uri.Resource != "stones" || uri.Version != "v1" || uri.ResourceName != "mystones" || uri.Op != "label" {
		t.Fatal("1 not expected returns")
	}
}

func Test_parseOp(t *testing.T) {
	parser := NewURIParser()

	type expectedStruct struct {
		uri string
		uriStruct
	}

	expectedStructs := []expectedStruct{
		{"/workload/api/v1/pods", uriStruct{Service: "workload", API: "api", Namespace: "", Version: "v1", Resource: "pods", ResourceName: ""}},
		{"/workload/api/v1/namespaces", uriStruct{Service: "workload", API: "api", Namespace: "", Version: "v1", Resource: "namespaces", ResourceName: ""}},
		{"/workload/api/v1/namespaces/im", uriStruct{Service: "workload", API: "api", Namespace: "im", Version: "v1", Resource: "namespaces", ResourceName: "im"}},
		{"/workload/api/v1/namespaces/im/op/label", uriStruct{Service: "workload", API: "api", Namespace: "im", Version: "v1", Resource: "namespaces", ResourceName: "im", Op: "label"}},
		{"/workload/api/v1/namespaces/im/pods", uriStruct{Service: "workload", Version: "v1", API: "api", Namespace: "im", Resource: "pods", ResourceName: ""}},
		{"/workload/api/v1/namespaces/im/pods/mypod", uriStruct{Service: "workload", Version: "v1", API: "api", Namespace: "im", Resource: "pods", ResourceName: "mypod"}},
		{"/workload/api/v1/namespaces/im/pods/mypod/op/log", uriStruct{Service: "workload", Version: "v1", API: "api", Namespace: "im", Resource: "pods", ResourceName: "mypod", Op: "log"}},
		{"/workload/apis/nuwa.nip.io/v1/stones", uriStruct{Service: "workload", Version: "v1", Group: "nuwa.nip.io", API: "apis", Namespace: "", Resource: "stones", ResourceName: ""}},
		{"/workload/apis/nuwa.nip.io/v1/namespaces/im/stones", uriStruct{Service: "workload", Version: "v1", Group: "nuwa.nip.io", API: "apis", Namespace: "im", Resource: "stones", ResourceName: ""}},
		{"/workload/apis/nuwa.nip.io/v1/namespaces/im/stones/mystones", uriStruct{Service: "workload", Version: "v1", Group: "nuwa.nip.io", API: "apis", Namespace: "im", Resource: "stones", ResourceName: "mystones"}},
		{"/workload/apis/nuwa.nip.io/v1/namespaces/im/stones/mystones/op/label", uriStruct{Service: "workload", Version: "v1", Group: "nuwa.nip.io", API: "apis", Namespace: "im", Resource: "stones", ResourceName: "mystones", Op: "label"}},
	}

	for _, value := range expectedStructs {
		if op, err := parser.ParseOp(value.uri); err != nil || !reflect.DeepEqual(&op.uriStruct, &value.uriStruct) {
			t.Fatalf("test parse uri (%s) error (%v) \nexpected (%s) \n   real （%s)", value.uri, err, &value.uriStruct, &op.uriStruct)
		}
	}
}
