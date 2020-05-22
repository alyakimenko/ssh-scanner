package main

import (
	"log"
	"os/user"
	"sync"
	"time"

	"gitlab.com/astroproxy/ssh-scanner/pkg/handler"
	"gitlab.com/astroproxy/ssh-scanner/pkg/interactive"
	"gitlab.com/astroproxy/ssh-scanner/pkg/keys"
	"gitlab.com/astroproxy/ssh-scanner/pkg/parser"
	"golang.org/x/crypto/ssh"
)

func main() {
	flags := parseFlags()

	currentUser, err := user.Current()
	if err != nil {
		log.Fatal("Error occurred during getting current user: ", err)
	}

	rsa, err := keys.TryRSA(currentUser, *flags.RSAPassphrase)
	if err != nil {
		log.Fatal("Error occurred during parsing RSA: ", err)
	}

	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.Password("root"),
			ssh.PublicKeys(rsa),
			ssh.KeyboardInteractive(interactive.SSHInteractive),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout: 2000*time.Millisecond,
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
