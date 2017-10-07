package server

import (
	"net/http"

	"github.com/ReformedDevs/anonbot/db"
	"github.com/flosch/pongo2"
)

func (s *Server) users(w http.ResponseWriter, r *http.Request) {
	users := []*db.User{}
	if err := s.database.C.Find(&users).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.render(w, r, "users.html", pongo2.Context{
		"title": "Users",
		"users": users,
	})
}
