package db

// Setting represents a configurable value in the database.
type Setting struct {
	Key   string `gorm:"primary_key"`
	Value string `gorm:"not null"`
}

// Setting returns the value of the specified setting, using the default value
// if the key does not exist.
func (c *Connection) Setting(key, defaultValue string) string {
	s := &Setting{}
	if err := c.C.Where("key = ?", key).Find(&s).Error; err != nil {
		return defaultValue
	}
	return s.Value
}
