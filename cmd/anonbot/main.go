package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ReformedDevs/anonbot/db"
	"github.com/ReformedDevs/anonbot/server"
	"github.com/ReformedDevs/anonbot/tweeter"
	"github.com/howeyc/gopass"
	"github.com/sirupsen/logrus"
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
		cli.BoolFlag{
			Name:   "debug",
			EnvVar: "DEBUG",
			Usage:  "enable debug logging",
		},
		cli.StringFlag{
			Name:   "secret-key",
			EnvVar: "SECRET_KEY",
			Usage:  "secret key",
		},
		cli.StringFlag{
			Name:   "server-addr",
			Value:  ":8000",
			EnvVar: "SERVER_ADDR",
			Usage:  "server driver",
		},
		cli.StringFlag{
			Name:   "server-url",
			EnvVar: "SERVER_URL",
			Usage:  "FQDN of the server",
		},
		cli.StringFlag{
			Name:   "twitter-consumer-key",
			EnvVar: "TWITTER_CONSUMER_KEY",
			Usage:  "Twitter API consumer key",
		},
		cli.StringFlag{
			Name:   "twitter-consumer-secret",
			EnvVar: "TWITTER_CONSUMER_SECRET",
			Usage:  "Twitter API consumer secret ",
		},
	}
	app.Before = func(c *cli.Context) error {
		if c.Bool("debug") {
			logrus.SetLevel(logrus.DebugLevel)
		}
		return nil
	}
	app.Commands = []cli.Command{
		{
			Name:  "createadmin",
			Usage: "create an admin user",
			Action: func(c *cli.Context) error {

				// Create the database connection
				d, err := db.Connect(&db.Config{
					Driver: c.GlobalString("db-driver"),
					Args:   c.GlobalString("db-args"),
				})
				if err != nil {
					return err
				}
				defer d.Close()

				// Prompt for username
				var username string
				fmt.Print("Username? ")
				fmt.Scanln(&username)

				// Prompt for the password, hiding the input
				fmt.Print("Password? ")
				p, err := gopass.GetPasswd()
				if err != nil {
					return err
				}

				// Prompt for email address
				var email string
				fmt.Print("Email? ")
				fmt.Scanln(&email)

				// Create the new user
				u := &db.User{
					Username: username,
					Email:    email,
					IsAdmin:  true,
				}
				if err := u.SetPassword(string(p)); err != nil {
					return err
				}

				// Store the user in the database
				if err := d.C.Create(u).Error; err != nil {
					return err
				}

				return nil
			},
		},
		{
			Name:  "migrate",
			Usage: "perform all database migrations",
			Action: func(c *cli.Context) error {

				// Create the database connection
				d, err := db.Connect(&db.Config{
					Driver: c.GlobalString("db-driver"),
					Args:   c.GlobalString("db-args"),
				})
				if err != nil {
					return err
				}
				defer d.Close()

				// Perform all migrations
				return d.Migrate()
			},
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

		// Create the tweeter
		t := tweeter.New(&tweeter.Config{
			ConsumerKey:    c.String("twitter-consumer-key"),
			ConsumerSecret: c.String("twitter-consumer-secret"),
			ServerURL:      c.String("server-url"),
			Database:       d,
		})
		defer t.Close()

		// Create the server
		s, err := server.New(&server.Config{
			Addr:      c.String("server-addr"),
			SecretKey: c.String("secret-key"),
			Database:  d,
			Tweeter:   t,
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
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Fatal: %s\n", err.Error())
	}
}
