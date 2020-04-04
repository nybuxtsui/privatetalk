package main

import (
	"os"

	"privatetalk/server"

	"github.com/nybuxtsui/log"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "addr",
				Value: ":80",
				Usage: "port",
			},
		},
		Action: server.Run,
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Error("exit: %s", err.Error())
	}
}
