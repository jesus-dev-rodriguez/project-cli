package process

import (
	"context"
	"io"
	"os/exec"
)

type Runner interface {
	Run(ctx context.Context, name string, args ...string) error
}

type ExecRunner struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

func NewRunner(stdin io.Reader, stdout, stderr io.Writer) ExecRunner {
	return ExecRunner{stdin: stdin, stdout: stdout, stderr: stderr}
}

func (r ExecRunner) Run(ctx context.Context, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdin = r.stdin
	cmd.Stdout = r.stdout
	cmd.Stderr = r.stderr
	return cmd.Run()
}
