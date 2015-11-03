package gitolite

import "golang.org/x/crypto/ssh"

type Config struct {
	HostKeys []ssh.Signer

	PublicKeyCallback PublicKeyCallback

	PasswordCallback PasswordCallback
}

func DefaultConfig() *Config {
	return &Config{
		HostKeys:          make([]ssh.Signer, 0, 1),
		PublicKeyCallback: nil,
		PasswordCallback:  nil,
	}
}

func (c *Config) AddHostKey(privateKey ssh.Signer) {
	c.HostKeys = append(c.HostKeys, privateKey)
}
