package yagu

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/kevinburke/ssh_config"
)

type SSHMachine struct {
	User string
	Keys *ssh.PublicKeys
}

func SSHCredentials(machine string) (SSHMachine, error) {
	fmt.Println("ssh.CredsLookup")
	pub := ""
	usr := "git"

	// first look for a usr override, can be used for key var or default location
	if u := os.Getenv("HOF_SSHUSR"); u != "" {
		usr = u
	}

	// look for env var key location
	if key := os.Getenv("HOF_SSHKEY"); key != "" {
		fmt.Println("ssh.SSHEnvVar")
		pks, err := ssh.NewPublicKeysFromFile(usr, key, "")
		if err != nil {
			return SSHMachine{}, err
		}
		return SSHMachine{usr, pks}, nil
	}

	// try to get homedir
	hdir, err := os.UserHomeDir()
	if err != nil {
		// no home dir?
		return SSHMachine{}, err
	}

	// try sshconfig
	_, uerr := os.Lstat(filepath.Join(hdir, ".ssh", "config"))
	_, serr := os.Lstat(filepath.Join("etc", "ssh", "ssh_config"))
	if uerr == nil || serr == nil {
		fmt.Println("ssh.SSHConfig")
		return getSSHConfigVals(machine)
	}

	// fallback on default pubkey
	fmt.Println("ssh.DefaultPubkey")
	pub = filepath.Join(hdir, ".ssh", "id_rsa")
	pks, err := ssh.NewPublicKeysFromFile(usr, pub, "")
	fmt.Println("err:", err)
	if err != nil {
		return SSHMachine{}, err
	}

	return SSHMachine{usr, pks}, nil
}

func getSSHConfigVals(machine string) (SSHMachine, error) {
	// try to lookup the machine in config
	pub, err := ssh_config.GetStrict(machine, "IdentityFile")
	if err != nil {
		fmt.Println(pub, err)
		return SSHMachine{}, err
	}

	// replace if key location has ~
	if strings.HasPrefix(pub, "~") {
		// we already validated homedir from calling function
		hdir, _ := os.UserHomeDir()
		pub = strings.Replace(pub, "~", hdir, 1)
	}

	// override user if defined in config
	usr := ssh_config.Get(machine, "User")
	if usr == "" {
		usr = "git"
	}

	// get key from filename
	pks, err := ssh.NewPublicKeysFromFile(usr, pub, "")
	if err != nil {
		return SSHMachine{}, err
	}

	return SSHMachine{usr, pks}, nil
}
