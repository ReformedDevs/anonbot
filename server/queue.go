package server

import (
	"net/http"

	"github.com/ReformedDevs/anonbot/db"
	"github.com/flosch/pongo2"
)

func (s *Server) queue(w http.ResponseWriter, r *http.Request) {
	queue := []*db.QueueItem{}
	if err := s.database.C.
		Order("account_id, date").
		Preload("User").
		Preload("Account").
		Find(&queue).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.render(w, r, "queue.html", pongo2.Context{
		"title": "Queue",
		"queue": queue,
	})
}
