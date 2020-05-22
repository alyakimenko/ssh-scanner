package keys

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os/user"
)

func TryRSA(user *user.User, passphrase string) (ssh.Signer, error) {
	key, err := ioutil.ReadFile(fmt.Sprintf("%s/.ssh/id_rsa", user.HomeDir))
	if err != nil {
		return nil, err
	}

	var signer ssh.Signer
	if passphrase != "" {
		signer, err = ssh.ParsePrivateKeyWithPassphrase(key, []byte(passphrase))
		if err != nil {
			return nil, err
		}
	} else {
		signer, err = ssh.ParsePrivateKey(key)
		if err != nil {
			return nil, err
		}
	}

	return signer, nil
}
