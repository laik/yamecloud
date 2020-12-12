package gateway

import "encoding/json"

type UserConfig struct {
	LensVersion       string   `json:"lensVersion"`
	LensTheme         string   `json:"lensTheme"`
	UserName          string   `json:"userName"`
	Token             string   `json:"token"`
	AllowedNamespaces []string `json:"allowedNamespaces"`
	IsClusterAdmin    bool     `json:"isClusterAdmin"`
	IsTenantOwner     bool     `json:"isTenantOwner"`
	ChartEnable       bool     `json:"chartEnable"`
	KubectlAccess     bool     `json:"kubectlAccess"`
	DefaultNamespace  string   `json:"defaultNamespace"`
}

func (uc *UserConfig) String() string {
	bytesData, _ := json.Marshal(uc)
	return string(bytesData)
}

func NewUserConfig(user string, token string, allowedNamespaces []string, defaultNamespace string) *UserConfig {
	isClusterAdmin := false
	if user == "admin" {
		isClusterAdmin = true
	}
	return &UserConfig{
		LensVersion:       "1.0",
		LensTheme:         "",
		UserName:          user,
		Token:             token,
		AllowedNamespaces: allowedNamespaces,
		IsClusterAdmin:    isClusterAdmin,
		ChartEnable:       true,
		KubectlAccess:     true,
		DefaultNamespace:  defaultNamespace,
	}
}
