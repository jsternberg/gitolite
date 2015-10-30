package gitolite

import "net"

type Server struct{}

func New(config *Config) *Server {
	return nil
}

func (s *Server) Serve(l net.Listener) error {
	return nil
}

func (s *Server) ListenAndServe(addr string) error {
	return nil
}
