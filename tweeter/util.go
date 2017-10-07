package tweeter

import (
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/ReformedDevs/anonbot/db"
)

func (t *Tweeter) selectQueuedItem(c *db.Connection) (*db.Account, *db.QueueItem) {
	a := &db.Account{}
	if err := c.C.
		Where("queue_length > 0").
		Where("last_tweet + tweet_interval <= ?", time.Now().Unix()).
		First(a).Error; err != nil {
		return nil, nil
	}
	q := &db.QueueItem{}
	if err := c.C.
		Order("date").
		Where("account_id = ?", a.ID).
		First(q).Error; err != nil {
		return nil, nil
	}
	return a, q
}

func (t *Tweeter) tweet(c *db.Connection, a *db.Account, q *db.QueueItem) error {
	t.log.Infof("tweeting from %s...", a.Name)
	anaconda.SetConsumerKey(a.ConsumerKey)
	anaconda.SetConsumerSecret(a.ConsumerSecret)
	api := anaconda.NewTwitterApi(a.AccessToken, a.AccessSecret)
	_, err := api.PostTweet(q.Text, nil)
	if err != nil {
		return err
	}
	if err := c.C.Delete(q).Error; err != nil {
		return err
	}
	a.QueueLength--
	a.LastTweet = time.Now().Unix()
	if err := c.C.Save(a).Error; err != nil {
		return err
	}
	return nil
}

func (t *Tweeter) nextTweetCh(c *db.Connection) <-chan time.Time {
	a := &db.Account{}
	if err := c.C.
		Select("*, last_tweet + tweet_interval AS next_tweet_time").
		Order("next_tweet_time").
		Where("queue_length > 0").
		First(a).Error; err != nil {
		return nil
	}
	nextTweet := time.Unix(a.LastTweet+a.TweetInterval, 0).Sub(time.Now())
	t.log.Debugf("next tweet from %s in %s", a.Name, nextTweet.String())
	return time.After(nextTweet)
}
