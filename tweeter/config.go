package tweeter

import (
	"github.com/ReformedDevs/anonbot/db"
)

type Config struct {
	ConsumerKey    string
	ConsumerSecret string
	ServerURL      string
	Database       *db.Connection
}
