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

type ajaxActionVote struct {
	SuggestionID string `json:"suggestion_id"`
}

type ajaxActionTweet struct {
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
	u := r.Context().Value(contextUser).(*db.User)
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
	case "vote":
		vote := &ajaxActionVote{}
		if err := json.Unmarshal(action.Data, vote); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		s.database.Transaction(func(c *db.Connection) error {
			s := &db.Suggestion{}
			if err := c.C.
				Preload("Votes", "user_id = ?", u.ID).
				Find(s, vote.SuggestionID).Error; err != nil {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return nil
			}
			if len(s.Votes) == 0 {
				v := &db.Vote{
					UserID: u.ID,
				}
				if err := c.C.Model(s).Association("Votes").Append(v).Error; err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return nil
				}
				s.VoteCount++
			} else {
				if err := c.C.Model(s).Association("Votes").Delete(s.Votes[0]).Error; err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return nil
				}
				s.VoteCount--
			}
			if err := c.C.Save(s).Error; err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return nil
			}
			w.Header().Set("Content-Length", "0")
			w.WriteHeader(http.StatusOK)
			return nil
		})
	case "tweet":
		if !u.IsAdmin {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		tweet := &ajaxActionTweet{}
		if err := json.Unmarshal(action.Data, tweet); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		a := &db.Account{}
		if err := s.database.C.Find(a, tweet.AccountID).Error; err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		if err := s.tweeter.Reply(a, tweet.TweetID, tweet.Text); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Length", "0")
		w.WriteHeader(http.StatusOK)
	case "like":
		if !u.IsAdmin {
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
		w.Header().Set("Content-Length", "0")
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}
