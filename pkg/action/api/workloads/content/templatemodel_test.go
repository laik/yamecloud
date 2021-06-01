package content

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"testing"
)

func Test_TemplateModel(t *testing.T) {
	d := NewTemplateModel()
	d.AddMetadata("devops", "my-deployment", "myuuid").
		AddContainer("abc", "nginx").
		AddCommand("start", "nginx").
		AddArgs("--config", "/etc/config").
		AddEnvironment("ENV", "DEV").
		AddEnvironment("PORT", "1024").
		AddResourceLimits(100, 200, 30, 100).
		AddVolumeMounts("etc-nginx-conf-d-mall-gw-conf", "/etc/nginx/conf.d/mall-gw.conf", "mall-gw.conf").
		AddVolumeMounts("data1", "/data1", "").
		AddVolumeMounts("data2", "/data2", "").
		SetImagePullPolicy(IfNotPresent)

	d.AddVolumes("data-www-wwwroot-config-production-app-php").
		AddConfigMap("mall-hhb-bbc", "app-php", "app.php")

	d.AddVolumes("data-www-wwwroot-config-production-app-php1").
		AddConfigMap("mall-hhb-bbc1", "app-php1", "app.php1")

	d.AddService("1", ClusterIP).AddServiceSpec(TCP, 80, 80)
	d.AddService("2", "").AddServiceSpec(TCP, 443, 443)
	d.AddService("3", "").AddServiceSpec(UDP, 443, 443)

	d.AddVolumeClaimTemplate("data1", "rook-ceph", 1024)
	d.AddVolumeClaimTemplate("data2", "rook-ceph", 1024)

	d.AddCoordinate("cg").AddReplicas(2).
		AddZoneSet("cg", "r1", "master01").
		AddZoneSet("cg", "r1", "master02").
		AddZoneSet("cg", "r1", "master03")

	x, err := yaml.Marshal(d)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("args \n %s\n", x)

	unstructuredData, err := Render(d, stoneTpl)
	if err != nil {
		t.Fatal(err)
	}

	out, err := yaml.Marshal(unstructuredData.Object)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("-------------------\n")
	fmt.Printf("%s\n", out)
}
