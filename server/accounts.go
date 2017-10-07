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

type editAccountForm struct {
	Name           string
	ConsumerKey    string
	ConsumerSecret string
	AccessToken    string
	AccessSecret   string
	TweetInterval  int64
}

func (s *Server) editAccount(w http.ResponseWriter, r *http.Request) {
	s.database.Transaction(func(c *db.Connection) error {
		var (
			id   = mux.Vars(r)["id"]
			a    = &db.Account{}
			form = &editAccountForm{
				TweetInterval: 86400,
			}
			ctx = pongo2.Context{
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
			s.tweeter.Trigger()
			http.Redirect(w, r, "/accounts", http.StatusFound)
			return nil
		}
		s.render(w, r, "editaccount.html", ctx)
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