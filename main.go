package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/datianshi/cfnetwork/curlv2"
)

func main() {
	app := cli.NewApp()
	app.Name = "cfnetwork"
	app.Usage = "Trouble Shooting cf network issues"

	app.Commands = []cli.Command{
		{
			Name:    "curlv2",
			Aliases: []string{"c"},
			Usage:   "curl",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "domain",
					Usage: "domain for cloudfoundry",
				},
				cli.BoolFlag{
					Name:  "https",
					Usage: "using https",
				},
			},
			Action: curlv2.CliDomainAction,
			Subcommands: []cli.Command{
				{
					Name:  "router",
					Usage: "curl router directly",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "domain",
							Usage: "cloudfoundry domain",
						},
						cli.StringFlag{
							Name:  "ip",
							Usage: "router ip",
						},
					},
					Action: curlv2.CliRouterAction,
				},
			},
		},
	}

	app.Run(os.Args)
}
