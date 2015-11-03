package gitolite

import (
	"log"
	"net"
)

type Server struct {
	config *Config
}

// Instantiates a Server with the given Config.
// The Config should not be modified after being given to the server
// to avoid any race conditions as the Config is used to store information
// at runtime.
//
// This does not cause the server to start listening or performing any other
// actions at the moment.
func New(config *Config) *Server {
	return &Server{config: config}
}

// This serves the git repository over SSH on the passed in listener.
// The method will block until the listener is closed or an error occurs
// during the 'Accept()' call.
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

// This is a shortcut method for opening a TCP socket and calling Serve.
func (s *Server) ListenAndServe(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return s.Serve(l)
}
