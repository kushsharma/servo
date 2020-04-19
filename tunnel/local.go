package tunnel

type LocalTunnel struct {
}

func (tnl *LocalTunnel) Run(cmd string) error {

	return nil
}

func (tnl *LocalTunnel) RunWithOutput(cmd string) (string, error) {

	return "", nil
}

func (tnl *LocalTunnel) Close() error {
	return nil
}

func NewLocalTunnel() *LocalTunnel {
	tnl := new(LocalTunnel)
	return tnl
}
