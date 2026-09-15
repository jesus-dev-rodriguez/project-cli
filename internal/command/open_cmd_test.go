package command

import (
	"context"
	"testing"

	"github.com/jesus-dev-rodriguez/project-cli/internal/services"
	"github.com/spf13/cobra"
)

type fakeProjectService struct {
	opened string
}

func (s *fakeProjectService) List() ([]string, error) {
	return []string{"alpha", "zeta"}, nil
}

func (s *fakeProjectService) Open(_ context.Context, name string) error {
	s.opened = name
	return nil
}

func (s *fakeProjectService) Create(_ context.Context, _ string) error {
	return nil
}

func (s *fakeProjectService) Path(_ string) (string, error) {
	return "/tmp/project", nil
}

var _ services.ProjectService = (*fakeProjectService)(nil)

func TestOpenCmd(t *testing.T) {
	service := new(fakeProjectService)
	cmd := OpenCmd(service)
	cmd.SetArgs([]string{"my-app"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if service.opened != "my-app" {
		t.Fatalf("opened project = %q, want %q", service.opened, "my-app")
	}
}

func TestOpenCmdCompletion(t *testing.T) {
	service := new(fakeProjectService)
	cmd := OpenCmd(service)
	got, directive := cmd.ValidArgsFunction(cmd, nil, "")

	if directive != cobra.ShellCompDirectiveNoFileComp {
		t.Fatalf("completion directive = %v, want no-file-comp", directive)
	}
	if len(got) != 2 || got[0] != "alpha" || got[1] != "zeta" {
		t.Fatalf("completion = %v, want [alpha zeta]", got)
	}
}
