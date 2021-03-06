package server

import (
	"net/http"

	"github.com/ReformedDevs/anonbot/db"
	"github.com/flosch/pongo2"
	"github.com/gorilla/mux"
)

func (s *Server) accounts(w http.ResponseWriter, r *http.Request) {
	accounts := []*db.Account{}
	if err := s.database.C.Find(&accounts).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.render(w, r, "accounts.html", pongo2.Context{
		"title":    "Accounts",
		"accounts": accounts,
	})
}

func (s *Server) viewAccount(w http.ResponseWriter, r *http.Request) {
	var (
		id = mux.Vars(r)["id"]
		a  = &db.Account{}
	)
	if err := s.database.C.Find(a, id).Error; err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	m, err := s.tweeter.Mentions(a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t := []*db.Tweet{}
	if err := s.database.C.
		Preload("User").
		Order("date desc").
		Where("account_id = ?", id).
		Find(&t).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.render(w, r, "viewaccount.html", pongo2.Context{
		"title":    a.Name,
		"account":  a,
		"mentions": m,
		"tweets":   t,
	})
}

type editAccountForm struct {
	Name string
}

func (s *Server) editAccount(w http.ResponseWriter, r *http.Request) {
	s.database.Transaction(func(c *db.Connection) error {
		var (
			id   = mux.Vars(r)["id"]
			a    = &db.Account{}
			form = &editAccountForm{}
			ctx  = pongo2.Context{
				"form": form,
			}
		)
		if len(id) != 0 {
			if err := c.C.Find(a, id).Error; err != nil {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return nil
			}
			s.copyStruct(a, form)
			ctx["title"] = "Edit Account"
			ctx["action"] = "Save"
		} else {
			ctx["title"] = "New Account"
			ctx["action"] = "Create"
		}
		if r.Method == http.MethodPost {
			s.populateStruct(r.Form, form)
			s.copyStruct(form, a)
			if err := c.C.Save(a).Error; err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return nil
			}
			http.Redirect(w, r, "/accounts", http.StatusFound)
			return nil
		}
		s.render(w, r, "editaccount.html", ctx)
		return nil
	})
}

func (s *Server) advanceAccount(w http.ResponseWriter, r *http.Request) {
	var (
		id = mux.Vars(r)["id"]
		a  = &db.Account{}
	)
	if err := s.database.C.Find(a, id).Error; err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if r.Method == http.MethodPost {
		s.tweeter.Trigger(a)
		http.Redirect(w, r, "/accounts", http.StatusFound)
		return
	}
	s.render(w, r, "confirm.html", pongo2.Context{
		"title":  "Advance Account",
		"action": "advance the account",
	})
}

func (s *Server) authorizeAccount(w http.ResponseWriter, r *http.Request) {
	var (
		id = mux.Vars(r)["id"]
		a  = &db.Account{}
	)
	if err := s.database.C.Find(a, id).Error; err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	u, accessToken, accessSecret, err := s.tweeter.Authorize(a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	session, _ := s.store.Get(r, sessionName)
	session.Values[sessionAccessToken] = accessToken
	session.Values[sessionAccessSecret] = accessSecret
	session.Save(r, w)
	http.Redirect(w, r, u, http.StatusFound)
}

func (s *Server) completeAccount(w http.ResponseWriter, r *http.Request) {
	s.database.Transaction(func(c *db.Connection) error {
		var (
			id = mux.Vars(r)["id"]
			a  = &db.Account{}
		)
		if err := c.C.Find(a, id).Error; err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return nil
		}
		var (
			session, _      = s.store.Get(r, sessionName)
			accessToken, _  = session.Values[sessionAccessToken].(string)
			accessSecret, _ = session.Values[sessionAccessSecret].(string)
		)
		accessToken, accessSecret, err := s.tweeter.Complete(
			accessToken,
			accessSecret,
			r.FormValue("oauth_verifier"),
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return nil
		}
		a.AccessToken = accessToken
		a.AccessSecret = accessSecret
		if err := c.C.Save(a).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return nil
		}
		s.addAlert(w, r, "Account has been authorized.")
		http.Redirect(w, r, "/accounts", http.StatusFound)
		return nil
	})
}

func (s *Server) deleteAccount(w http.ResponseWriter, r *http.Request) {
	s.database.Transaction(func(c *db.Connection) error {
		var (
			id = mux.Vars(r)["id"]
			a  = &db.Account{}
		)
		if err := c.C.Find(a, id).Error; err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return nil
		}
		if r.Method == http.MethodPost {
			if err := c.C.Delete(a).Error; err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return nil
			}
			http.Redirect(w, r, "/accounts", http.StatusFound)
			return nil
		}
		s.render(w, r, "delete.html", pongo2.Context{
			"title": "Delete Account",
			"name":  a.Name,
		})
		return nil
	})
}
