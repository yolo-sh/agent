package env

import (
	"encoding/json"
	"os"
)

// VSCodeWorkspaceConfig matches .code-workspace schema.
// See: https://code.visualstudio.com/docs/editor/multi-root-workspaces#_workspace-file-schema
type VSCodeWorkspaceConfig struct {
	Folders    []VSCodeWorkspaceConfigFolder   `json:"folders"`
	Settings   map[string]interface{}          `json:"settings"`
	Extensions VSCodeWorkspaceConfigExtensions `json:"extensions"`
}

type VSCodeWorkspaceConfigFolder struct {
	Path string `json:"path"`
}

type VSCodeWorkspaceConfigExtensions struct {
	Recommendations []string `json:"recommendations"`
}

func buildInitialVSCodeWorkspaceConfig() VSCodeWorkspaceConfig {
	return VSCodeWorkspaceConfig{
		Folders: []VSCodeWorkspaceConfigFolder{},
		Settings: map[string]interface{}{
			"remote.autoForwardPorts":      true,
			"remote.restoreForwardedPorts": true,
			// Auto-detect (using "/proc") and forward opened port.
			// Way better than "output" that parse terminal output.
			// See: https://github.com/microsoft/vscode/issues/143958#issuecomment-1050959241
			"remote.autoForwardPortsSource": "process",
			// We overwrite the $PATH environment variable in integrated terminal
			// because RVM displays warnings when VSCode changes the order of the paths.
			// See: https://github.com/microsoft/vscode/issues/70248
			"terminal.integrated.env.linux": map[string]interface{}{
				"PATH": "${env:PATH}",
			},
		},
		Extensions: VSCodeWorkspaceConfigExtensions{
			Recommendations: []string{},
		},
	}
}

func saveVSCodeWorkspaceConfigAsFile(
	vscodeWorkspaceConfigFilePath string,
	vscodeWorkspaceConfig VSCodeWorkspaceConfig,
) error {

	vscodeWorkspaceConfigAsJSON, err := json.Marshal(&vscodeWorkspaceConfig)

	if err != nil {
		return err
	}

	return os.WriteFile(
		vscodeWorkspaceConfigFilePath,
		vscodeWorkspaceConfigAsJSON,
		os.FileMode(0644),
	)
}
