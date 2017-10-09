package server

import (
	"net/http"
)

func (s *Server) addAlert(w http.ResponseWriter, r *http.Request, text string) {
	session, _ := s.store.Get(r, sessionName)
	session.AddFlash(text)
	session.Save(r, w)
}

func (s *Server) getAlerts(w http.ResponseWriter, r *http.Request) interface{} {
	session, _ := s.store.Get(r, sessionName)
	defer session.Save(r, w)
	return session.Flashes()
}
