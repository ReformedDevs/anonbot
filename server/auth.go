package server

import (
	"context"
	"net/http"
	"net/url"

	"github.com/ReformedDevs/anonbot/db"
	"github.com/flosch/pongo2"
)

const (
	contextUser = "user"

	sessionName         = "session"
	sessionUserID       = "userID"
	sessionAccessToken  = "accessToken"
	sessionAccessSecret = "accessSecret"

	invalidCredentials = "invalid username or password"
)

func (s *Server) loadUser(r *http.Request) *http.Request {
	session, _ := s.store.Get(r, sessionName)
	v, _ := session.Values[sessionUserID]
	if v != nil {
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
			u := url.QueryEscape(r.URL.RequestURI())
			http.Redirect(w, r, "/login?url="+u, http.StatusFound)
			return
		}
		fn(w, r)
	}
}

func (s *Server) requireAdmin(fn http.HandlerFunc) http.HandlerFunc {
	return s.requireLogin(func(w http.ResponseWriter, r *http.Request) {
		if !r.Context().Value(contextUser).(*db.User).IsAdmin {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		fn(w, r)
	})
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	ctx := pongo2.Context{
		"title": "Login",
	}
	if r.Method == http.MethodPost {
		for {
			var (
				username = r.Form.Get("username")
				password = r.Form.Get("password")
				u        = &db.User{}
			)
			ctx["username"] = username
			ctx["password"] = password
			if err := s.database.C.Where("username = ?", username).First(&u).Error; err != nil {
				ctx["error"] = invalidCredentials
				break
			}
			if !u.IsActive {
				ctx["error"] = "account is not active"
				break
			}
			if err := u.Authenticate(password); err != nil {
				ctx["error"] = invalidCredentials
				break
			}
			session, _ := s.store.Get(r, sessionName)
			session.Values[sessionUserID] = u.ID
			session.Save(r, w)
			redir := r.URL.Query().Get("url")
			if len(redir) == 0 {
				redir = "/"
			}
			http.Redirect(w, r, redir, http.StatusFound)
			return
		}
	}
	s.render(w, r, "login.html", ctx)
}

func (s *Server) register(w http.ResponseWriter, r *http.Request) {
	ctx := pongo2.Context{
		"title": "Register",
	}
	if r.Method == http.MethodPost {
		for {
			var (
				username  = r.Form.Get("username")
				password  = r.Form.Get("password")
				password2 = r.Form.Get("password2")
				email     = r.Form.Get("email")
				u         = &db.User{
					Username: username,
					Email:    email,
				}
			)
			ctx["username"] = username
			ctx["email"] = email
			if len(username) == 0 || len(password) == 0 || len(email) == 0 {
				ctx["error"] = "invalid input"
				break
			}
			if password != password2 {
				ctx["error"] = "passwords do not match"
				break
			}
			if err := u.SetPassword(password); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if err := s.database.C.Create(u).Error; err != nil {
				ctx["error"] = "unable to create user"
				break
			}
			s.addAlert(w, r, "Thank you for registering. Please wait for an admin to activate your account.")
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
	}
	s.render(w, r, "register.html", ctx)
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	session, _ := s.store.Get(r, sessionName)
	session.Values[sessionUserID] = nil
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}
