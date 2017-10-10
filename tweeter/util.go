package tweeter

import (
	"time"

	"github.com/ReformedDevs/anonbot/db"
)

func (t *Tweeter) selectAccount(c *db.Connection) *db.Account {
	a := &db.Account{}
	if err := c.C.
		Where("queue_length > 0").
		Where("last_tweet + tweet_interval <= ?", time.Now().Unix()).
		First(a).Error; err != nil {
		return nil
	}
	return a
}

func (t *Tweeter) selectQueuedItem(c *db.Connection, a *db.Account) *db.QueueItem {
	q := &db.QueueItem{}
	if err := c.C.
		Order("date").
		Where("account_id = ?", a.ID).
		First(q).Error; err != nil {
		return nil
	}
	return q
}

func (t *Tweeter) nextTweetCh(c *db.Connection) <-chan time.Time {
	a := &db.Account{}
	if err := c.C.
		Select("*, last_tweet + tweet_interval AS next_tweet_time").
		Order("next_tweet_time").
		Where("queue_length > 0").
		First(a).Error; err != nil {
		t.log.Debug("no upcoming tweets")
		return nil
	}
	nextTweet := time.Unix(a.LastTweet+a.TweetInterval, 0).Sub(time.Now())
	t.log.Debugf("next tweet from %s in %s", a.Name, nextTweet.String())
	return time.After(nextTweet)
}
