package db

// Account represents an individual Twitter account with API credentials.
type Account struct {
	ID             int64
	Name           string `gorm:"not null"`
	ConsumerKey    string `gorm:"not null"`
	ConsumerSecret string `gorm:"not null"`
	AccessToken    string `gorm:"not null"`
	AccessSecret   string `gorm:"not null"`
}
