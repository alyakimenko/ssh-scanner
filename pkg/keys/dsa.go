package keys

import (
	"fmt"
	"io/ioutil"
	"os/user"

	"golang.org/x/crypto/ssh"
)

func TryDSA(user *user.User) (ssh.Signer, error) {
	key, err := ioutil.ReadFile(fmt.Sprintf("%s/.ssh/id_dsa", user.HomeDir))
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}

	return signer, nil
}
