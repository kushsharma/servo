package tunnel

import (
	"github.com/kushsharma/servo/internal"
	"github.com/kushsharma/servo/sshtunnel"
)

type SSHTunnel struct {
	client *sshtunnel.Client
}

func (tnl *SSHTunnel) Run(cmd string) error {

	return nil
}

func (tnl *SSHTunnel) RunWithOutput(cmd string) (string, error) {
	return "", nil
}

func (tnl *SSHTunnel) Close() error {
	return tnl.client.Close()
}

func NewSSHTunnel(authConfig internal.SSHAuthConfig) (*SSHTunnel, error) {
	tnl := new(SSHTunnel)
	sshclient, err := sshtunnel.ConnectWithKeyPassphrase(authConfig)
	if err != nil {
		return nil, err
	}
	tnl.client = sshclient
	return tnl, nil
}
