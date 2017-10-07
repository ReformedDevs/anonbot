package db

import (
	"time"
)

// Account represents an individual Twitter account with API credentials.
type Account struct {
	ID             int64
	Name           string `gorm:"not null"`
	ConsumerKey    string `gorm:"not null"`
	ConsumerSecret string `gorm:"not null"`
	AccessToken    string `gorm:"not null"`
	AccessSecret   string `gorm:"not null"`
	QueueLength    int64
	TweetInterval  int64
	LastTweet      int64
}

// LastTweetDate returns the time of the last tweet as a time.Time.
func (a *Account) LastTweetDate() time.Time {
	return time.Unix(a.LastTweet, 0)
}
