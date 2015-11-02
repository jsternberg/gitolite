package gitolite

import (
	"log"
	"net"
)

type Server struct {
	config *Config
}

func New(config *Config) *Server {
	return &Server{config: config}
}

func (s *Server) Serve(l net.Listener) error {
	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}

		// The serve method here will launch the necessary go procs
		if err = s.serve(conn); err != nil {
			log.Println(err)
		}
	}
	return nil
}

func (s *Server) ListenAndServe(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return s.Serve(l)
}
