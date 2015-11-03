package gitolite

import "golang.org/x/crypto/ssh"

type Config struct {
	// Host keys that identify this server. At least one host key is needed.
	HostKeys []ssh.Signer

	// A callback function for handling public key authentication.
	PublicKeyCallback PublicKeyCallback

	// A callback function for handling password authentication.
	PasswordCallback PasswordCallback

	// The handler for the git protocol. The handler should implement at least
	// one function from UploadPack or ReceivePack, preferably both. Any interfaces
	// not implemented will be unavailable.
	Handler interface{}
}

func DefaultConfig() *Config {
	return &Config{
		HostKeys:          make([]ssh.Signer, 0, 1),
		PublicKeyCallback: nil,
		PasswordCallback:  nil,
		Handler:           nil,
	}
}

func (c *Config) AddHostKey(privateKey ssh.Signer) {
	c.HostKeys = append(c.HostKeys, privateKey)
}
