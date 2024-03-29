package handler

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/alyakimenko/ssh-scanner/pkg/proxy"
	"golang.org/x/crypto/ssh"
)

type Data struct {
	Addr        string
	IsAvailable bool
	Err         error
}

func HandleSSH(wg *sync.WaitGroup, addrChan chan<- *Data, proxyAddr string, sshAddr string, config *ssh.ClientConfig) {
	defer wg.Done()

	var client *ssh.Client

	if proxyAddr != "" {
		var err error
		client, err = proxy.ProxiedSSHClient(proxyAddr, sshAddr, config)
		if err != nil {
			addrChan <- &Data{Addr: sshAddr, Err: err}
			return
		}
	} else {
		var err error
		client, err = ssh.Dial("tcp", sshAddr, config)
		if err != nil {
			addrChan <- &Data{Addr: sshAddr, Err: err}
			return
		}
	}
	defer client.Close()

	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := client.NewSession()
	if err != nil {
		addrChan <- &Data{Addr: sshAddr, Err: err}
		return
	}
	defer session.Close()

	// Once a Session is created, execute a single command on
	// the remote side using the Run method.
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("cat /etc/hostname"); err != nil {
		addrChan <- &Data{Addr: sshAddr, Err: err}
		return
	}
	fmt.Printf("[cat /etc/hostname] Hostname of %s is %s", sshAddr, b.String())
	addrChan <- &Data{Addr: sshAddr, IsAvailable: true}
}
