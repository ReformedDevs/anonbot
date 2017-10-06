package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ReformedDevs/anonbot/db"
	"github.com/ReformedDevs/anonbot/server"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "anonbot"
	app.Usage = "Twitter anon account bot"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "db-args",
			Value:  "dbname=postgres user=postgres",
			EnvVar: "DB_ARGS",
			Usage:  "database arguments",
		},
		cli.StringFlag{
			Name:   "db-driver",
			Value:  "postgres",
			EnvVar: "DB_DRIVER",
			Usage:  "database driver",
		},
		cli.StringFlag{
			Name:   "server-addr",
			Value:  ":8000",
			EnvVar: "SERVER_ADDR",
			Usage:  "server driver",
		},
	}
	app.Action = func(c *cli.Context) error {

		// Create the database connection
		d, err := db.Connect(&db.Config{
			Driver: c.String("db-driver"),
			Args:   c.String("db-args"),
		})
		if err != nil {
			return err
		}
		defer d.Close()

		// Create the server
		s, err := server.New(&server.Config{
			Addr:     c.String("server-addr"),
			Database: d,
		})
		if err != nil {
			return err
		}
		defer s.Close()

		// Wait for SIGINT or SIGTERM
		sigChan := make(chan os.Signal)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		return nil
	}
	app.Run(os.Args)
}
