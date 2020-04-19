package logtool

import (
	"fmt"
	"strings"

	"github.com/kushsharma/servo/sshtunnel"
)

// Service that implements logtool.Service interface
type ShellService struct {
	ssh *sshtunnel.Client
}

// Fetch extracts the contents of file
func (svc *ShellService) Fetch(path, filename string) (string, error) {
	cmdLine := fmt.Sprintf(`cat "%s"`, path)
	shellCmd := svc.ssh.Cmd(cmdLine)
	output, err := shellCmd.RunWithOutput()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

//List return file names in the directory
func (svc *ShellService) List(path string) ([]string, error) {
	cmdLine := fmt.Sprintf(`find "%s" -type f`, path)
	shellCmd := svc.ssh.Cmd(cmdLine)
	output, err := shellCmd.RunWithOutput()
	if err != nil {
		return nil, err
	}

	files := strings.Split(string(output), "\n")
	return files, nil
}

// Delete remove a single file in the provided absolute path
func (svc *ShellService) Delete(path string) error {
	cmdLine := fmt.Sprintf(`rm "%s"`, path)
	shellCmd := svc.ssh.Cmd(cmdLine)
	_, err := shellCmd.RunWithOutput()
	if err != nil {
		return err
	}

	return nil
}

// Clean removes all the files older than provided days in given directory
func (svc *ShellService) Clean(path string, daysold int) error {
	cmdLine := fmt.Sprintf(`find "%s" -type f -mtime +%d -delete;`, path, daysold)
	shellCmd := svc.ssh.Cmd(cmdLine)
	_, err := shellCmd.RunWithOutput()
	if err != nil {
		return err
	}

	return nil
}

// DryClean only list files that can be removed instead of actually removing them
func (svc *ShellService) DryClean(path string, daysold int) ([]string, error) {
	cmdLine := fmt.Sprintf(`find "%s" -type f -mtime +%d -print;`, path, daysold)
	shellCmd := svc.ssh.Cmd(cmdLine)
	output, err := shellCmd.RunWithOutput()
	if err != nil {
		return nil, err
	}

	return strings.Split(string(output), "\n"), nil
}

// NewService returns a instance of ShellService that implements LogMangager over shell
func NewService(client *sshtunnel.Client) *ShellService {
	return &ShellService{
		ssh: client,
	}
}
