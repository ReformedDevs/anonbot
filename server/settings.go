package server

import (
	"net/http"

	"github.com/ReformedDevs/anonbot/db"
	"github.com/flosch/pongo2"
)

func (s *Server) settings(w http.ResponseWriter, r *http.Request) {
	var (
		tweetInterval = s.database.Setting(db.TweetInterval, db.TweetIntervalDefault)
	)
	s.render(w, r, "settings.html", pongo2.Context{
		"title":         "Settings",
		"tweetInterval": tweetInterval,
	})
}
