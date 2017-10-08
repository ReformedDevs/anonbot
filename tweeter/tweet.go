package tweeter

import (
	"net/url"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/ReformedDevs/anonbot/db"
)

func (t *Tweeter) tweet(c *db.Connection, a *db.Account, q *db.QueueItem) error {
	t.log.Infof("tweeting from %s...", a.Name)
	anaconda.SetConsumerKey(a.ConsumerKey)
	anaconda.SetConsumerSecret(a.ConsumerSecret)
	var (
		api = anaconda.NewTwitterApi(a.AccessToken, a.AccessSecret)
		v   = url.Values{}
	)
	if len(q.MediaURL) != 0 {
		d, err := t.retrieveMedia(q.MediaURL)
		if err != nil {
			return err
		}
		m, err := api.UploadMedia(d)
		if err != nil {
			return err
		}
		v.Set("media_ids", m.MediaIDString)
	}
	_, err := api.PostTweet(q.Text, v)
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
