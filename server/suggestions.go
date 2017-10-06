package server

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ReformedDevs/anonbot/db"
	"github.com/flosch/pongo2"
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
			http.Redirect(w, r, "/suggestions", http.StatusTemporaryRedirect)
			return
		}
	}
	s.render(w, r, "newsuggestion.html", ctx)
}
