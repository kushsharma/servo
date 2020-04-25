package tunnel

import (
	"strings"

	"github.com/kushsharma/servo/internal"
	sshtunnel "github.com/kushsharma/servo/tunnel/ssh"
	"github.com/pkg/errors"
)

type SSHTunnel struct {
	client *sshtunnel.Client
}

func (tnl *SSHTunnel) Run(cmd string) error {
	execCmd := tnl.client.Cmd(cmd)
	return execCmd.Run()
}

func (tnl *SSHTunnel) RunWithOutput(cmd string) (string, error) {
	execCmd := tnl.client.Cmd(cmd)
	output, err := execCmd.RunWithOutput()
	return strings.TrimSpace(string(output)), err
}

func (tnl *SSHTunnel) Close() error {
	return tnl.client.Close()
}

func NewSSHTunnel(authConfig internal.SSHAuthConfig) (*SSHTunnel, error) {
	tnl := new(SSHTunnel)

	if authConfig.Host == "" || authConfig.User == "" {
		return nil, errors.New("invalid machine auth config")
	}
	sshclient, err := sshtunnel.ConnectWithKeyPassphrase(authConfig)
	if err != nil {
		return nil, err
	}
	tnl.client = sshclient

	return tnl, nil
}
