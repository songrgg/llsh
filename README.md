# Light-weight parallel remote ssh tool
LLSH is for running shell commands on remote machines through SSH, it's quite common that sometimes
operators need to run some commands on multiple servers.

LLSH is a super simple and direct tool with only command line arguments, no need to use conf file.
You can simply run the SSH commands through two different ways:
* create a process running `/bin/ssh` command to run the remote commands
* use SSH protocol to access remote servers

The first approach works fine when all of the settings are already ready in local SSH, but it has
extra costs to create a process, the second one means faster and less resource and supports password
and public key to authenticate.

## Usages

### Use ssh command to access host1 and host2 running ls.
```bash
$ llsh --hosts hosts1,hosts2 \
 --user sojiang \
 -c "ls"
```

### Use host file
Hhost file is INI-like file, the format is like this, file name can be `test.hosts`
```ini
host1
host2

[app]
app1
app2

[web]
web1
web2
``` 

If you want to run `ls` on all of the servers, 
```bash
$ llsh --host_file test.hosts \
 --user sojiang \
 -c "ls"
```

If you only want to manipulate app section, it will choose only `app1` and `app2`
```bash
$ llsh --host_file test.hosts \
 --section app \
 --user sojiang \
 -c "ls"
```

### Help
```bash
$ llsh --help
llsh executes the remote shell command on multiple remote servers

Usage:
  llsh [flags]

Flags:
  -c, --command string            The command you need to run on the servers
  -h, --help                      help for llsh
      --host_file string          The path of host file, the content should be hosts separated by new line, if specified, hosts will be omitted
      --hosts string              The host names for remote server, split by comma, for example, host1,host2,host3
  -p, --password                  If we need to input password from terminal
  -k, --publickey string          The file path of the public key, like ~/.ssh/id_rsa.pub
      --section string            The section of the host_file, default to all hosts when not specified or empty
      --ssh_command_path string   The file path of the ssh command, like /bin/ssh (default "ssh")
      --use_ssh_command           Use the ssh executable, like /bin/ssh (default true)
  -u, --user string               The username for login user
```
