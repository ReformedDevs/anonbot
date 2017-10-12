package tweeter

import (
	"net/url"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/ReformedDevs/anonbot/db"
)

func (t *Tweeter) tweet(c *db.Connection, a *db.Account, q *db.QueueItem) error {
	if err := c.C.Delete(q).Error; err != nil {
		return err
	}
	a.QueueLength--
	a.LastTweet = time.Now().Unix()
	if err := c.C.Save(a).Error; err != nil {
		return err
	}
	var (
		api = anaconda.NewTwitterApi(a.AccessToken, a.AccessSecret)
		v   = url.Values{}
	)
	if len(q.MediaURL) != 0 {
		t.log.Infof("uploading %s...", q.MediaURL)
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
	t.log.Infof("posting tweet for %s...", a.Name)
	i, err := api.PostTweet(q.Text, v)
	if err != nil {
		return err
	}
	var mediaURL string
	if len(i.Entities.Media) != 0 {
		mediaURL = i.Entities.Media[0].Media_url_https
	}
	return c.C.Save(&db.Tweet{
		TweetID:  i.Id,
		Date:     time.Now(),
		Text:     q.Text,
		MediaURL: mediaURL,
		UserID:   q.UserID,
	}).Error
}
