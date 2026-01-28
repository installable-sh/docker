package install

import (
	"bytes"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want *Install
	}{
		{
			name: "empty args",
			args: []string{},
			want: &Install{},
		},
		{
			name: "help flag",
			args: []string{"--help"},
			want: &Install{ShowHelp: true},
		},
		{
			name: "help short flag",
			args: []string{"-h"},
			want: &Install{ShowHelp: true},
		},
		{
			name: "version flag",
			args: []string{"--version"},
			want: &Install{ShowVersion: true},
		},
		{
			name: "version short flag",
			args: []string{"-v"},
			want: &Install{ShowVersion: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args)

			if got.ShowHelp != tt.want.ShowHelp {
				t.Errorf("ShowHelp = %v, want %v", got.ShowHelp, tt.want.ShowHelp)
			}
			if got.ShowVersion != tt.want.ShowVersion {
				t.Errorf("ShowVersion = %v, want %v", got.ShowVersion, tt.want.ShowVersion)
			}
		})
	}
}

func TestExec_Help(t *testing.T) {
	var stdout bytes.Buffer
	cmd := &Install{
		ShowHelp: true,
		Stdout:   &stdout,
		Stderr:   &bytes.Buffer{},
	}

	err := cmd.Exec()
	if err != nil {
		t.Errorf("Exec() error = %v", err)
	}

	if !strings.Contains(stdout.String(), "usage:") {
		t.Errorf("Exec() stdout = %q, want usage info", stdout.String())
	}
}

func TestExec_Version(t *testing.T) {
	var stdout bytes.Buffer
	cmd := &Install{
		ShowVersion: true,
		Stdout:      &stdout,
		Stderr:      &bytes.Buffer{},
	}

	err := cmd.Exec()
	if err != nil {
		t.Errorf("Exec() error = %v", err)
	}
}

func TestExec_NotImplemented(t *testing.T) {
	var stderr bytes.Buffer
	cmd := &Install{
		Stdout: &bytes.Buffer{},
		Stderr: &stderr,
	}

	err := cmd.Exec()
	if err == nil {
		t.Error("Exec() expected error for not implemented")
	}

	if !strings.Contains(stderr.String(), "coming soon") {
		t.Errorf("Exec() stderr = %q, want 'coming soon'", stderr.String())
	}
}
