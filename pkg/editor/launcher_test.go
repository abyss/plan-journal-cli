package editor

import (
	"reflect"
	"testing"
)

func TestParseCommand(t *testing.T) {
	tests := []struct {
		name   string
		cmdStr string
		want   []string
	}{
		{
			name:   "simple command",
			cmdStr: "vim file.txt",
			want:   []string{"vim", "file.txt"},
		},
		{
			name:   "command with flags",
			cmdStr: "code --goto file.txt:10:5",
			want:   []string{"code", "--goto", "file.txt:10:5"},
		},
		{
			name:   "command with double quotes",
			cmdStr: `code --goto "file with spaces.txt:10:5"`,
			want:   []string{"code", "--goto", "file with spaces.txt:10:5"},
		},
		{
			name:   "command with single quotes",
			cmdStr: `vim +'call cursor(10, 5)' file.txt`,
			want:   []string{"vim", "+call cursor(10, 5)", "file.txt"},
		},
		{
			name:   "command with multiple spaces",
			cmdStr: "vim   +10    file.txt",
			want:   []string{"vim", "+10", "file.txt"},
		},
		{
			name:   "complex nvim command",
			cmdStr: `nvim -c "call cursor(10, 5)" file.txt`,
			want:   []string{"nvim", "-c", "call cursor(10, 5)", "file.txt"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseCommand(tt.cmdStr)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}
