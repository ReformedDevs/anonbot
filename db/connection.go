package db

import (
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Connection represents a socket connection to the database. The C member may
// be used directly to perform queries.
type Connection struct {
	C *gorm.DB
}

// Connect makes an attempt to connect to the database using the provided
// configuration.
func Connect(cfg *Config) (*Connection, error) {
	d, err := gorm.Open(cfg.Driver, cfg.Args)
	if err != nil {
		return nil, err
	}
	return &Connection{
		C: d,
	}, nil
}

// Migrate performs all pending database migrations.
func (c *Connection) Migrate() error {
	return c.C.AutoMigrate(
		&User{},
		&Account{},
		&Suggestion{},
		&QueueItem{},
		&Tweet{},
	).Error
}

// Transaction executes the provided callback in a transaction, automatically
// rolling back the database if an error occurs and committing the changes if
// none occurs.
func (c *Connection) Transaction(fn func(*Connection) error) error {
	d := c.C.Begin()
	if err := fn(&Connection{C: d}); err != nil {
		d.Rollback()
		return err
	}
	d.Commit()
	return nil
}