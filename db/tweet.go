package db

import (
	"time"
)

// Tweet represents a suggestion that has been tweeted.
type Tweet struct {
	ID           int64
	TweetID      int64
	Date         time.Time
	Suggestion   *Suggestion `gorm:"ForeignKey:SuggestionID"`
	SuggestionID int64       `sql:"type:int REFERENCES suggestions(id) ON DELETE CASCADE"`
}
