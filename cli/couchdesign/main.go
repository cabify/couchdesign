package main

import (
	"os"

	"github.com/apex/log"
	"github.com/cabify/couchdesign"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "CouchDB design document admin tool"
	app.UsageText = "$ couchdesign [COMMAND] [OPTIONS]"
	app.Version = "0.0.1"
	app.Authors = []cli.Author{{Name: "Carlos Alonso", Email: "carlos.alonso@cabify.com"}}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "server",
			Usage: "Server to connect to",
			Value: "127.0.0.1",
		},
		cli.StringFlag{
			Name:  "username",
			Usage: "Username to authenticate",
		},
		cli.StringFlag{
			Name:  "password",
			Usage: "Password for the username to authenticate",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "status",
			Usage: "Describe the current status of design docs",
			Action: func(c *cli.Context) {
				serverUrl := c.GlobalString("server")
				log.WithField("server", serverUrl).Info("Retrieving design docs statses...")
				authHttpReq := buildAuthHttpReq(c)
			},
		},
	}

	app.Run(os.Args)
}

func buildAuthHttpReq(c *cli.Context) *couchdesign.AuthHttpRequester {
	ahr, err := couchdesign.NewAuthHttpRequester(c.GlobalString("admin"), c.GlobalString("password"), c.GlobalString("server"))
	if err != nil {
		log.WithError(err).
			WithFields(log.Fields{"admin": c.GlobalString("admin"), "password": c.GlobalString("password"), "server": c.GlobalString("server")}).
			Fatal("Couldn't create the AuthHttpRequester")
	}
	return ahr
}
