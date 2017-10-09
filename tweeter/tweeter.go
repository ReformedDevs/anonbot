package tweeter

import (
	"net/url"
	"sync"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/ReformedDevs/anonbot/db"
	"github.com/sirupsen/logrus"
)

// Tweeter takes care of sending tweets at regularly scheduled intervals. In
// addition, the type may be used to interact with the Twitter API.
type Tweeter struct {
	mutex     sync.Mutex
	database  *db.Connection
	log       *logrus.Entry
	triggerCh chan bool
	stoppedCh chan bool
}

func (t *Tweeter) activate(a *db.Account) *anaconda.TwitterApi {
	anaconda.SetConsumerKey(a.ConsumerKey)
	anaconda.SetConsumerSecret(a.ConsumerSecret)
	return anaconda.NewTwitterApi(a.AccessToken, a.AccessSecret)
}

func (t *Tweeter) run() {
	defer close(t.stoppedCh)
	defer t.log.Info("tweeter has stopped")
	t.log.Info("starting tweeter...")
	for {
		var nextTweetCh <-chan time.Time
		err := t.database.Transaction(func(c *db.Connection) error {
			a, q := t.selectQueuedItem(c)
			if a != nil && q != nil {
				if err := t.tweet(c, a, q); err != nil {
					return err
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
		case _, ok := <-t.triggerCh:
			if !ok {
				return
			}
		}
	}
}

// New creates a new tweeter instance from the specified configuration.
func New(cfg *Config) *Tweeter {
	t := &Tweeter{
		database:  cfg.Database,
		log:       logrus.WithField("context", "tweeter"),
		triggerCh: make(chan bool),
		stoppedCh: make(chan bool),
	}
	go t.run()
	return t
}

// Trigger hints to the tweeter that a new tweet is available in the database.
func (t *Tweeter) Trigger() {
	t.triggerCh <- true
}

// Mentions returns recent mentions for the account.
func (t *Tweeter) Mentions(a *db.Account) ([]anaconda.Tweet, error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return t.activate(a).GetMentionsTimeline(nil)
}

// Reply replies to the tweet with the specified ID.
func (t *Tweeter) Reply(a *db.Account, id, text string) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	_, err := t.activate(a).PostTweet(text, url.Values{
		"in_reply_to_status_id": []string{id},
	})
	return err
}

// Close shuts down the tweeter.
func (t *Tweeter) Close() {
	close(t.triggerCh)
	<-t.stoppedCh
}
