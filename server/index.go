package server

import (
	"net/http"

	"github.com/flosch/pongo2"
)

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	s.render(w, r, "index.html", pongo2.Context{})
}
