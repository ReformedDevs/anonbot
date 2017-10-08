package server

import (
	"net"
	"net/http"

	"github.com/ReformedDevs/anonbot/db"
	"github.com/ReformedDevs/anonbot/tweeter"
	"github.com/flosch/pongo2"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"

	_ "github.com/flosch/pongo2-addons"
)

// Server provides the web UI for interacting with the application. Users can
// login, post suggestions, and queue items if they have the appropriate
// permissions.
type Server struct {
	listener    net.Listener
	database    *db.Connection
	tweeter     *tweeter.Tweeter
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
		router = mux.NewRouter()
		s      = &Server{
			listener:    l,
			database:    cfg.Database,
			tweeter:     cfg.Tweeter,
			router:      mux.NewRouter(),
			store:       sessions.NewCookieStore([]byte(cfg.SecretKey)),
			templateSet: pongo2.NewSet("", &b0xLoader{}),
			log:         logrus.WithField("context", "server"),
			stoppedCh:   make(chan bool),
		}
		server = http.Server{
			Handler: router,
		}
	)
	s.router.HandleFunc("/", s.index)
	s.router.HandleFunc("/login", s.login)
	s.router.HandleFunc("/register", s.register)
	s.router.HandleFunc("/logout", s.requireLogin(s.logout))
	s.router.HandleFunc("/suggestions", s.requireLogin(s.suggestions))
	s.router.HandleFunc("/suggestions/new", s.requireLogin(s.editSuggestion))
	s.router.HandleFunc("/suggestions/{id:[0-9]+}/edit", s.requireLogin(s.editSuggestion))
	s.router.HandleFunc("/suggestions/{id:[0-9]+}/queue", s.requireAdmin(s.queueSuggestion))
	s.router.HandleFunc("/suggestions/{id:[0-9]+}/delete", s.requireLogin(s.deleteSuggestion))
	s.router.HandleFunc("/queue", s.requireLogin(s.queue))
	s.router.HandleFunc("/queue/{id:[0-9]+}/delete", s.requireAdmin(s.deleteQueueItem))
	s.router.HandleFunc("/accounts", s.requireAdmin(s.accounts))
	s.router.HandleFunc("/accounts/new", s.requireAdmin(s.editAccount))
	s.router.HandleFunc("/accounts/{id:[0-9]+}/edit", s.requireAdmin(s.editAccount))
	s.router.HandleFunc("/accounts/{id:[0-9]+}/delete", s.requireAdmin(s.deleteAccount))
	s.router.HandleFunc("/users", s.requireAdmin(s.users))
	s.router.HandleFunc("/users/new", s.requireAdmin(s.editUser))
	s.router.HandleFunc("/users/{id:[0-9]+}/edit", s.requireAdmin(s.editUser))
	s.router.HandleFunc("/users/{id:[0-9]+}/delete", s.requireAdmin(s.deleteUser))
	router.PathPrefix("/static").Handler(http.FileServer(HTTP))
	router.PathPrefix("/").Handler(s)
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
