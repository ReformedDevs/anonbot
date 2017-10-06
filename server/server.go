package server

import (
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

// Server provides the web UI for interacting with the application. Users can
// login, post suggestions, and queue items if they have the appropriate
// permissions.
type Server struct {
	listener net.Listener
	store    *sessions.CookieStore
	log      *logrus.Entry
	stopped  chan bool
}

// New creates a new server with the specified configuration.
func New(cfg *Config) (*Server, error) {
	l, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		return nil, err
	}
	var (
		router = mux.NewRouter()
		server = http.Server{
			Handler: router,
		}
		s = &Server{
			listener: l,
			store:    sessions.NewCookieStore([]byte(cfg.SecretKey)),
			log:      logrus.WithField("context", "server"),
			stopped:  make(chan bool),
		}
	)
	router.PathPrefix("/static").Handler(http.FileServer(HTTP))
	go func() {
		defer close(s.stopped)
		defer s.log.Info("web server has stopped")
		s.log.Info("starting web server...")
		if err := server.Serve(l); err != nil {
			s.log.Error(err.Error())
		}
	}()
	return s, nil
}

// Close shuts down the server and waits for it to complete.
func (s *Server) Close() {
	s.listener.Close()
	<-s.stopped
}
