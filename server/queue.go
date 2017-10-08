package server

import (
	"net/http"

	"github.com/ReformedDevs/anonbot/db"
	"github.com/flosch/pongo2"
	"github.com/gorilla/mux"
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

func (s *Server) deleteQueueItem(w http.ResponseWriter, r *http.Request) {
	s.database.Transaction(func(c *db.Connection) error {
		var (
			id = mux.Vars(r)["id"]
			q  = &db.QueueItem{}
		)
		if err := c.C.Find(q, id).Error; err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return nil
		}
		if r.Method == http.MethodPost {
			if err := c.C.Delete(q).Error; err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return nil
			}
			http.Redirect(w, r, "/queue", http.StatusFound)
			return nil
		}
		s.render(w, r, "delete.html", pongo2.Context{
			"title": "Delete Queue Item",
			"name":  "this queue item",
		})
		return nil
	})
}
