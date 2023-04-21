# Remote SSH Command Execution in Go
This is a small Go program that allows you to execute a command on a remote host, using an SSH connection. It reads your SSH configuration file and connects to the specified host using the provided alias.

## Prerequisites
Before you can use the program, you must have the following installed on your system:
- Go (version 1.16 or later)

## Setup
You can clone the full easy-going repo to have access to all examples. If you want to use only this example, you can copy the code but ensure to have the required packages installed:

- `go get -u golang.org/x/crypto/ssh`

- `go get -u github.com/kevinburke/ssh_config`


### Build the program:
`go build -o ssh_exec`

### Run the program with the host alias and command you want to execute:

`./ssh_exec <host_alias> "<command>"`

Replace <host_alias> with the alias of the host in your SSH config file, and <command> with the command you want to execute on the remote host. For example:

`./ssh_exec my_server "uname -a"`

## Notes
This program uses the `InsecureIgnoreHostKey()` function for the HostKeyCallback configuration, which skips host key verification. For production use, you should replace it with a proper host key verification method, such as `knownhosts.New()` from the `golang.org/x/crypto/ssh/knownhosts` package.

The program assumes your private key is located at `$HOME/.ssh/id_rsa`. If your key is in a different location or has a different name, update the keyPath variable accordingly.

If your private key is password-protected, set the password in the `ssh.ParsePrivateKeyWithPassphrase()` function. If your private key does not require a password, use the `ssh.ParsePrivateKey()` function instead. Remember to remove any hardcoded passwords before committing your code.