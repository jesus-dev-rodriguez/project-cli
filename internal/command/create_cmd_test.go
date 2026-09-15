package command

import (
	"bytes"
	"context"
	"testing"
)

func TestCreateCmd(t *testing.T) {
	service := new(fakeProjectService)
	var output bytes.Buffer
	cmd := CreateCmd(service)
	cmd.SetArgs([]string{"new-app"})
	cmd.SetOut(&output)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if got, want := output.String(), "Proyecto \"new-app\" creado y abierto.\n"; got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}

var _ interface {
	Create(context.Context, string) error
} = (*fakeProjectService)(nil)
