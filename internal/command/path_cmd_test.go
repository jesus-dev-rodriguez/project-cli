package command

import (
	"bytes"
	"testing"
)

func TestPathCmd(t *testing.T) {
	var output bytes.Buffer
	cmd := PathCmd(new(fakeProjectService))
	cmd.SetArgs([]string{"my-app"})
	cmd.SetOut(&output)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if got, want := output.String(), "/tmp/project\n"; got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}
