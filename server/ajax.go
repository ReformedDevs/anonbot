package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ReformedDevs/anonbot/db"
	"github.com/gorilla/mux"
)

type ajaxAction struct {
	Action string          `json:"action"`
	Data   json.RawMessage `json:"data"`
}

type ajaxActionReply struct {
	TweetID string `json:"tweet_id"`
	Text    string `json:"text"`
}

type ajaxActionLike struct {
	TweetID string `json:"tweet_id"`
}

func (s *Server) ajaxAccount(w http.ResponseWriter, r *http.Request) {
	var (
		id = mux.Vars(r)["id"]
		a  = &db.Account{}
	)
	if err := s.database.C.Find(a, id).Error; err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	action := &ajaxAction{}
	if err := json.NewDecoder(r.Body).Decode(action); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	switch action.Action {
	case "reply":
		reply := &ajaxActionReply{}
		if err := json.Unmarshal(action.Data, reply); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := s.tweeter.Reply(a, reply.TweetID, reply.Text); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "like":
		like := &ajaxActionLike{}
		if err := json.Unmarshal(action.Data, like); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		tweetID, _ := strconv.ParseInt(like.TweetID, 10, 64)
		if err := s.tweeter.Like(a, tweetID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Length", "0")
	w.WriteHeader(http.StatusOK)
}
