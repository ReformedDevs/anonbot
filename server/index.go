package server

import (
	"net/http"

	"github.com/ReformedDevs/anonbot/db"
	"github.com/flosch/pongo2"
)

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	var (
		numUsers       int64
		numAccounts    int64
		numSuggestions int64
		numQueuedItems int64
	)
	if err := s.database.C.Model(&db.User{}).Count(&numUsers).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := s.database.C.Model(&db.Account{}).Count(&numAccounts).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := s.database.C.Model(&db.Suggestion{}).Count(&numSuggestions).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := s.database.C.Model(&db.QueueItem{}).Count(&numQueuedItems).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.render(w, r, "index.html", pongo2.Context{
		"num_users":        numUsers,
		"num_accounts":     numAccounts,
		"num_suggestions":  numSuggestions,
		"num_queued_items": numQueuedItems,
	})
}
