package server

import (
	"net/http"
	"strconv"

	"github.com/ReformedDevs/anonbot/db"
	"github.com/flosch/pongo2"
)

func (s *Server) accounts(w http.ResponseWriter, r *http.Request) {
	accounts := []*db.Account{}
	if err := s.database.C.Find(&accounts).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.render(w, r, "accounts.html", pongo2.Context{
		"title":    "Accounts",
		"accounts": accounts,
	})
}

func (s *Server) newAccount(w http.ResponseWriter, r *http.Request) {
	ctx := pongo2.Context{
		"title": "New Account",
	}
	if r.Method == http.MethodPost {
		for {
			var (
				name                = r.Form.Get("name")
				consumerKey         = r.Form.Get("consumer_key")
				consumerSecret      = r.Form.Get("consumer_secret")
				accessToken         = r.Form.Get("access_token")
				accessSecret        = r.Form.Get("access_secret")
				tweetInterval       = r.Form.Get("tweet_interval")
				tweetIntervalInt, _ = strconv.ParseInt(tweetInterval, 10, 64)
				a                   = &db.Account{
					Name:           name,
					ConsumerKey:    consumerKey,
					ConsumerSecret: consumerSecret,
					AccessToken:    accessToken,
					AccessSecret:   accessSecret,
					TweetInterval:  tweetIntervalInt,
				}
			)
			ctx["name"] = name
			ctx["consumer_key"] = consumerKey
			ctx["consumer_secret"] = consumerSecret
			ctx["access_token"] = accessToken
			ctx["access_secret"] = accessSecret
			ctx["tweet_interval"] = tweetInterval
			if len(name) == 0 {
				ctx["error"] = "invalid name"
				break
			}
			if err := s.database.C.Create(a).Error; err != nil {
				ctx["error"] = "unable to create account"
				break
			}
			http.Redirect(w, r, "/accounts", http.StatusTemporaryRedirect)
			return
		}
	}
	s.render(w, r, "newaccount.html", ctx)
}

func (s *Server) deleteAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if err := s.database.C.Where("id = ?", r.Form.Get("id")).Delete(&db.Account{}).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(w, r, "/accounts", http.StatusTemporaryRedirect)
}
