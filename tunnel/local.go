package tunnel

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
	"time"
)

const (
	CommandTimeout = 15 * time.Minute
)

type LocalTunnel struct {
	stderr bytes.Buffer
	stdout bytes.Buffer
}

func (tnl *LocalTunnel) Run(command string) error {
	ctx, cancel := context.WithTimeout(context.Background(), CommandTimeout)
	defer cancel()

	parts := strings.Split(command, " ")
	if len(parts) == 1 {
		//make sure parts has min len of 2
		parts = append(parts, "")
	}
	return exec.CommandContext(ctx, parts[0], parts[:1]...).Run()
}

func (tnl *LocalTunnel) RunWithOutput(command string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CommandTimeout)
	defer cancel()

	parts := strings.Split(command, " ")
	if len(parts) == 1 {
		//make sure parts has min len of 2
		parts = append(parts, "")
	}
	execCmd := exec.CommandContext(ctx, parts[0], parts[1:]...)
	execCmd.Stdout = &tnl.stdout
	execCmd.Stderr = &tnl.stderr
	err := execCmd.Run()

	if err != nil {
		errStr := strings.TrimSpace(tnl.stderr.String())
		tnl.stderr.Reset()
		return errStr, err
	}
	output := strings.TrimSpace(tnl.stdout.String())
	tnl.stdout.Reset()
	return output, err
}

func (tnl *LocalTunnel) Close() error {
	return nil
}

func NewLocalTunnel() *LocalTunnel {
	tnl := new(LocalTunnel)
	return tnl
}
