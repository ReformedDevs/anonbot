package server

import (
	"net/http"
	"time"

	"github.com/ReformedDevs/anonbot/db"
	"github.com/flosch/pongo2"
	"github.com/gorilla/mux"
)

func (s *Server) suggestions(w http.ResponseWriter, r *http.Request) {
	var (
		u           = r.Context().Value(contextUser).(*db.User)
		order       = "account_id, date"
		suggestions = []*db.Suggestion{}
	)
	if r.FormValue("order") == "votes" {
		order = "account_id, vote_count desc"
	}
	if err := s.database.C.
		Preload("User").
		Preload("Account").
		Preload("Votes", "user_id = ?", u.ID).
		Order(order).
		Find(&suggestions).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.render(w, r, "suggestions.html", pongo2.Context{
		"title":       "Suggestions",
		"suggestions": suggestions,
	})
}

type editSuggestionForm struct {
	Text      string
	MediaURL  string
	AccountID int64
}

func (s *Server) editSuggestion(w http.ResponseWriter, r *http.Request) {
	s.database.Transaction(func(c *db.Connection) error {
		var (
			id = mux.Vars(r)["id"]
			u  = r.Context().Value(contextUser).(*db.User)
			su = &db.Suggestion{
				Date:   time.Now(),
				UserID: u.ID,
			}
			form = &editSuggestionForm{}
			ctx  = pongo2.Context{
				"form": form,
			}
		)
		if len(id) != 0 {
			if err := c.C.Find(su, id).Error; err != nil {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return nil
			}
			if su.UserID != u.ID && !u.IsAdmin {
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return nil
			}
			s.copyStruct(su, form)
			ctx["title"] = "Edit Suggestion"
			ctx["action"] = "Save"
		} else {
			ctx["title"] = "New Suggestion"
			ctx["action"] = "Create"
		}
		if r.Method == http.MethodPost {
			s.populateStruct(r.Form, form)
			s.copyStruct(form, su)
			if err := c.C.Save(su).Error; err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return nil
			}
			http.Redirect(w, r, "/suggestions", http.StatusFound)
			return nil
		}
		accounts := []*db.Account{}
		if err := c.C.Find(&accounts).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return nil
		}
		ctx["accounts"] = accounts
		s.render(w, r, "editsuggestion.html", ctx)
		return nil
	})
}

func (s *Server) queueSuggestion(w http.ResponseWriter, r *http.Request) {
	s.database.Transaction(func(c *db.Connection) error {
		var (
			id = mux.Vars(r)["id"]
			su = &db.Suggestion{}
		)
		if err := c.C.Preload("Account").Find(su, id).Error; err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return nil
		}
		if r.Method == http.MethodPost {
			for {
				if err := c.C.Delete(su).Error; err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return nil
				}
				q := &db.QueueItem{
					Date:      time.Now(),
					Text:      su.Text,
					MediaURL:  su.MediaURL,
					UserID:    su.UserID,
					AccountID: su.AccountID,
				}
				if err := c.C.Save(q).Error; err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return nil
				}
				su.Account.QueueLength++
				if err := c.C.Save(su.Account).Error; err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return nil
				}
				s.tweeter.Trigger(nil)
				http.Redirect(w, r, "/suggestions", http.StatusFound)
				return nil
			}
		}
		s.render(w, r, "confirm.html", pongo2.Context{
			"title":  "Queue Suggestion",
			"action": "queue this tweet",
		})
		return nil
	})
}

func (s *Server) deleteSuggestion(w http.ResponseWriter, r *http.Request) {
	s.database.Transaction(func(c *db.Connection) error {
		var (
			id = mux.Vars(r)["id"]
			u  = r.Context().Value(contextUser).(*db.User)
			su = &db.Suggestion{}
		)
		if err := c.C.Find(su, id).Error; err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return nil
		}
		if su.UserID != u.ID && !u.IsAdmin {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return nil
		}
		if r.Method == http.MethodPost {
			if err := c.C.Delete(su).Error; err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return nil
			}
			http.Redirect(w, r, "/suggestions", http.StatusFound)
			return nil
		}
		s.render(w, r, "delete.html", pongo2.Context{
			"title": "Delete Suggestion",
			"name":  "this suggestion",
		})
		return nil
	})
}
