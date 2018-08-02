package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
)

const usage = `pdocker is a simple container runtime implementation.
			   The purpose of this project is to learn how docker works. 
			   example: 
			   # sudo ./pDocker run -ti /bin/sh, 
			   you may need run 
			   # sudo mount -t proc proc /proc
			   before launch pDocker`

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		initCommand,
		runCommand,
	}
	app.Usage = usage
	app.Before = func(context *cli.Context) error {
		log.SetFormatter(&log.JSONFormatter{})
		log.SetOutput(os.Stdout)
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
