package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

var initCommand = cli.Command{
	Name:  "init",
	Usage: "Init container process run user's process in container. Do not call it outside",
	Action: func(c *cli.Context) error {
		log.Info("in init")
		return nil
	},
}

var runCommand = cli.Command{
	Name:  "run",
	Usage: "Create a container with namespace and cgroups limit",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "ti",
			Usage: "enable tty",
		},
	},
	Action: func(c *cli.Context) error {
		log.Info("in commit")
		return nil
	},
}
