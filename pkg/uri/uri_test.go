package uri

import (
	"bytes"
	"net/http"
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
	const testMethod = http.MethodGet
	const testURL = `/workload/api/v1/namespaces/im-ops`
	uri, err := parse(testMethod, testURL)
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
	const testMethod = http.MethodGet
	const testURL = `/workload/api/v1/pods`
	uri, err := parse(testMethod, testURL)
	if err != nil {
		t.Fatal("0 not expected returns")
	}
	if uri.Namespace != "" || uri.Service != "workload" || uri.Resource != "pods" {
		t.Fatal("1 not expected returns")
	}
}

func Test_APIFourURLParse2(t *testing.T) {
	const testMethod = http.MethodGet
	const testURL = `/workload/api/v1/namespaces`
	uri, err := parse(testMethod, testURL)
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
	const testMethod = http.MethodGet
	const testURL = `/workload/api/v1/namespaces/im`
	uri, err := parse(testMethod, testURL)
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
	const testMethod = http.MethodGet
	const testURL = `/workload/api/v1/namespaces/im/pods`
	uri, err := parse(testMethod, testURL)
	if err != nil {
		t.Fatal("0 not expected returns")
	}
	if uri.Namespace != "im" || uri.Service != "workload" || uri.Resource != "pods" {
		t.Fatal("1 not expected returns")
	}
}

func Test_APISevenURLParse(t *testing.T) {
	const testMethod = http.MethodGet
	const testURL = `/workload/api/v1/namespaces/im/pods/mypod`
	uri, err := parse(testMethod, testURL)
	if err != nil {
		t.Fatal("0 not expected returns")
	}
	if uri.Namespace != "im" || uri.Service != "workload" || uri.Resource != "pods" || uri.ResourceName != "mypod" {
		t.Fatal("1 not expected returns")
	}
}

func Test_APIEightURLParse(t *testing.T) {
	const testMethod = http.MethodPost
	const testURL = `/workload/api/v1/namespaces/im/pods/mypod/log`
	uri, err := parse(testMethod, testURL)
	if err != nil {
		t.Fatal("0 not expected returns")
	}
	if uri.Namespace != "im" || uri.Service != "workload" || uri.Resource != "pods" || uri.ResourceName != "mypod" || uri.Op != "log" {
		t.Fatal("1 not expected returns")
	}
}

// /workload/apis/nuwa.nip.io/v1/stones
func Test_APISFiveURLParse1(t *testing.T) {
	const testMethod = http.MethodGet
	const testURL = `/workload/apis/nuwa.nip.io/v1/stones`
	uri, err := parse(testMethod, testURL)
	if err != nil {
		t.Fatal("0 not expected returns")
	}
	if uri.Namespace != "" || uri.Service != "workload" || uri.Resource != "stones" || uri.Version != "v1" {
		t.Fatal("1 not expected returns")
	}
}

// /workload/apis/nuwa.nip.io/v1/namespaces/im/stones
func Test_APISSixURLParse(t *testing.T) {
	const testMethod = http.MethodGet
	const testURL = `/workload/apis/nuwa.nip.io/v1/namespaces/im/stones`
	uri, err := parse(testMethod, testURL)
	if err != nil {
		t.Fatal("0 not expected returns")
	}
	if uri.Namespace != "im" || uri.Service != "workload" || uri.Resource != "stones" || uri.Version != "v1" {
		t.Fatal("1 not expected returns")
	}
}

// /workload/apis/nuwa.nip.io/v1/namespaces/im/stones/mystones
func Test_APISSevenURLParse(t *testing.T) {
	const testMethod = http.MethodGet
	const testURL = `/workload/apis/nuwa.nip.io/v1/namespaces/im/stones/mystones`
	uri, err := parse(testMethod, testURL)
	if err != nil {
		t.Fatal("0 not expected returns")
	}
	if uri.Namespace != "im" || uri.Service != "workload" || uri.Resource != "stones" || uri.Version != "v1" || uri.ResourceName != "mystones" {
		t.Fatal("1 not expected returns")
	}
}

// /workload/apis/nuwa.nip.io/v1/namespaces/im/stones/mystones/label
func Test_APISEightURLParse(t *testing.T) {
	const testMethod = http.MethodPost
	const testURL = `/workload/apis/nuwa.nip.io/v1/namespaces/im/stones/mystones/label`
	uri, err := parse(testMethod, testURL)
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
		method string
		uri    string
		uriStruct
	}

	expectedStructs := []expectedStruct{
		{http.MethodPost, "/workload/api/metrics?start=1615188480&end=1615192080&step=60&kubernetes_namespace=im", uriStruct{Service: "workload", API: "api", Resource: "metrics", Namespace: "im", Op: "metrics"}},
		{http.MethodPost, "/workload/metrics?start=1615188480&end=1615192080&step=60&kubernetes_namespace=im", uriStruct{Service: "workload", API: "", Resource: "metrics", Namespace: "im", Op: "metrics"}},
		{http.MethodGet, "/workload/attach/namespace/devops/pod/nginx-0-cg-0/container/nginx/common", uriStruct{Service: "workload", API: "", Resource: "pod", Namespace: "devops", Op: "attach", ResourceName: "nginx-0-cg-0"}},
		{http.MethodGet, "/workload/api/v1/pods", uriStruct{Service: "workload", API: "api", Namespace: "", Version: "v1", Resource: "pods", ResourceName: "", Op: "view"}},
		{http.MethodGet, "/workload/api/v1/namespaces", uriStruct{Service: "workload", API: "api", Namespace: "", Version: "v1", Resource: "namespaces", ResourceName: "", Op: "view"}},
		{http.MethodGet, "/workload/api/v1/namespaces/im", uriStruct{Service: "workload", API: "api", Namespace: "im", Version: "v1", Resource: "namespaces", ResourceName: "im", Op: "view"}},
		{http.MethodGet, "/workload/api/v1/namespaces/im/op/label", uriStruct{Service: "workload", API: "api", Namespace: "im", Version: "v1", Resource: "namespaces", ResourceName: "im", Op: "label"}},
		{http.MethodGet, "/workload/api/v1/namespaces/im/pods", uriStruct{Service: "workload", Version: "v1", API: "api", Namespace: "im", Resource: "pods", ResourceName: "", Op: "view"}},
		{http.MethodGet, "/workload/api/v1/namespaces/im/pods/mypod", uriStruct{Service: "workload", Version: "v1", API: "api", Namespace: "im", Resource: "pods", ResourceName: "mypod", Op: "view"}},
		{http.MethodGet, "/workload/api/v1/namespaces/im/pods/mypod/op/log", uriStruct{Service: "workload", Version: "v1", API: "api", Namespace: "im", Resource: "pods", ResourceName: "mypod", Op: "log"}},
		{http.MethodGet, "/workload/apis/nuwa.nip.io/v1/stones", uriStruct{Service: "workload", Version: "v1", Group: "nuwa.nip.io", API: "apis", Namespace: "", Resource: "stones", ResourceName: "", Op: "view"}},
		{http.MethodGet, "/workload/apis/nuwa.nip.io/v1/namespaces/im/stones", uriStruct{Service: "workload", Version: "v1", Group: "nuwa.nip.io", API: "apis", Namespace: "im", Resource: "stones", ResourceName: "", Op: "view"}},
		{http.MethodGet, "/workload/apis/nuwa.nip.io/v1/namespaces/im/stones/mystones", uriStruct{Service: "workload", Version: "v1", Group: "nuwa.nip.io", API: "apis", Namespace: "im", Resource: "stones", ResourceName: "mystones", Op: "view"}},
		{http.MethodGet, "/workload/apis/nuwa.nip.io/v1/namespaces/im/stones/mystones/op/label", uriStruct{Service: "workload", Version: "v1", Group: "nuwa.nip.io", API: "apis", Namespace: "im", Resource: "stones", ResourceName: "mystones", Op: "label"}},
	}

	for _, value := range expectedStructs {
		if op, err := parser.ParseOp(value.method, value.uri); err != nil || !reflect.DeepEqual(&op.uriStruct, &value.uriStruct) {
			t.Fatalf("test parse uri (%s) error (%v) \nexpected (%s) \n   real （%s)", value.uri, err, &value.uriStruct, &op.uriStruct)
		}
	}
}

type uriStruct2 struct {
	uriStruct
	ResourceVersion string
}

// TODO
func Test_parseWatch(t *testing.T) {
	parser := NewURIParser()

	type expectedStruct struct {
		uri         string
		uriStruct2s []uriStruct2
	}

	expectedStructs := []expectedStruct{
		{`http://compass.ym/api/watch?api=%2Fapis%2Ftekton.dev%2Fv1alpha1%2Fpipelineruns%3Fwatch%3D1%26resourceVersion%3D167634323&api=%2Fapis%2Ftekton.dev%2Fv1alpha1%2Ftaskruns%3Fwatch%3D1%26resourceVersion%3D167634323&api=%2Fapis%2Ftekton.dev%2Fv1alpha1%2Ftasks%3Fwatch%3D1%26resourceVersion%3D167634323&api=%2Fapis%2Ftekton.dev%2Fv1alpha1%2Fpipelines%3Fwatch%3D1%26resourceVersion%3D167634323&api=%2Fapis%2Ftekton.dev%2Fv1alpha1%2Fpipelineresources%3Fwatch%3D1%26resourceVersion%3D167634323&api=%2Fapi%2Fv1%2Fnamespaces%3Fwatch%3D1%26resourceVersion%3D167634332`,
			[]uriStruct2{
				{ResourceVersion: "", uriStruct: uriStruct{}},
				{ResourceVersion: "", uriStruct: uriStruct{}},
				{ResourceVersion: "", uriStruct: uriStruct{}},
				{ResourceVersion: "", uriStruct: uriStruct{}},
				{ResourceVersion: "", uriStruct: uriStruct{}},
				{ResourceVersion: "", uriStruct: uriStruct{}},
			},
		},
		{`http://compass.ym/api/watch?api=%2Fapis%2Ftekton.dev%2Fv1alpha1%2Fnamespaces%2Fsw-userops%2Fpipelineruns%3Fwatch%3D1%26resourceVersion%3D167664896&api=%2Fapis%2Ftekton.dev%2Fv1alpha1%2Fnamespaces%2Fchewu-dev%2Fpipelineruns%3Fwatch%3D1%26resourceVersion%3D167664896&api=%2Fapis%2Ftekton.dev%2Fv1alpha1%2Fnamespaces%2Fyce-cloud-extensions-ops%2Fpipelineruns%3Fwatch%3D1%26resourceVersion%3D167664896&api=%2Fapis%2Ftekton.dev%2Fv1alpha1%2Fnamespaces%2Fsw-userops%2Ftaskruns%3Fwatch%3D1%26resourceVersion%3D167664896&api=%2Fapis%2Ftekton.dev%2Fv1alpha1%2Fnamespaces%2Fchewu-dev%2Ftaskruns%3Fwatch%3D1%26resourceVersion%3D167664896&api=%2Fapis%2Ftekton.dev%2Fv1alpha1%2Fnamespaces%2Fyce-cloud-extensions-ops%2Ftaskruns%3Fwatch%3D1%26resourceVersion%3D167664896&api=%2Fapis%2Ftekton.dev%2Fv1alpha1%2Fnamespaces%2Fsw-userops%2Ftasks%3Fwatch%3D1%26resourceVersion%3D167664896&api=%2Fapis%2Ftekton.dev%2Fv1alpha1%2Fnamespaces%2Fchewu-dev%2Ftasks%3Fwatch%3D1%26resourceVersion%3D167664896&api=%2Fapis%2Ftekton.dev%2Fv1alpha1%2Fnamespaces%2Fyce-cloud-extensions-ops%2Ftasks%3Fwatch%3D1%26resourceVersion%3D167664896&api=%2Fapis%2Ftekton.dev%2Fv1alpha1%2Fnamespaces%2Fsw-userops%2Fpipelines%3Fwatch%3D1%26resourceVersion%3D167664896&api=%2Fapis%2Ftekton.dev%2Fv1alpha1%2Fnamespaces%2Fchewu-dev%2Fpipelines%3Fwatch%3D1%26resourceVersion%3D167664896&api=%2Fapis%2Ftekton.dev%2Fv1alpha1%2Fnamespaces%2Fyce-cloud-extensions-ops%2Fpipelines%3Fwatch%3D1%26resourceVersion%3D167664896&api=%2Fapis%2Ftekton.dev%2Fv1alpha1%2Fnamespaces%2Fsw-userops%2Fpipelineresources%3Fwatch%3D1%26resourceVersion%3D167664896&api=%2Fapis%2Ftekton.dev%2Fv1alpha1%2Fnamespaces%2Fchewu-dev%2Fpipelineresources%3Fwatch%3D1%26resourceVersion%3D167664896&api=%2Fapis%2Ftekton.dev%2Fv1alpha1%2Fnamespaces%2Fyce-cloud-extensions-ops%2Fpipelineresources%3Fwatch%3D1%26resourceVersion%3D167664896`,
			[]uriStruct2{
				{ResourceVersion: "", uriStruct: uriStruct{}},
				{ResourceVersion: "", uriStruct: uriStruct{}},
				{ResourceVersion: "", uriStruct: uriStruct{}},
				{ResourceVersion: "", uriStruct: uriStruct{}},
				{ResourceVersion: "", uriStruct: uriStruct{}},
				{ResourceVersion: "", uriStruct: uriStruct{}},
			},
		},
	}

	_ = parser
	_ = expectedStructs
	uris, err := parser.ParseWatch(expectedStructs[0].uri)
	if err != nil {
		t.Fatal(err)
	}

	////expectedUris := []*URI{
	////
	////}
	_ = uris

}
