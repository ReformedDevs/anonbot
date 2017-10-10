package db

// Vote represents a vote for a specific suggestion.
type Vote struct {
	ID           int64
	User         *User       `gorm:"ForeignKey:UserID"`
	UserID       int64       `sql:"type:int REFERENCES users(id) ON DELETE CASCADE"`
	Suggestion   *Suggestion `gorm:"ForeignKey:SuggestionID"`
	SuggestionID int64       `sql:"type:int REFERENCES suggestions(id) ON DELETE CASCADE"`
}
