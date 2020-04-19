package tunnel

type Executioner interface {
	Run(cmd string) error
	RunWithOutput(cmd string) (string, error)
	Close() error
}
