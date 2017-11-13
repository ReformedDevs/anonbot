package tweeter

import (
	"time"

	"github.com/ReformedDevs/anonbot/db"
)

// selectSchedule attempts to find a schedule for an account whose next tweet
// should have been sent in the past or present. The Account member is filled
// in.
func (t *Tweeter) selectSchedule(c *db.Connection) *db.Schedule {
	s := &db.Schedule{}
	if err := c.C.
		Preload("Account").
		Order("next_run").
		Where("next_run <= ?", time.Now()).
		First(s).Error; err != nil {
		return nil
	}
	return s
}

// selectQueuedItem retrieves the next available suggestion for the specified
// schedule. If no item is available, the schedule is updated to its next time
// slot.
func (t *Tweeter) selectQueuedItem(c *db.Connection, s *db.Schedule) (*db.QueueItem, error) {
	q := &db.QueueItem{}
	if err := c.C.
		Order("date").
		Where("account_id = ?", s.AccountID).
		First(q).Error; err != nil {
		if s.ID != 0 {
			if err := s.Calculate(); err != nil {
				return nil, err
			}
			if err := c.C.Save(s).Error; err != nil {
				return nil, err
			}
		}
		return nil, nil
	}
	return q, nil
}

// nextTweetCh creates a channel that sends when the next tweet should be sent.
func (t *Tweeter) nextTweetCh(c *db.Connection) <-chan time.Time {
	s := &db.Schedule{}
	if err := c.C.
		Preload("Account").
		Joins("LEFT JOIN accounts ON schedules.account_id = accounts.id").
		Where("accounts.queue_length > 0").
		First(s).Error; err != nil {
		return nil
	}
	nextTweet := s.NextRun.Sub(time.Now())
	t.log.Debugf("next tweet from %s in %s", s.Account.Name, nextTweet.String())
	return time.After(nextTweet)
}
