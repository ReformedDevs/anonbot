package server

import (
	"net"
	"net/http"

	"github.com/ReformedDevs/anonbot/db"
	"github.com/flosch/pongo2"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

// Server provides the web UI for interacting with the application. Users can
// login, post suggestions, and queue items if they have the appropriate
// permissions.
type Server struct {
	listener    net.Listener
	database    *db.Connection
	router      *mux.Router
	store       *sessions.CookieStore
	templateSet *pongo2.TemplateSet
	log         *logrus.Entry
	stoppedCh   chan bool
}

// New creates a new server with the specified configuration.
func New(cfg *Config) (*Server, error) {
	l, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		return nil, err
	}
	var (
		s = &Server{
			listener:    l,
			database:    cfg.Database,
			router:      mux.NewRouter(),
			store:       sessions.NewCookieStore([]byte(cfg.SecretKey)),
			templateSet: pongo2.NewSet("", &b0xLoader{}),
			log:         logrus.WithField("context", "server"),
			stoppedCh:   make(chan bool),
		}
		server = http.Server{
			Handler: s,
		}
	)
	s.router.HandleFunc("/", s.index)
	s.router.HandleFunc("/login", s.login)
	s.router.HandleFunc("/register", s.register)
	s.router.HandleFunc("/logout", s.requireLogin(s.logout))
	s.router.HandleFunc("/suggestions", s.requireLogin(s.suggestions))
	s.router.HandleFunc("/suggestions/new", s.requireLogin(s.newSuggestion))
	s.router.HandleFunc("/accounts", s.requireAdmin(s.accounts))
	s.router.HandleFunc("/accounts/new", s.requireAdmin(s.newAccount))
	s.router.HandleFunc("/accounts/delete", s.requireAdmin(s.deleteAccount))
	s.router.HandleFunc("/settings", s.requireAdmin(s.settings))
	s.router.PathPrefix("/static").Handler(http.FileServer(HTTP))
	go func() {
		defer close(s.stoppedCh)
		defer s.log.Info("web server has stopped")
		s.log.Info("starting web server...")
		if err := server.Serve(l); err != nil {
			s.log.Error(err.Error())
		}
	}()
	return s, nil
}

// ServeHTTP does preparatory work for the handlers. It attempts to load the
// user from the database if authenticated and ensures that POST requests have
// their forms parsed.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	r = s.loadUser(r)
	s.router.ServeHTTP(w, r)
}

// Close shuts down the server and waits for it to complete.
func (s *Server) Close() {
	s.listener.Close()
	<-s.stoppedCh
}
