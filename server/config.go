package server

import (
	"github.com/ReformedDevs/anonbot/db"
)

// Config provides parameters for hosting the website.
type Config struct {
	Addr     string
	Database *db.Connection
}
