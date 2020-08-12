package service

import ()

type LastConfigReqParam struct {
	Detail *Detail `json:"info"`
}

type Detail struct {
	AddonID  string `json:"addon_id"`
	TenantID string `json:"tenant_id"`
	NetType  string `json:"net"`
	Product  string `json:"product"`
}

type LastConfigRespParam struct {
	Result *LastConfigResult `json:"result"`
}

type LastConfigResult struct {
	SDKURL         string     `json:"sdk_url"`
	SDKConfig      *SDKConfig `json:"sdk_config"`
	ServerPortPath string     `json:"server_port_path"`
}

type SDKConfig struct {
	// TODO
}
