package server

import (
	"context"
	"net/http"

	"github.com/ReformedDevs/anonbot/db"
	"github.com/flosch/pongo2"
)

const (
	contextUser = "user"

	sessionName   = "session"
	sessionUserID = "userID"
)

func (s *Server) loadUser(r *http.Request) *http.Request {
	session, _ := s.store.Get(r, sessionName)
	v, _ := session.Values[sessionUserID]
	if v != "" {
		u := &db.User{}
		if err := s.database.C.First(u, v).Error; err == nil {
			r = r.WithContext(context.WithValue(r.Context(), contextUser, u))
		}
	}
	return r
}

func (s *Server) requireLogin(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Context().Value(contextUser) == nil {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}
		fn(w, r)
	}
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	ctx := pongo2.Context{
		"title": "Login",
	}
	if r.Method == http.MethodPost {
		var (
			username = r.Form.Get("username")
			password = r.Form.Get("password")
			u        = &db.User{}
		)
		if err := s.database.C.Where("username = ?", username).First(&u).Error; err == nil {
			if err := u.Authenticate(password); err == nil {
				session, _ := s.store.Get(r, sessionName)
				session.Values[sessionUserID] = u.ID
				session.Save(r, w)
				http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
				return
			}
		}
		ctx["error"] = "invalid username or password"
	}
	s.render(w, r, "login.html", ctx)
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	session, _ := s.store.Get(r, sessionName)
	session.Values[sessionUserID] = ""
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
