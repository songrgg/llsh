package ssh

import (
	"bytes"
	"fmt"
	"os/exec"

	"golang.org/x/crypto/ssh"
)

// Execute the SSH task on remote host.
func (t *Task) Execute() (string, error) {
	if !t.UseSSHCommand {
		return t.execute()
	}

	return t.executeWithCommand()
}

func (t *Task) executeWithCommand() (string, error) {
	cmd := exec.Command(t.SSHCommandPath, t.Host, t.Script)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "fail to exec command: " + string(out), err
	}

	return string(out), nil
}

func (t *Task) execute() (string, error) {
	config := new(ssh.ClientConfig)
	config.SetDefaults()
	config.User = t.Username
	config.Timeout = t.Timeout
	config.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	var authMethods []ssh.AuthMethod
	for _, m := range t.AuthMethods {
		switch m.Type {
		case AuthByPassword:
			authMethods = append(authMethods, ssh.Password(m.Content))
		case AuthByPublicKey:
			pk, err := publicKey([]byte(m.Content))
			if err != nil {
				return "", err
			}
			authMethods = []ssh.AuthMethod{pk}
		default:
			return "", fmt.Errorf("invalid authentication type")
		}
	}

	config.Auth = authMethods
	if len(authMethods) < 1 {
		return "", fmt.Errorf("none available auth methods")
	}

	if t.Port < 1 {
		t.Port = 20
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", t.Host, t.Port), config)
	if err != nil {
		return "", err
	}

	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	// Once a Session is created, you can execute a single command on
	// the remote side using the Run method.
	var b bytes.Buffer
	session.Stdout = &b

	if err := session.Run(t.Script); err != nil {
		return "", fmt.Errorf("failed to run: %v", err)
	}
	return b.String(), nil
}

// publicKey reads the private key's content.
func publicKey(content []byte) (ssh.AuthMethod, error) {
	key, err := ssh.ParsePrivateKey(content)
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(key), nil
}
