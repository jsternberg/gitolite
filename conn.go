package gitolite

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os/user"
	"path/filepath"

	"golang.org/x/crypto/ssh"
)

type sshConn struct {
	conn       *ssh.ServerConn
	newChannel <-chan ssh.NewChannel
	requests   <-chan *ssh.Request
}

func (s *Server) serve(conn net.Conn) error {
	config := ssh.ServerConfig{
		PublicKeyCallback: noAuthentication,
	}

	user, err := user.Current()
	if err != nil {
		return err
	}
	privateKeyPath := filepath.Join(user.HomeDir, ".ssh/id_rsa")

	pemBytes, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return fmt.Errorf("unable to load server private key: %s", err)
	}

	privateKey, err := ssh.ParsePrivateKey(pemBytes)
	if err != nil {
		return err
	}
	config.AddHostKey(privateKey)

	serverConn, channelRequestCh, globalRequestCh, err := ssh.NewServerConn(conn, &config)
	if err != nil {
		return err
	}

	// connection succeeded at this point. create the ssh connection and start the go procs.
	newConn := sshConn{
		conn:       serverConn,
		newChannel: channelRequestCh,
		requests:   globalRequestCh,
	}
	go newConn.processChannelRequests()
	go newConn.processGlobalRequests()
	return nil
}

func (c *sshConn) processChannelRequests() {
	for req := range c.newChannel {
		if req.ChannelType() != "session" {
			req.Reject(ssh.UnknownChannelType, "")
			continue
		}

		channel, requestCh, err := req.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		s := &session{channel, requestCh}
		go s.loop()
	}
}

func (c *sshConn) processGlobalRequests() {
	for req := range c.requests {
		if req.WantReply {
			req.Reply(false, nil)
		}
	}
}
