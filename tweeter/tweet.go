package tweeter

import (
	"net/url"
	"time"

	"github.com/ReformedDevs/anonbot/db"
)

func (t *Tweeter) tweet(c *db.Connection, a *db.Account, q *db.QueueItem) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.log.Infof("tweeting from %s...", a.Name)
	var (
		api = t.activate(a)
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
