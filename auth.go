package gitolite

import (
	"errors"

	"golang.org/x/crypto/ssh"
)

type PublicKeyCallback func(conn ssh.ConnMetadata, key ssh.PublicKey) bool

type PasswordCallback func(conn ssh.ConnMetadata, password []byte) bool

func (cb PublicKeyCallback) wrap() func(ssh.ConnMetadata, ssh.PublicKey) (*ssh.Permissions, error) {
	if cb != nil {
		return func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			return checkAuthentication(cb(conn, key))
		}
	}
	return nil
}

func (cb PasswordCallback) wrap() func(ssh.ConnMetadata, []byte) (*ssh.Permissions, error) {
	if cb != nil {
		return func(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
			return checkAuthentication(cb(conn, password))
		}
	}
	return nil
}

func checkAuthentication(ok bool) (*ssh.Permissions, error) {
	var err error
	if !ok {
		err = errors.New("Authentication failed")
	}
	return nil, err
}

func AllowAll(conn ssh.ConnMetadata, key ssh.PublicKey) bool {
	return true
}
