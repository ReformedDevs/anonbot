package server

import (
	"net/http"

	"github.com/ReformedDevs/anonbot/db"
	"github.com/flosch/pongo2"
)

const (
	sessionName   = "session"
	sessionUserID = "userID"
)

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
