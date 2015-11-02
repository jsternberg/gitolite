package gitolite

import (
	"io"
	"log"

	"golang.org/x/crypto/ssh"
)

type session struct {
	channel  ssh.Channel
	requests <-chan *ssh.Request
}

func (s *session) loop() {
	defer s.channel.Close()
	for req := range s.requests {
		switch req.Type {
		case "exec":
			if err := s.handleExec(req); err != nil {
				log.Println(err)
				req.Reply(false, nil)
				return
			}
			req.Reply(true, nil)
		case "pty-req":
			// do not allocate a pty since we just want to allow the upcoming shell request
			// so we can return a version string
			req.Reply(true, nil)
		case "shell":
			io.WriteString(s.channel, "Authentication succeeded.\n")
			req.Reply(true, nil)
			return
		default:
			if req.WantReply {
				req.Reply(false, nil)
			}
		}
	}
}

func (s *session) handleExec(req *ssh.Request) error {
	req.Reply(false, nil)
	return nil
}
