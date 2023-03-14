package terraform

import "encoding/json"

func GetInternalStateFromBytes(stateBytes []byte) (*terraformStateInternal, error) {
	res := &terraformStateInternal{}

	err := json.Unmarshal(stateBytes, res)

	if err != nil {
		return nil, err
	}

	return res, nil
}

type terraformStateInternal struct {
	Version          int                              `json:"version"`
	TerraformVersion string                           `json:"terraform_version"`
	Serial           int                              `json:"serial"`
	Lineage          string                           `json:"lineage"`
	Outputs          interface{}                      `json:"outputs"`
	Resources        []terraformStateInternalResource `json:"resources"`
}

type terraformStateInternalResource struct {
	Instances []terraformStateInternalInstance `json:"instances"`
	Mode      string                           `json:"mode"`
	Name      string                           `json:"name"`
	Provider  string                           `json:"provider"`
	Type      string                           `json:"type"`
}

type terraformStateInternalInstance struct {
	IndexKey      uint                   `json:"index_key,omitempty"`
	SchemaVersion uint                   `json:"schema_version,omitempty"`
	Attributes    map[string]interface{} `json:"attributes,omitempty"`
	Dependencies  []string               `json:"dependencies,omitempty"`
}
