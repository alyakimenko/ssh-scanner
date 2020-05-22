package main

import (
	"log"
	"os/user"
	"sync"

	"gitlab.com/astroproxy/ssh-scanner/pkg/handler"
	"gitlab.com/astroproxy/ssh-scanner/pkg/interactive"
	"gitlab.com/astroproxy/ssh-scanner/pkg/keys"
	"gitlab.com/astroproxy/ssh-scanner/pkg/parser"
	"golang.org/x/crypto/ssh"
)

func main() {
	flags := parseFlags()

	var rsaSigner ssh.Signer
	if *flags.RSA {
		currentUser, err := user.Current()
		if err != nil {
			log.Fatal("Error occurred during getting current user: ", err)
		}

		rsaSigner, err = keys.TryRSA(currentUser, *flags.RSAPassphrase)
		if err != nil {
			log.Fatal("Error occurred during parsing RSA: ", err)
		}
	}

	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.Password("root"),
			ssh.KeyboardInteractive(interactive.SSHInteractive),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	if rsaSigner != nil {
		config.Auth = append(config.Auth, ssh.PublicKeys(rsaSigner))
	}

	addresses, err := parser.ParseAddrFile(*flags.AddrFile)
	if err != nil {
		log.Fatal("Error occurred during parsing addr file: ", err)
	}

	var wg sync.WaitGroup
	addrChan := make(chan *handler.Data)

	go func() {
		wg.Wait()
		close(addrChan)
	}()

	for _, addr := range addresses {
		wg.Add(1)
		go handler.HandleSSH(&wg, addrChan, *flags.Proxy, addr, config)
	}

	var availableAddr []string
	for data := range addrChan {
		if data.Err != nil {
			log.Println(data.Addr + " is not available")
		}
		if data.IsAvailable {
			availableAddr = append(availableAddr, data.Addr)
		}
	}
	log.Println("Available addresses: ", availableAddr)
}
