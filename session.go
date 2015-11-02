package gitolite

import (
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
		case "shell":
			req.Reply(false, nil)
			return
		default:
			if req.WantReply {
				req.Reply(false, nil)
			}
			return
		}
	}
}

func (s *session) handleExec(req *ssh.Request) error {
	req.Reply(false, nil)
	return nil
}
