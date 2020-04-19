package tunnel

import "github.com/kushsharma/servo/sshtunnel"

type SSHTunnel struct {
	client *sshtunnel.Client
}

func (tnl *SSHTunnel) Run(cmd string) error {

	return nil
}

func (tnl *SSHTunnel) RunWithOutput(cmd string) (string, error) {

	return "", nil
}

func NewSSHTunnel(c *sshtunnel.Client) *SSHTunnel {
	tnl := new(SSHTunnel)
	tnl.client = c

	return tnl
}
