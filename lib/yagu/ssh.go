package yagu

import (
	"os"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/kevinburke/ssh_config"
)

type SSHMachine struct {
	User string
	Keys *ssh.PublicKeys
}

func SSHCredentials(machine string) (SSHMachine, error) {
	pk, err := ssh_config.GetStrict(machine, "IdentityFile")
	if err != nil {
		return SSHMachine{}, err
	}
	if strings.HasPrefix(pk, "~") {
		if hdir, err := os.UserHomeDir(); err == nil {
			pk = strings.Replace(pk, "~", hdir, 1)
		}
	}
	usr := ssh_config.Get(machine, "User")
	if usr == "" {
		usr = "git"
	}

	pks, err := ssh.NewPublicKeysFromFile(usr, pk, "")
	if err != nil {
		return SSHMachine{}, err
	}

	return SSHMachine{usr, pks}, nil
}
