package tweeter

import (
	"fmt"
	"net/url"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/ReformedDevs/anonbot/db"
	"github.com/garyburd/go-oauth/oauth"
	"github.com/sirupsen/logrus"
)

// Tweeter takes care of sending tweets at regularly scheduled intervals. In
// addition, the type may be used to interact with the Twitter API.
type Tweeter struct {
	serverURL string
	database  *db.Connection
	log       *logrus.Entry
	triggerCh chan *db.Account
	stoppedCh chan bool
}

func (t *Tweeter) run() {
	defer close(t.stoppedCh)
	defer t.log.Info("tweeter has stopped")
	t.log.Info("starting tweeter...")
	var s *db.Schedule
	for {
		var nextTweetCh <-chan time.Time
		err := t.database.Transaction(func(c *db.Connection) error {
			if s == nil {
				s = t.selectSchedule(c)
			}
			if s != nil {
				q := t.selectQueuedItem(c, s)
				if q != nil {
					if err := t.tweet(c, s, q); err != nil {
						return err
					}
				}
			}
			nextTweetCh = t.nextTweetCh(c)
			return nil
		})
		if err != nil {
			t.log.Error(err.Error())
			nextTweetCh = time.After(30 * time.Second)
		}
		select {
		case <-nextTweetCh:
		case forceAccount, ok := <-t.triggerCh:
			if !ok {
				return
			}
			if forceAccount != nil {
				s = &db.Schedule{
					Account:   forceAccount,
					AccountID: forceAccount.ID,
				}
				continue
			}
		}
		s = nil
	}
}

// New creates a new tweeter instance from the specified configuration.
func New(cfg *Config) *Tweeter {
	anaconda.SetConsumerKey(cfg.ConsumerKey)
	anaconda.SetConsumerSecret(cfg.ConsumerSecret)
	t := &Tweeter{
		serverURL: cfg.ServerURL,
		database:  cfg.Database,
		log:       logrus.WithField("context", "tweeter"),
		triggerCh: make(chan *db.Account),
		stoppedCh: make(chan bool),
	}
	go t.run()
	return t
}

// Trigger hints to the tweeter that a new tweet is available or forces an
// account to tweet its next queued item.
func (t *Tweeter) Trigger(a *db.Account) {
	t.triggerCh <- a
}

// Authorize begins the authorization process for an account. The URL to
// redirect the user to and the temporary credentials are returned.
func (t *Tweeter) Authorize(a *db.Account) (string, string, string, error) {
	u, c, err := anaconda.AuthorizationURL(
		fmt.Sprintf("%s/accounts/%d/complete", t.serverURL, a.ID),
	)
	if err != nil {
		return "", "", "", err
	}
	return u, c.Token, c.Secret, nil
}

// Complete finishes the OAuth process. The access token and secret are
// returned.
func (t *Tweeter) Complete(accessToken, accessSecret, verifier string) (string, string, error) {
	c, _, err := anaconda.GetCredentials(&oauth.Credentials{
		Token:  accessToken,
		Secret: accessSecret,
	}, verifier)
	if err != nil {
		return "", "", err
	}
	return c.Token, c.Secret, nil
}

// Like favorites the specified tweet.
func (t *Tweeter) Like(a *db.Account, id int64) error {
	_, err := anaconda.NewTwitterApi(a.AccessToken, a.AccessSecret).Favorite(id)
	return err
}

// Mentions returns recent mentions for the account.
func (t *Tweeter) Mentions(a *db.Account) ([]anaconda.Tweet, error) {
	return anaconda.NewTwitterApi(a.AccessToken, a.AccessSecret).GetMentionsTimeline(nil)
}

// Tweet sends a tweet, optionally replying to a tweet with the specified ID.
func (t *Tweeter) Reply(a *db.Account, id, text string) error {
	v := url.Values{}
	if len(id) != 0 {
		v.Set("in_reply_to_status_id", id)
	}
	_, err := anaconda.NewTwitterApi(a.AccessToken, a.AccessSecret).PostTweet(text, v)
	return err
}

// Close shuts down the tweeter.
func (t *Tweeter) Close() {
	close(t.triggerCh)
	<-t.stoppedCh
}
