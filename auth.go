package gitolite

import "golang.org/x/crypto/ssh"

// Temporary function unti we implement a way to add custom authentication methods.
func noAuthentication(ssh.ConnMetadata, ssh.PublicKey) (*ssh.Permissions, error) {
	return nil, nil
}
