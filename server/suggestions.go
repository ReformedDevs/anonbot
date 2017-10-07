package server

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ReformedDevs/anonbot/db"
	"github.com/flosch/pongo2"
	"github.com/gorilla/mux"
)

func (s *Server) suggestions(w http.ResponseWriter, r *http.Request) {
	suggestions := []*db.Suggestion{}
	if err := s.database.C.Preload("User").Preload("Account").Find(&suggestions).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.render(w, r, "suggestions.html", pongo2.Context{
		"title":       "Suggestions",
		"suggestions": suggestions,
	})
}

func (s *Server) newSuggestion(w http.ResponseWriter, r *http.Request) {
	accounts := []*db.Account{}
	if err := s.database.C.Find(&accounts).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx := pongo2.Context{
		"title":    "New Suggestion",
		"accounts": accounts,
	}
	if r.Method == http.MethodPost {
		for {
			var (
				text         = r.Form.Get("text")
				accountIDStr = r.Form.Get("account_id")
				accountID, _ = strconv.ParseInt(accountIDStr, 10, 64)
				su           = &db.Suggestion{
					Date:      time.Now(),
					Text:      text,
					UserID:    r.Context().Value(contextUser).(*db.User).ID,
					AccountID: accountID,
				}
			)
			ctx["text"] = text
			ctx["account_id"] = accountID
			if len(text) == 0 {
				ctx["error"] = "invalid text"
				break
			}
			if err := s.database.C.Create(su).Error; err != nil {
				ctx["error"] = "unable to create suggestion"
				break
			}
			http.Redirect(w, r, "/suggestions", http.StatusFound)
			return
		}
	}
	s.render(w, r, "newsuggestion.html", ctx)
}

func (s *Server) queueSuggestion(w http.ResponseWriter, r *http.Request) {
	s.database.Transaction(func(c *db.Connection) error {
		var (
			id  = mux.Vars(r)["id"]
			su  = &db.Suggestion{}
			ctx = pongo2.Context{
				"title":      "Queue Suggestion",
				"suggestion": su,
			}
		)
		if err := c.C.Preload("Account").Find(su, id).Error; err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return nil
		}
		ctx["text"] = su.Text
		if r.Method == http.MethodPost {
			for {
				if err := c.C.Delete(su).Error; err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return nil
				}
				q := &db.QueueItem{
					Order:     0,
					Text:      r.Form.Get("text"),
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
				http.Redirect(w, r, "/suggestions", http.StatusFound)
				return nil
			}
		}
		s.render(w, r, "queuesuggestion.html", ctx)
		return nil
	})
}
