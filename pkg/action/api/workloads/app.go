package workloads

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yametech/yamecloud/pkg/action/api/common"
	"github.com/yametech/yamecloud/pkg/helm"
	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/releaseutil"
	helmtime "helm.sh/helm/v3/pkg/time"
	"k8s.io/apimachinery/pkg/runtime"
	k8sjson "k8s.io/apimachinery/pkg/util/json"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/scheme"
	"net/http"
	"strconv"
	"strings"
)

func (w *workloadServer) ListCharts(g *gin.Context) {
	result := make(map[string]map[string]Detail)

	_, _ = listRepoDetails("", w)

	for repoName, _ := range Repos {
		details, err := listRepoDetails(repoName, w)
		if err != nil {
			common.InternalServerError(g, err, err)
			return
		}
		detail, err := details.LastVersion()
		if err != nil {
			common.InternalServerError(g, err, err)
			return
		}
		result[repoName] = map[string]Detail{repoName: *detail}
	}
	g.JSON(http.StatusOK, result)
}

func (w *workloadServer) GetCharts(g *gin.Context) {
	repoName := g.Param("repo")
	chartName := g.Param("chart")
	version := g.Query("version")

	if len(chartName) == 0 || len(repoName) == 0 {
		common.RequestParametersError(g, fmt.Errorf("params not obtain or params parse not chart args"))
		return
	}

	details, err := getRepoDetails(repoName, chartName, w)
	if err != nil {
		common.InternalServerError(g, err, err)
		return
	}

	url, ver, err := details.GetURLAndVer(chartName, version)
	if err != nil {
		common.InternalServerError(g, err, err)
		return
	}

	archive, err := getContentByURL(repoName, chartName, ver, url)
	if err != nil {
		common.InternalServerError(g, err, err)
		return
	}

	_chart := &Chart{}
	for _, raw := range archive.Raw {
		if raw.Name == "README.md" {
			_chart.Readme = string(raw.Data)
		}
	}
	_chart.Versions = details

	g.JSON(http.StatusOK, _chart)
}

func (w *workloadServer) GetChartValues(g *gin.Context) {
	repoName := g.Param("repo")
	chartName := g.Param("chart")
	version := g.Query("version")

	if len(repoName) == 0 || len(chartName) == 0 || version == "" {
		common.RequestParametersError(g, fmt.Errorf("not enough parameters"))
		return
	}

	details, err := getRepoDetails(repoName, chartName, w)
	if err != nil {
		common.InternalServerError(g, err, err)
		return
	}

	url, ver, err := details.GetURLAndVer(chartName, version)
	if err != nil {
		common.InternalServerError(g, err, err)
		return
	}

	archive, err := getContentByURL(repoName, chartName, ver, url)
	if err != nil {
		common.InternalServerError(g, err, err)
		return
	}
	chart := &Chart{}
	for _, raw := range archive.Raw {
		if raw.Name == "values.yaml" {
			chart.Readme = string(raw.Data)
		}
	}
	chart.Versions = details

	g.JSON(http.StatusOK, chart.Readme)
}

type installChart struct {
	Chart     string      `json:"chart" form:"chart"`
	Values    helm.Values `json:"values" form:"values"`
	Name      string      `json:"name" form:"name"`
	Namespace string      `json:"namespace" form:"namespace"`
	Version   string      `json:"version" form:"version"`
}

func (w *workloadServer) InstallChart(g *gin.Context) {
	installChart := &installChart{}
	if err := g.BindJSON(installChart); err != nil {
		common.RequestParametersError(g, err)
		return
	}

	if len(installChart.Name) == 0 {
		common.RequestParametersError(g, fmt.Errorf("name is required"))
		return
	}

	repoNameAndChartName := strings.Split(installChart.Chart, "/")
	if len(repoNameAndChartName) < 2 {
		common.RequestParametersError(g, fmt.Errorf("install chart name or repo name not match"))
		return
	}

	details, err := getRepoDetails(repoNameAndChartName[0], repoNameAndChartName[1], w)
	if err != nil {
		common.InternalServerError(g, err, err)
		return
	}

	url, ver, err := details.GetURLAndVer(repoNameAndChartName[1], installChart.Version)
	if err != nil {
		common.InternalServerError(g, err, err)
		return
	}

	values, err := yaml.Marshal(installChart.Values)
	if err != nil {
		common.InternalServerError(g, err, err)
		return
	}

	archive, err := getContentByURL(repoNameAndChartName[0], repoNameAndChartName[1], ver, url)
	if err != nil {
		common.InternalServerError(g, err, err)
		return
	}

	release, err := helm.CreateRelease(w.HelmAction(installChart.Namespace), installChart.Name, installChart.Namespace, string(values), archive, nil)
	if err != nil {
		common.InternalServerError(g, err, err)
		return
	}

	g.JSON(http.StatusOK, gin.H{
		"log": release,
		"release": gin.H{
			"name": release.Name,
		},
	})

}

func (w *workloadServer) ListRelease(g *gin.Context) {
	namespace := g.Query("namespace")

	cmd := action.NewList(w.HelmAction(namespace))
	if namespace == "" {
		cmd.AllNamespaces = true
		cmd.StateMask = action.ListAll
	}
	cmd.Limit = 10000
	releasesResult, err := cmd.Run()
	if err != nil {
		common.RequestParametersError(g, err)
		return
	}

	releases := make([]Release, 0)
	for _, _release := range releasesResult {
		release := &Release{}
		release.Chart = _release.Chart.Metadata.Name + "-" + _release.Chart.Metadata.Version
		release.Name = _release.Name
		release.AppVersion = _release.Chart.AppVersion()
		release.Updated = _release.Info.LastDeployed
		release.Namespace = _release.Namespace
		release.Status = string(_release.Info.Status)
		release.Revision = strconv.Itoa(_release.Version)
		releases = append(releases, *release)
	}

	g.JSON(http.StatusOK, releases)
}

func (w *workloadServer) ReleasesByNamespace(g *gin.Context) {
	ns := g.Param("namespace")
	if len(ns) == 0 {
		common.RequestParametersError(g, fmt.Errorf("illegal namespace specified"))
		return
	}
	cmd := action.NewList(w.HelmAction(ns))
	//cmd.AllNamespaces = true
	cmd.StateMask = action.ListAll
	cmd.Limit = 10000
	releases, err := cmd.Run()
	if err != nil {
		common.RequestParametersError(g, err)
		return
	}

	rels := make([]Release, 0)
	for _, re := range releases {
		rel := new(Release)
		rel.Chart = re.Chart.Metadata.Name + "-" + re.Chart.Metadata.Version
		rel.Name = re.Name
		rel.AppVersion = re.Chart.AppVersion()
		rel.Updated = re.Info.LastDeployed
		rel.Namespace = re.Namespace
		rel.Status = string(re.Info.Status)
		rel.Revision = strconv.Itoa(re.Version)
		rels = append(rels, *rel)
	}

	g.JSON(http.StatusOK, rels)
}

type ReleaseDetail struct {
	Info      *release.Info          `json:"info,omitempty"`
	Manifest  string                 `json:"manifest,omitempty"`
	Name      string                 `json:"name,omitempty"`
	Config    map[string]interface{} `json:"config,omitempty"`
	Version   int                    `json:"version,omitempty"`
	NameSpace string                 `json:"namespace,omitempty"`
	Resources *Resources             `json:"resources,omitempty"`
}

type Resources struct {
	items []runtime.Object `json:"items,omitempty"`
}

func (w *workloadServer) ReleaseByNamespace(g *gin.Context) {
	ns := g.Param("namespace")
	name := g.Param("release")
	if len(ns) == 0 || len(name) == 0 {
		common.RequestParametersError(g, fmt.Errorf("illegal namespace name not specified"))
		return
	}

	cmd := action.NewList(w.HelmAction(ns))
	cmd.AllNamespaces = false
	cmd.Limit = 10000
	cmd.StateMask = action.ListAll
	releases, err := cmd.Run()
	if err != nil {
		common.RequestParametersError(g, err)
		return
	}

	rel := new(ReleaseDetail)
	for _, release := range releases {
		if release.Name != name {
			continue
		}

		rel.Name = name
		rel.Manifest = release.Manifest
		rel.Info = release.Info
		rel.Config = release.Config
		rel.Version = release.Version
		rel.NameSpace = release.Namespace

		f := cmdutil.NewFactory(helm.NewConfigFlagsFromCluster(release.Namespace, w.config))
		result := *f.NewBuilder().
			WithScheme(scheme.Scheme, scheme.Scheme.PrioritizedVersionsAllGroups()...).
			Stream(bytes.NewBufferString(release.Manifest), "input").
			Flatten().ContinueOnError().Do()

		if err := result.Err(); err != nil {
			common.InternalServerError(g, "", err)
			return
		}

		items, err := result.Infos()
		if err != nil {
			common.InternalServerError(g, "", err)
			return
		}
		rel.Resources = new(Resources)
		for _, item := range items {
			obj, err := w.ListGVR(item.Namespace, item.Mapping.Resource, "")
			if err != nil {
				common.InternalServerError(g, "", err)
				return
			}
			rel.Resources.items = append(rel.Resources.items, obj)
		}
		break
	}

	e, err := k8sjson.Marshal(rel.Resources.items)
	if err != nil {
		common.InternalServerError(g, err, err)
		return
	}

	g.JSON(http.StatusOK, gin.H{
		"info":      rel.Info,
		"manifest":  rel.Manifest,
		"name":      rel.Name,
		"config":    rel.Config,
		"version":   rel.Version,
		"namespace": rel.NameSpace,
		"resources": gin.H{
			"items": string(e),
		},
	})
}

func (w *workloadServer) ReleaseValueByName(g *gin.Context) {
	ns := g.Param("namespace")
	name := g.Param("release")
	if len(ns) == 0 || len(name) == 0 {
		common.RequestParametersError(g, fmt.Errorf("illegal namespace name not specified"))
		return
	}

	cmd := action.NewList(w.HelmAction(ns))
	cmd.AllNamespaces = false
	cmd.Limit = 10000
	cmd.StateMask = action.ListAll
	releases, err := cmd.Run()
	if err != nil {
		common.InternalServerError(g, "", err)
		return
	}

	rel := make(map[string]interface{})
	for _, release := range releases {
		if release.Name == name {
			rel = release.Chart.Values
		}
	}

	b, err := yaml.Marshal(rel)
	if err != nil {
		common.InternalServerError(g, "", err)
		return
	}

	g.String(http.StatusOK, string(b))
}

func (w *workloadServer) DeleteRelease(g *gin.Context) {
	ns := g.Param("namespace")
	name := g.Param("release")
	if len(ns) == 0 || len(name) == 0 {
		common.RequestParametersError(g, fmt.Errorf("illegal namespace name not specified"))
		return
	}

	_, err := helm.DeleteRelease(w.HelmAction(ns), name, false)
	if err != nil {
		common.InternalServerError(g, "", err)
		return
	}
	g.JSON(http.StatusOK, nil)

}

func (w *workloadServer) UpgradeRelease(g *gin.Context) {
	ns := g.Param("namespace")
	name := g.Param("release")
	rawData, err := g.GetRawData()
	if err != nil {
		common.RequestParametersError(g, fmt.Errorf("illegal namespace name not specified"))
		return
	}

	installChart := &installChart{}
	err = json.Unmarshal(rawData, &installChart)
	if err != nil {
		common.RequestParametersError(g, err)
		return
	}

	repoNameAndChartName := strings.Split(installChart.Chart, "/")
	if len(repoNameAndChartName) < 2 {
		common.RequestParametersError(g, fmt.Errorf("install chart name or repo name not match"))
		return
	}

	details, err := getRepoDetails(repoNameAndChartName[0], repoNameAndChartName[1], w)
	if err != nil {
		common.InternalServerError(g, err, err)
		return
	}

	url, ver, err := details.GetURLAndVer(repoNameAndChartName[1], installChart.Version)
	if err != nil {
		common.InternalServerError(g, err, err)
		return
	}

	values, err := yaml.Marshal(installChart.Values)
	if err != nil {
		common.InternalServerError(g, err, err)
		return
	}

	archive, err := getContentByURL(repoNameAndChartName[0], repoNameAndChartName[1], ver, url)
	if err != nil {
		common.InternalServerError(g, err, err)
		return
	}

	rel, err := helm.UpgradeRelease(w.HelmAction(ns), name, string(values), archive, nil)
	if err != nil {
		common.InternalServerError(g, "", err)
		return
	}

	g.JSON(http.StatusOK, gin.H{
		"log":     rel,
		"release": rel,
	})
}

type RollBackRelease struct {
	Revision int `json:"revision,omitempty"`
}

func (w *workloadServer) RollbackRelease(g *gin.Context) {
	ns := g.Param("namespace")
	name := g.Param("release")
	rawData, err := g.GetRawData()
	if err != nil {
		common.RequestParametersError(g, fmt.Errorf("illegal namespace name not specified"))
		return
	}
	obj := &RollBackRelease{}
	err = json.Unmarshal(rawData, &obj)
	if err != nil {
		common.RequestParametersError(g, err)
		return
	}
	rel, err := helm.RollbackRelease(w.HelmAction(ns), name, obj.Revision)
	if err != nil {
		common.InternalServerError(g, "", err)
		return
	}
	g.JSON(http.StatusOK, gin.H{
		"message": rel,
	})
}

func (w *workloadServer) HistoryRelease(g *gin.Context) {
	ns := g.Param("namespace")
	name := g.Param("release")
	if len(ns) == 0 || len(name) == 0 {
		common.RequestParametersError(g, fmt.Errorf("illegal namespace name not specified"))

		return
	}

	history := action.NewHistory(w.HelmAction(ns))
	history.Max = 256
	his, err := getHistory(history, name)

	if err != nil {
		common.InternalServerError(g, "", err)
		return
	}
	g.JSON(http.StatusOK, his)

}

type releaseInfo struct {
	Revision    int           `json:"revision"`
	Updated     helmtime.Time `json:"updated"`
	Status      string        `json:"status"`
	Chart       string        `json:"chart"`
	AppVersion  string        `json:"app_version"`
	Description string        `json:"description"`
}

type releaseHistory []releaseInfo

func getHistory(client *action.History, name string) (releaseHistory, error) {
	hist, err := client.Run(name)
	if err != nil {
		return nil, err
	}

	releaseutil.Reverse(hist, releaseutil.SortByRevision)

	var rels []*release.Release
	for i := 0; i < min(len(hist), client.Max); i++ {
		rels = append(rels, hist[i])
	}

	if len(rels) == 0 {
		return releaseHistory{}, nil
	}

	releaseHistory := getReleaseHistory(rels)

	return releaseHistory, nil
}

func getReleaseHistory(rls []*release.Release) (history releaseHistory) {
	for i := len(rls) - 1; i >= 0; i-- {
		r := rls[i]
		c := formatChartName(r.Chart)
		s := r.Info.Status.String()
		v := r.Version
		d := r.Info.Description
		a := formatAppVersion(r.Chart)

		rInfo := releaseInfo{
			Revision:    v,
			Status:      s,
			Chart:       c,
			AppVersion:  a,
			Description: d,
		}
		if !r.Info.LastDeployed.IsZero() {
			rInfo.Updated = r.Info.LastDeployed

		}
		history = append(history, rInfo)
	}

	return history
}

func formatChartName(c *chart.Chart) string {
	if c == nil || c.Metadata == nil {
		// This is an edge case that has happened in prod, though we don't
		// know how: https://github.com/helm/helm/issues/1347
		return "MISSING"
	}
	return fmt.Sprintf("%s-%s", c.Name(), c.Metadata.Version)
}

func formatAppVersion(c *chart.Chart) string {
	if c == nil || c.Metadata == nil {
		// This is an edge case that has happened in prod, though we don't
		// know how: https://github.com/helm/helm/issues/1347
		return "MISSING"
	}
	return c.AppVersion()
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
