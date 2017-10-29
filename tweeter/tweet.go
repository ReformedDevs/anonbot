package tweeter

import (
	"net/url"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/ReformedDevs/anonbot/db"
)

func (t *Tweeter) tweet(c *db.Connection, s *db.Schedule, q *db.QueueItem) error {
	if err := c.C.Delete(q).Error; err != nil {
		return err
	}
	s.Account.QueueLength--
	if err := c.C.Save(s.Account).Error; err != nil {
		return err
	}
	if s.ID != nil {
		if err := s.Calculate(); err != nil {
			return err
		}
		if err := c.C.Save(s).Error; err != nil {
			return err
		}
	}
	var (
		api = anaconda.NewTwitterApi(s.Account.AccessToken, s.Account.AccessSecret)
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
	t.log.Infof("posting tweet for %s...", s.Account.Name)
	i, err := api.PostTweet(q.Text, v)
	if err != nil {
		return err
	}
	var mediaURL string
	if len(i.Entities.Media) != 0 {
		mediaURL = i.Entities.Media[0].Media_url_https
	}
	return c.C.Save(&db.Tweet{
		TweetID:   i.Id,
		Date:      time.Now(),
		Text:      q.Text,
		MediaURL:  mediaURL,
		UserID:    q.UserID,
		AccountID: s.AccountID,
	}).Error
}
