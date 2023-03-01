package monitors

import (
	_ "embed"
)

//go:embed presets/drift_detection.rego
var PresetDriftDetectionPolicy []byte
