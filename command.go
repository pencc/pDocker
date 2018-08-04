package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	"pDocker/cgroups"
	"pDocker/container"
)

var initCommand = cli.Command{
	Name:  "init",
	Usage: "Init container process run user's process in container. Do not call it outside",
	Action: func(context *cli.Context) error {
		log.Info("in init")
		container.RunContainerInitProcess(context.Args().Get(0), nil)
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
		cli.StringFlag{
			Name:  "m",
			Usage: "memory limit",
		},
		cli.StringFlag{
			Name:  "cpuset",
			Usage: "cpu number",
		},
		cli.StringFlag{
			Name:  "cpushare",
			Usage: "cpu share",
		},
	},
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("Missing container command")
		}

		rc := cgroups.ResourceConfig{
			MemoryLimit: context.String("m"),
			CpuShare:    context.String("cpushare"),
			CpuSet:      context.String("cpuset"),
		}

		log.Info("cmd:" + context.Args().Get(0))
		container.Run(context.Bool("ti"), context.Args().Get(0), rc)
		return nil
	},
}
