package gitolite

import (
	"log"
	"net"

	"golang.org/x/crypto/ssh"
)

type sshConn struct {
	conn       *ssh.ServerConn
	newChannel <-chan ssh.NewChannel
	requests   <-chan *ssh.Request
}

func (s *Server) serve(conn net.Conn) error {
	config := ssh.ServerConfig{
		PublicKeyCallback: s.config.PublicKeyCallback.wrap(),
		PasswordCallback:  s.config.PasswordCallback.wrap(),
	}

	for _, hostKey := range s.config.HostKeys {
		config.AddHostKey(hostKey)
	}

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
