package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
		assert.Equal(t, tc.want, role.name)
	}
}
