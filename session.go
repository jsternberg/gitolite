package gitolite

import (
	"errors"
	"fmt"
	"io"
	"log"
	"path/filepath"

	"golang.org/x/crypto/ssh"
)

type session struct {
	handler  interface{}
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
		log.Println("handleExec:", err)
		return fmt.Errorf("%s: %s", err, cmd.Command)
	}
	return nil
}

func (s *session) exec(args []string) error {
	if len(args) != 2 {
		return errors.New("Invalid command")
	}

	// Take the command from the first argument and the path from the second.
	// Join the path with the root path and clean it up. This should prevent
	// any attacks where ../.. are used to try and gain additional access to
	// the system because we are always assuming the path referenced is always
	// relative to the root. We want to do this clean up before we pass it to
	// the handler so the handler doesn't have to worry about validation.
	// This also means that the path passed to the underlying handler will
	// always have a slash at the beginning.
	command := args[0]
	path := filepath.Clean(filepath.Join("/", args[1]))

	w := channelWriter{s.channel}
	r := channelReader{s.channel}

	switch command {
	case "git-upload-pack":
		switch handler := s.handler.(type) {
		case UploadPack:
			return handler.UploadPack(path, &r, &w)
		default:
			return errors.New("Unsupported operation")
		}
	case "git-receive-pack":
		switch handler := s.handler.(type) {
		case ReceivePack:
			return handler.ReceivePack(path, &r, &w)
		default:
			return errors.New("Unsupported operation")
		}
	default:
		return errors.New("Invalid command")
	}
}
