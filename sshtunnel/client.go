package sshtunnel

import (
	"bytes"
	"io/ioutil"
	"net"

	"golang.org/x/crypto/ssh"
)

type SSHAuthConfig struct {
	Address      string `yaml:"address"`
	User         string `yaml:"user"`
	AuthPassword string `yaml:"authpassword"`
	KeyFile      string `yaml:"keyfile"`
	KeyPassword  string `yaml:"keypassword"`
}

type remoteScriptType byte
type remoteShellType byte

const (
	cmdLine remoteScriptType = iota
	rawScript

	interactiveShell remoteShellType = iota
	nonInteractiveShell
)

// Client is a wrapper over sshclient
type Client struct {
	client *ssh.Client
}

// ConnectWithPasswd starts a client connection to the given SSH server with passwd authmethod.
func ConnectWithPasswd(addr, user, passwd string) (*Client, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(passwd),
		},
		HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
	}

	return Dial("tcp", addr, config)
}

// ConnectWithKey starts a client connection to the given SSH server with key authmethod.
func ConnectWithKey(addr, user, keyfile string) (*Client, error) {
	key, err := ioutil.ReadFile(keyfile)
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
	}

	return Dial("tcp", addr, config)
}

// ConnectWithKeyPassphrase same as ConnectWithKey but with a passphrase to decrypt the private key
/*
# To generate the private key, run the following:
$ openssl genrsa -des3 -out private.pem 4096
Now you have your private key. Now you need to generate the public key.

# To generate the *public* key from your private key, run the following:
$ openssl rsa -in private.pem -outform PEM -pubout -out public.pem
Now you have a PEM format for your public key. Nice! This can’t be used with SSH’s authorized_keys file though, so we’ll have to do one more conversion:

# To generate the ssh-rsa public key format, run the following:
$ ssh-keygen -f public.pem -i -mPKCS8 > id_rsa.pub
*/
func ConnectWithKeyPassphrase(auth SSHAuthConfig) (*Client, error) {
	key, err := ioutil.ReadFile(auth.KeyFile)
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKeyWithPassphrase(key, []byte(auth.KeyPassword))
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User: auth.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
	}

	return Dial("tcp", auth.Address, config)
}

// Dial starts a client connection to the given SSH server.
// This is a wrapper on ssh.Dial
func Dial(network, addr string, config *ssh.ClientConfig) (*Client, error) {
	client, err := ssh.Dial(network, addr, config)
	if err != nil {
		return nil, err
	}
	return &Client{
		client: client,
	}, nil
}

func (c *Client) Close() error {
	return c.client.Close()
}

// Cmd create a command on client
func (c *Client) Cmd(cmd string) *RemoteCommand {
	return &RemoteCommand{
		_type:  cmdLine,
		client: c.client,
		script: bytes.NewBufferString(cmd + "\n"),
	}
}

// Script
func (c *Client) Script(script string) *RemoteCommand {
	return &RemoteCommand{
		_type:  rawScript,
		client: c.client,
		script: bytes.NewBufferString(script + "\n"),
	}
}

// ScriptFile
func (c *Client) ScriptFile(fname string) (*RemoteCommand, error) {
	content, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}
	return c.Script(string(content)), nil
}
