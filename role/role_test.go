package role

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRole(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		name      string
		role      string
		variables map[string]string
		want      string
	}{
		"shell": {
			name: "Shell Command Generator",
			role: "Running on {{.OS}} with {{.Shell}} shell",
			variables: map[string]string{
				"OS":    "darwin",
				"Shell": "bash",
			},
			want: "You are Shell Command Generator\nRunning on darwin with bash shell",
		},
		"describeShell": {
			name:      "Shell Command Descriptor",
			role:      "Describea shell command",
			variables: map[string]string{},
			want:      "You are Shell Command Descriptor\nDescribea shell command",
		},
		"code": {
			name:      "Code Generator",
			role:      "Generating code",
			variables: map[string]string{},
			want:      "You are Code Generator\nGenerating code",
		},
		"default": {
			name: "ShellGPT",
			role: "Running on {{.OS}} with {{.Shell}} shell",
			variables: map[string]string{
				"OS":    "linux",
				"Shell": "zsh",
			},
			want: "You are ShellGPT\nRunning on linux with zsh shell",
		},
	}

	for _, tc := range tests {
		role, _ := NewRole(tc.name, tc.role, tc.variables)
		assert.Equal(t, tc.want, role.Role)
	}
}

func TestCheckGet(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		shell         bool
		describeShell bool
		code          bool
		want          string
	}{
		"shell": {
			shell:         true,
			describeShell: false,
			code:          false,
			want:          "Shell Command Generator",
		},
		"describeShell": {
			shell:         false,
			describeShell: true,
			code:          false,
			want:          "Shell Command Descriptor",
		},
		"code": {
			shell:         false,
			describeShell: false,
			code:          true,
			want:          "Code Generator",
		},
		"default": {
			shell:         false,
			describeShell: false,
			code:          false,
			want:          "ShellGPT",
		},
	}

	for _, tc := range tests {
		role, _ := CheckGet(tc.shell, tc.describeShell, tc.code)
		assert.Equal(t, tc.want, role.Name)
	}
}
