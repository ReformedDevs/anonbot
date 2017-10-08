package db

import (
	"time"
)

// QueueItem represents a suggestion that has been approved and is queued for
// tweeting.
type QueueItem struct {
	ID        int64
	Date      time.Time
	Text      string   `gorm:"not null"`
	MediaURL  string   `gorm:"not null"`
	User      *User    `gorm:"ForeignKey:UserID"`
	UserID    int64    `sql:"type:int REFERENCES users(id) ON DELETE CASCADE"`
	Account   *Account `gorm:"ForeignKey:AccountID"`
	AccountID int64    `sql:"type:int REFERENCES accounts(id) ON DELETE CASCADE"`
}
