package tweeter

import (
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/ReformedDevs/anonbot/db"
)

func (t *Tweeter) selectQueuedItem(c *db.Connection) (*db.Account, *db.QueueItem) {
	a := &db.Account{}
	if err := t.database.C.
		Where("queue_length > 0").
		Where("last_tweet + tweet_interval <= ?", time.Now().Unix()).
		First(a).Error; err != nil {
		return nil, nil
	}
	q := &db.QueueItem{}
	if err := t.database.C.
		Order("order").
		Where("account_id = ?", a.ID).
		First(q).Error; err != nil {
		return nil, nil
	}
	return a, q
}

func (t *Tweeter) tweet(a *db.Account, q *db.QueueItem) error {
	anaconda.SetConsumerKey(a.ConsumerKey)
	anaconda.SetConsumerSecret(a.ConsumerSecret)
	api := anaconda.NewTwitterApi(a.AccessToken, a.AccessSecret)
	_, err := api.PostTweet(q.Text, nil)
	if err != nil {
		return err
	}
	if err := t.database.C.Delete(q).Error; err != nil {
		return err
	}
	a.QueueLength--
	a.LastTweet = time.Now().Unix()
	if err := t.database.C.Save(a).Error; err != nil {
		return err
	}
	return nil
}
