package gateway

import "encoding/json"

type userConfig struct {
	LensVersion        string   `json:"lensVersion"`
	LensTheme          string   `json:"lensTheme"`
	UserName           string   `json:"userName"`
	Token              string   `json:"token"`
	AllowedNamespaces  []string `json:"allowedNamespaces"`
	IsClusterAdmin     bool     `json:"isClusterAdmin"`
	IsTenantOwner      bool     `json:"isTenantOwner"`
	IsCompassDeveloper bool     `json:"isCompassDeveloper"`
	IsDepartmentOwner  bool     `json:"isDepartmentOwner"`
	ChartEnable        bool     `json:"chartEnable"`
	KubectlAccess      bool     `json:"kubectlAccess"`
	DefaultNamespace   string   `json:"defaultNamespace"`
}

func (uc *userConfig) String() string {
	bytesData, _ := json.Marshal(uc)
	return string(bytesData)
}

func NewUserConfig(user string, token string, allowedNamespaces []string, defaultNamespace string) *userConfig {
	isClusterAdmin := false
	if user == "admin" {
		isClusterAdmin = true
	}
	return &userConfig{
		LensVersion:       "1.0",
		LensTheme:         "",
		UserName:          user,
		Token:             token,
		AllowedNamespaces: allowedNamespaces,
		IsClusterAdmin:    isClusterAdmin,
		IsTenantOwner:     true,
		ChartEnable:       true,
		KubectlAccess:     true,
		DefaultNamespace:  defaultNamespace,
	}
}
