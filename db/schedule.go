package db

import (
	"time"

	"github.com/robfig/cron"
)

// Schedule represents a specific time at which a queued tweet should be sent
// from an account. Cron notation is used for specifying the times. Each
// account may have multiple schedules.
type Schedule struct {
	ID        int64
	Cron      string    `gorm:"not null"`
	NextRun   time.Time `gorm:"not null"`
	Account   *Account  `gorm:"ForeignKey:AccountID"`
	AccountID int64     `sql:"type:int REFERENCES accounts(id) ON DELETE CASCADE"`
}

// Calculate determines when the next run should be based on the schedule. This
// method does NOT save the model.
func (s *Schedule) Calculate() error {
	p := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	v, err := p.Parse(s.Cron)
	if err != nil {
		return err
	}
	s.NextRun = v.Next(time.Now())
	return nil
}
