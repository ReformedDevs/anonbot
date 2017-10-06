package server

import (
	"github.com/ReformedDevs/anonbot/db"
	"github.com/ReformedDevs/anonbot/tweeter"
)

// Config provides parameters for hosting the website.
type Config struct {
	Addr      string
	SecretKey string
	Database  *db.Connection
	Tweeter   *tweeter.Tweeter
}
