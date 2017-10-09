package tweeter

import (
	"github.com/ReformedDevs/anonbot/db"
)

type Config struct {
	ConsumerKey    string
	ConsumerSecret string
	Database       *db.Connection
}
