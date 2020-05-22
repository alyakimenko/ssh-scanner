package main

import (
	"flag"
)

type config struct {
	AddrFile      *string
	LoginFile     *string
	RSA           *bool
	RSAPassphrase *string
	Proxy         *string
}

func parseFlags() *config {
	c := &config{}

	c.AddrFile = flag.String("addr-file", "ssh-addresses.txt", "Local file with SSH addresses")
	c.LoginFile = flag.String("login-file", "login-password.txt", "Local file with login/password combinations")
	c.RSA = flag.Bool("rsa", false, "Turn on/Turn off the RSA")
	c.RSAPassphrase = flag.String("rsa-pass", "", "RSA Passphrase")
	c.Proxy = flag.String("proxy", "", "Proxy address for SSH connections")

	flag.Parse()
	return c
}
