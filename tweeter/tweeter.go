package tweeter

import (
	"github.com/sirupsen/logrus"
)

// Tweeter takes care of sending tweets at regularly scheduled intervals.
type Tweeter struct {
	log       *logrus.Entry
	triggerCh chan bool
	stoppedCh chan bool
}

func (t *Tweeter) run() {
	defer close(t.stoppedCh)
	defer t.log.Info("tweeter has stopped")
	t.log.Info("starting tweeter...")
	for {
		select {
		case _, ok := <-t.triggerCh:
			if !ok {
				return
			}
		}
	}
}

// New creates a new tweeter instance.
func New() (*Tweeter, error) {
	t := &Tweeter{
		log:       logrus.WithField("context", "tweeter"),
		triggerCh: make(chan bool),
		stoppedCh: make(chan bool),
	}
	go t.run()
	return t, nil
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
