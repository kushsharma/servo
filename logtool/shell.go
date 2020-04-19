package logtool

import (
	"fmt"
	"strings"

	"github.com/kushsharma/servo/tunnel"
)

// ShellService that implements logtool.Service interface
type ShellService struct {
	tnl tunnel.Executioner
}

// Fetch extracts the contents of file
func (svc *ShellService) Fetch(path, filename string) (string, error) {
	cmdLine := fmt.Sprintf(`cat %s`, path)
	output, err := svc.tnl.RunWithOutput(cmdLine)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

//List return file names in the directory
func (svc *ShellService) List(path string) ([]string, error) {
	cmdLine := fmt.Sprintf(`find %s -type f`, path)
	output, err := svc.tnl.RunWithOutput(cmdLine)
	if err != nil {
		return nil, err
	}

	files := strings.Split(string(output), "\n")
	return files, nil
}

// Delete remove a single file in the provided absolute path
func (svc *ShellService) Delete(path string) error {
	cmdLine := fmt.Sprintf(`rm %s`, path)
	err := svc.tnl.Run(cmdLine)
	if err != nil {
		return err
	}

	return nil
}

// Clean removes all the files older than provided days in given directory
func (svc *ShellService) Clean(path string, daysold int) error {
	cmdLine := fmt.Sprintf(`find %s -type f -mtime +%d -delete`, path, daysold)
	err := svc.tnl.Run(cmdLine)
	if err != nil {
		return err
	}

	return nil
}

// DryClean only list files that can be removed instead of actually removing them
func (svc *ShellService) DryClean(path string, daysold int) ([]string, error) {
	cmdLine := fmt.Sprintf(`find %s -type f -mtime +%d -print`, path, daysold)
	output, err := svc.tnl.RunWithOutput(cmdLine)
	if err != nil {
		return []string{output}, err
	}

	return strings.Split(string(output), "\n"), nil
}

// NewService returns a instance of ShellService that implements LogMangager over shell
func NewService(tnl tunnel.Executioner) *ShellService {
	return &ShellService{
		tnl: tnl,
	}
}
