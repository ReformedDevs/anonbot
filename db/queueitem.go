package db

// QueueItem represents a suggestion that has been approved and is queued for
// tweeting.
type QueueItem struct {
	ID           int64
	Order        int
	Suggestion   *Suggestion `gorm:"ForeignKey:SuggestionID"`
	SuggestionID int64       `sql:"type:int REFERENCES suggestions(id) ON DELETE CASCADE"`
}
