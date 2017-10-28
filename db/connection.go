package db

import (
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Connection represents a socket connection to the database. The C member may
// be used directly to perform queries.
type Connection struct {
	C   *gorm.DB
	log *logrus.Entry
}

// Connect makes an attempt to connect to the database using the provided
// configuration.
func Connect(cfg *Config) (*Connection, error) {
	d, err := gorm.Open(cfg.Driver, cfg.Args)
	if err != nil {
		return nil, err
	}
	c := &Connection{
		C:   d,
		log: logrus.WithField("context", "db"),
	}
	c.log.Info("connected to database")
	return c, nil
}

// Migrate performs all pending database migrations.
func (c *Connection) Migrate() error {
	c.log.Info("performing migrations...")
	return c.C.AutoMigrate(
		&User{},
		&Account{},
		&Schedule{},
		&Suggestion{},
		&QueueItem{},
		&Tweet{},
		&Vote{},
	).Error
}

// Transaction executes the provided callback in a transaction, automatically
// rolling back the database if an error occurs and committing the changes if
// none occurs.
func (c *Connection) Transaction(fn func(*Connection) error) error {
	d := c.C.Begin()
	if err := fn(&Connection{C: d, log: c.log}); err != nil {
		d.Rollback()
		return err
	}
	d.Commit()
	return nil
}

// Close disconnects from the database.
func (c *Connection) Close() {
	c.C.Close()
}
