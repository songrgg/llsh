package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"gopkg.in/ini.v1"

	"github.com/fatih/color"
	"github.com/songrgg/llsh/pkg/helper"
	"github.com/songrgg/llsh/pkg/ssh"
	"github.com/spf13/cobra"
)

type options struct {
	Command        string
	Hosts          string
	HostFile       string
	Section        string
	User           string
	UseSSHCommand  bool
	UsePassword    bool
	PublicKeyPath  string
	SSHCommandPath string
}

type sshResult struct {
	Host   string
	Result string
	Error  error
}

// NewLLSHCommand creates `llsh` command.
func NewLLSHCommand() *cobra.Command {
	options := options{}

	cmds := &cobra.Command{
		Use:   "llsh",
		Short: "llsh executes the remote shell command on multiple remote servers",
		Run: func(cmd *cobra.Command, args []string) {
			options.validate()
			options.run()
		},
	}

	cmds.PersistentFlags().StringVar(&options.Hosts, "hosts", "",
		"The host names for remote server, split by comma, for example, host1,host2,host3")
	cmds.PersistentFlags().StringVar(&options.HostFile, "host_file", "",
		"The path of host file, the content should be hosts separated by new line, if specified, hosts will be omitted")
	cmds.PersistentFlags().StringVar(&options.Section, "section", "",
		"The section of the host_file, default to all hosts when not specified or empty")
	cmds.PersistentFlags().StringVarP(&options.User, "user", "u", "",
		"The username for login user")
	cmds.PersistentFlags().StringVarP(&options.Command, "command", "c", "",
		"The command you need to run on the servers")
	cmds.PersistentFlags().BoolVar(&options.UseSSHCommand, "use_ssh_command", true,
		"Use the ssh executable, like /bin/ssh")
	cmds.PersistentFlags().BoolVarP(&options.UsePassword, "password", "p", false,
		"If we need to input password from terminal")
	cmds.PersistentFlags().StringVarP(&options.PublicKeyPath, "publickey", "k", "",
		"The file path of the public key, like ~/.ssh/id_rsa.pub")
	cmds.PersistentFlags().StringVar(&options.SSHCommandPath, "ssh_command_path", "ssh",
		"The file path of the ssh command, like /bin/ssh")
	return cmds
}

func (o *options) authMethods() []ssh.AuthMethod {
	authMethods := make([]ssh.AuthMethod, 0)
	if o.UsePassword && !o.UseSSHCommand {
		password, err := helper.ReadPassword("Password: ")
		if err != nil {
			helper.Errorlq("password can't be read")
		}

		authMethods = append(authMethods, ssh.AuthMethod{
			Type:    ssh.AuthByPassword,
			Content: password,
		})
	}

	if o.PublicKeyPath != "" && !o.UseSSHCommand {
		content, err := ioutil.ReadFile(o.PublicKeyPath)
		if err != nil {
			helper.Errorlq("couldn't read the content of public key: ", err)
		}

		authMethods = append(authMethods, ssh.AuthMethod{
			Type:    ssh.AuthByPublicKey,
			Content: string(content),
		})
	}
	return authMethods
}

func (o *options) hosts() ([]string, error) {
	if o.HostFile != "" {
		//cfg, err := ini.Load(o.HostFile)
		cfg, err := ini.LoadSources(ini.LoadOptions{SkipUnrecognizableLines: true}, o.HostFile)
		if err != nil {
			return nil, err
		}

		cfg, err = ini.LoadSources(ini.LoadOptions{
			UnparseableSections: cfg.SectionStrings(),
		}, o.HostFile)
		if err != nil {
			return nil, err
		}

		// Set to all roles which is the ini `DEFAULT` section.
		var hosts []string
		var sections []*ini.Section
		if o.Section == "" {
			sections = cfg.Sections()
		} else {
			sec, err := cfg.GetSection(o.Section)
			if err != nil {
				return nil, err
			}
			sections = append(sections, sec)
		}

		for _, sec := range sections {
			for _, host := range strings.Split(string(sec.Body()), "\n") {
				if h := strings.Trim(host, " "); h != "" {
					hosts = append(hosts, h)
				}
			}
		}
		return hosts, nil
	}
	return strings.Split(o.Hosts, ","), nil
}

func (o *options) run() {
	authMethods := o.authMethods()

	hosts, err := o.hosts()
	if err != nil {
		color.Red("Invalid hosts: %v", err)
		os.Exit(1)
	}

	var wg sync.WaitGroup
	var results []sshResult
	var lock = sync.Mutex{}
	for _, host := range hosts {
		wg.Add(1)
		go func(host string) {
			t := &ssh.Task{
				Username:       o.User,
				AuthMethods:    authMethods,
				Host:           host,
				Script:         o.Command,
				UseSSHCommand:  o.UseSSHCommand,
				SSHCommandPath: o.SSHCommandPath,
			}

			result, err := t.Execute()
			if err != nil {
				result = fmt.Sprintf("fail to execute task: %s", result)
			}

			lock.Lock()
			results = append(results, sshResult{
				Host:   host,
				Result: result,
				Error:  err,
			})
			lock.Unlock()
			wg.Done()
		}(host)
	}
	wg.Wait()

	for _, result := range results {
		if result.Error != nil {
			color.Red("Host %s\n", result.Host)
		} else {
			color.Green("Host %s\n", result.Host)
		}

		_, _ = fmt.Println(result.Result)
	}
}

func (o *options) validate() {
	if o.HostFile == "" && o.Hosts == "" {
		helper.Errorlq("hosts can't be empty")
	}

	if o.PublicKeyPath != "" && o.UsePassword {
		helper.Errorlq("can't use both public key and password")
	}

	if (o.PublicKeyPath != "" || o.UsePassword) && o.UseSSHCommand {
		color.Yellow("Public key or password option will be ignored when using ssh command")
	}
}
