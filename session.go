package gitolite

import (
	"errors"
	"fmt"
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
			return
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

type execRequest struct {
	Command string
}

func (s *session) handleExec(req *ssh.Request) error {
	cmd := execRequest{}
	if err := ssh.Unmarshal(req.Payload, &cmd); err != nil {
		return err
	}

	args, err := shellArgs([]byte(cmd.Command))
	if err != nil {
		return err
	}

	if err := s.exec(args); err != nil {
		return fmt.Errorf("%s: %s", err, cmd.Command)
	}
	return nil
}

func (s *session) exec(args []string) error {
	if len(args) != 2 {
		return errors.New("Invalid command")
	}

	switch args[0] {
	case "git-upload-pack":
		return nil
	case "git-receive-pack":
		return nil
	default:
		return errors.New("Invalid command")
	}
}
