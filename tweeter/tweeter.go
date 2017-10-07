package tweeter

import (
	"time"

	"github.com/ReformedDevs/anonbot/db"
	"github.com/sirupsen/logrus"
)

// Tweeter takes care of sending tweets at regularly scheduled intervals.
type Tweeter struct {
	database  *db.Connection
	log       *logrus.Entry
	triggerCh chan bool
	stoppedCh chan bool
}

func (t *Tweeter) run() {
	defer close(t.stoppedCh)
	defer t.log.Info("tweeter has stopped")
	t.log.Info("starting tweeter...")
	for {
		var (
			nextTweetCh <-chan time.Time
		)
		err := t.database.Transaction(func(c *db.Connection) error {
			a, q := t.selectQueuedItem(c)
			if a != nil && q != nil {
				return t.tweet(a, q)
			}
			return nil
		})
		if err != nil {
			t.log.Error(err.Error())
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

// Close shuts down the tweeter.
func (t *Tweeter) Close() {
	close(t.triggerCh)
	<-t.stoppedCh
}
