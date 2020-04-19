package sshtunnel

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/ssh"
)

type RemoteCommand struct {
	client *ssh.Client
	_type  remoteScriptType
	script *bytes.Buffer
	err    error

	stdout io.Writer
	stderr io.Writer
}

// Run
func (rs *RemoteCommand) Run() error {
	if rs.err != nil {
		fmt.Println(rs.err)
		return rs.err
	}

	if rs._type == cmdLine {
		return rs.runCmds()
	} else if rs._type == rawScript {
		return rs.runScript()
	} else {
		return errors.New("not supported RemoteScript type")
	}
}

func (rs *RemoteCommand) RunWithOutput() ([]byte, error) {
	if rs.stdout != nil {
		return nil, errors.New("stdout already set")
	}
	if rs.stderr != nil {
		return nil, errors.New("stderr already set")
	}

	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)
	rs.stdout = &stdout
	rs.stderr = &stderr
	err := rs.Run()
	if err != nil {
		return stderr.Bytes(), err
	}
	return stdout.Bytes(), err
}

func (rs *RemoteCommand) SetCmd(cmd string) *RemoteCommand {
	_, err := rs.script.WriteString(cmd + "\n")
	if err != nil {
		rs.err = err
	}
	return rs
}

func (rs *RemoteCommand) SetStdio(stdout, stderr io.Writer) *RemoteCommand {
	rs.stdout = stdout
	rs.stderr = stderr
	return rs
}

func (rs *RemoteCommand) runCmd(cmd string) error {
	session, err := rs.client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	session.Stdout = rs.stdout
	session.Stderr = rs.stderr

	if err := session.Run(cmd); err != nil {
		return err
	}
	return nil
}

func (rs *RemoteCommand) runCmds() error {
	for {
		statment, err := rs.script.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if err := rs.runCmd(statment); err != nil {
			return err
		}
	}

	return nil
}

func (rs *RemoteCommand) runScript() error {
	session, err := rs.client.NewSession()
	if err != nil {
		return err
	}

	session.Stdin = rs.script
	session.Stdout = rs.stdout
	session.Stderr = rs.stderr

	if err := session.Shell(); err != nil {
		return err
	}
	if err := session.Wait(); err != nil {
		return err
	}

	return nil
}
