package ssh

import "time"

type AuthType int

const (
	// AuthByPassword indicates using password to authenticate.
	AuthByPassword AuthType = iota + 1
	// AuthByPublicKey indicates using public key to authenticate.
	AuthByPublicKey
)

// Task wraps the remote SSH task.
type Task struct {
	// Username is the username used to login the remote server.
	Username string `json:"username"`
	// AuthMethods includes the authentication methods used to login server.
	AuthMethods []AuthMethod `json:"auth_methods"`
	// Host indicates the remote server host.
	Host string `json:"host"`
	// Port indicates the remote server ssh port, default is 22.
	Port int `json:"port"`
	// Script is the shell command.
	Script string `json:"script"`
	// Timeout indicates the max command running time.
	Timeout time.Duration `json:"timeout"`
	// UseSSHCommand uses shell command to execute commands on remote machines.
	UseSSHCommand bool `json:"use_ssh_command"`
	// SSHCommandPath indicates the file location of the SSH command.
	SSHCommandPath string `json:"ssh_command_path"`
}

// AuthMethod is the authentication method.
type AuthMethod struct {
	Type    AuthType `json:"type"`
	Content string   `json:"content"`
}
