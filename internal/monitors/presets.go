package monitors

import (
	_ "embed"
)

//go:embed presets/state_test.rego
var PresetStateTestPolicy []byte
