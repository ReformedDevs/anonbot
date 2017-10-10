package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ReformedDevs/anonbot/db"
)

type ajaxAction struct {
	Action string          `json:"action"`
	Data   json.RawMessage `json:"data"`
}

type ajaxActionReply struct {
	AccountID string `json:"account_id"`
	TweetID   string `json:"tweet_id"`
	Text      string `json:"text"`
}

type ajaxActionLike struct {
	AccountID string `json:"account_id"`
	TweetID   string `json:"tweet_id"`
}

// TODO: this could probably be refactored a bit

func (s *Server) ajax(w http.ResponseWriter, r *http.Request) {
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
		if !r.Context().Value(contextUser).(*db.User).IsAdmin {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		reply := &ajaxActionReply{}
		if err := json.Unmarshal(action.Data, reply); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		a := &db.Account{}
		if err := s.database.C.Find(a, reply.AccountID).Error; err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		if err := s.tweeter.Reply(a, reply.TweetID, reply.Text); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "like":
		if !r.Context().Value(contextUser).(*db.User).IsAdmin {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		like := &ajaxActionLike{}
		if err := json.Unmarshal(action.Data, like); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		a := &db.Account{}
		if err := s.database.C.Find(a, like.AccountID).Error; err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		tweetID, _ := strconv.ParseInt(like.TweetID, 10, 64)
		if err := s.tweeter.Like(a, tweetID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Length", "0")
	w.WriteHeader(http.StatusOK)
}
