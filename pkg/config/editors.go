package config

import "fmt"

// EditorTemplate defines an editor with its command template
type EditorTemplate struct {
	Name    string
	Command string // Template with %file%, %line%, %column%
}

// BuiltInEditors contains predefined editor configurations
var BuiltInEditors = map[string]EditorTemplate{
	"vscode": {
		Name:    "Visual Studio Code",
		Command: "code --goto %file%:%line%:%column%",
	},
	"vim": {
		Name:    "Vim",
		Command: "vim +%line% %file%",
	},
}

// GetEditorTemplate returns the command template for a built-in editor
func GetEditorTemplate(editorName string) (string, error) {
	if template, ok := BuiltInEditors[editorName]; ok {
		return template.Command, nil
	}
	return "", fmt.Errorf("unknown editor: %s", editorName)
}
