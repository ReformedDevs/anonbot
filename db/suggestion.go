package db

import (
	"time"
)

// Suggestion represents a tweet proposed by a user but not yet queued for
// tweeting by an admin.
type Suggestion struct {
	ID        int64
	Date      time.Time
	Text      string   `gorm:"not null"`
	MediaURL  string   `gorm:"not null"`
	User      *User    `gorm:"ForeignKey:UserID"`
	UserID    int64    `sql:"type:int REFERENCES users(id) ON DELETE CASCADE"`
	Account   *Account `gorm:"ForeignKey:AccountID"`
	AccountID int64    `sql:"type:int REFERENCES accounts(id) ON DELETE CASCADE"`
	Votes     []*Vote  `gorm:"ForeignKey:SuggestionID"`
	VoteCount int64    `gorm:"not null"`
}
