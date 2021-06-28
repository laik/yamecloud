package workloads

import (
	"bytes"
	"fmt"
	"github.com/patrickmn/go-cache"
	"github.com/yametech/yamecloud/pkg/helm"
	helmchart "helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	helmtime "helm.sh/helm/v3/pkg/time"
	"io"
	"io/ioutil"
	helmhapichart "k8s.io/helm/pkg/proto/hapi/chart"
	"net/http"
	"sort"
	"time"
)

type Chart struct {
	Readme   string  `json:"readme"`
	Versions Details `json:"versions"`
}

type Detail struct {
	Annotations map[string]string `json:"annotations,omitempty"`

	APIVersion string    `json:"apiVersion"`
	AppVersion string    `json:"appVersion"`
	Created    time.Time `json:"created"`

	Dependencies []*helmchart.Dependency
	Maintainers  []*helmhapichart.Maintainer

	Description string   `json:"description"`
	Digest      string   `json:"digest"`
	Home        string   `json:"home"`
	Icon        string   `json:"icon"`
	Keywords    []string `json:"keywords"`
	Name        string   `json:"name"`
	Sources     []string `json:"sources"`
	Urls        []string `json:"urls"`
	Version     string   `json:"version"`
	Repo        string   `json:"repo"`
}

type Release struct {
	Name       string        `json:"name"`
	Namespace  string        `json:"namespace"`
	Revision   string        `json:"revision"`
	Updated    helmtime.Time `json:"updated"`
	Status     string        `json:"status"`
	Chart      string        `json:"chart"`
	AppVersion string        `json:"appVersion"`
}

var Repos = map[string]string{
	//"stabel":  "https://charts.helm.sh/stable",
	//"jenkins": "https://charts.jenkins.io/",
	//"github": "https://burdenbear.github.io/kube-charts-mirror/",
	//"aliyun": "https://kubernetes.oss-cn-hangzhou.aliyuncs.com/charts/",
	//"presslabs":          "https://presslabs.github.io/charts",
	//"banzaicloud-stable": "https://kubernetes-charts.banzaicloud.com",
	//"ingress": "https://kubernetes.github.io/ingress-nginx",
	//"fantastic-charts":   "https://fantastic-charts.storage.googleapis.com",
}

var (
	repository      = cache.New(3*time.Minute, 24*time.Hour)
	tgzContentCache = cache.New(3*time.Minute, 24*time.Hour)
)

var _ sort.Interface = &Details{}

type Details []Detail

func (d Details) Len() int { return len(d) }

func (d Details) Less(i, j int) bool { return d[i].Version < d[j].Version }

func (d Details) Swap(i, j int) { d[i], d[j] = d[j], d[i] }

func (d Details) LastVersion() (*Detail, error) {
	if d.Len() == 0 {
		return nil, fmt.Errorf("details not element")
	}
	//sort.Sort(d)
	return &d[0], nil
}

func (d Details) GetURLAndVer(chartName, version string) (url string, ver string, err error) {
	for _, _detail := range d {
		condition := true
		if version != "" {
			condition = _detail.Version == version
		}

		if condition && _detail.Name == chartName {
			if len(_detail.Urls) > 0 {
				return _detail.Urls[0], _detail.Version, nil
			}
		}
	}
	return "", "", fmt.Errorf("not found chart %s version %s", chartName, version)
}

func (d Details) Search(chartName, version string) (*Detail, error) {
	for _, _detail := range d {
		if _detail.Name == chartName && _detail.Version == version {
			return &_detail, nil
		}
	}
	return nil, fmt.Errorf("not found chart %s version %s", chartName, version)
}

func (d Details) SearchByChart(chartName string) (Details, error) {
	details := make(Details, 0)
	for _, _detail := range d {
		if _detail.Name == chartName {
			details = append(details, _detail)
		}
	}
	if details.Len() > 0 {
		return details, nil
	}
	return nil, fmt.Errorf("not found chart %s ", chartName)
}

func listRepoDetails(repoName string, w *workloadServer) (Details, error) {
	helmRepos, err := w.GetHelmRepos()
	if err != nil {
		fmt.Printf("load global config helmRepos error: %s", err)
	}

	for k, v := range helmRepos {
		Repos[k] = v
	}

	if repoName == "" {
		return nil, nil
	}

	repoURL, exist := Repos[repoName]
	if !exist {
		return nil, fmt.Errorf("repo %s not found", repoName)
	}

	detailsInterface, found := repository.Get(repoName)
	if found {
		return detailsInterface.(Details), nil
	}

	_, contentResult, err := helm.GetRepo("", "", repoURL, "")
	if err != nil {
		fmt.Printf("%s\n", err)
		return nil, nil
	}

	index, err := helm.ParseRepoIndex(contentResult)
	if err != nil {
		fmt.Printf("%s\n", err)
		return nil, nil
	}

	details := make(Details, 0)
	for _, chartVersions := range index.Entries {
		for _, value := range chartVersions {
			detail := Detail{}
			detail.Name = value.Name
			detail.Repo = repoName
			detail.APIVersion = value.ApiVersion
			detail.Version = value.Version
			detail.AppVersion = value.AppVersion
			detail.Created = value.Created
			detail.Description = value.Description
			detail.Digest = value.Digest
			detail.Home = value.Home
			detail.Icon = value.Icon
			detail.Keywords = value.Keywords
			detail.Maintainers = value.Maintainers
			detail.Sources = value.Sources
			detail.Urls = value.URLs
			details = append(details, detail)
		}
	}
	repository.Set(repoName, details, cache.DefaultExpiration)

	return details, nil
}

func getRepoDetails(repoName, chartName string, w *workloadServer) (Details, error) {
	var details Details

	detailsInterface, found := repository.Get(repoName)
	if !found {
		_details, err := listRepoDetails(repoName, w)
		if err != nil {
			return nil, err
		}
		details = _details
	} else {
		details = detailsInterface.(Details)
	}

	resultDetails, err := details.SearchByChart(chartName)
	if err != nil {
		return nil, err
	}

	return resultDetails, nil
}

func getContentByURL(repoName, chartName, version, url string) (*helmchart.Chart, error) {
	body, found := tgzContentCache.Get(fmt.Sprintf("%s-%s-%s", repoName, chartName, version))
	if found {
		return loader.LoadArchive(bytes.NewReader(body.([]byte)))
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := helm.NetClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(resp.Body)

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	tgzContentCache.Set(fmt.Sprintf("%s-%s-%s", repoName, chartName, version), data, cache.DefaultExpiration)

	return loader.LoadArchive(bytes.NewReader(data))
}
