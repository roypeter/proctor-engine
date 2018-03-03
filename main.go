package main

import (
	"os"

	"github.com/gojektech/proctor-engine/server"

	"github.com/urfave/cli"
)

func main() {
	proctor := cli.NewApp()
	proctor.Name = "Proctor"
	proctor.Usage = "Handle orchestration of automated tasks"
	proctor.Commands = []cli.Command{
		{
			Name:    "start",
			Aliases: []string{"s"},
			Usage:   "start server",
			Action: func(c *cli.Context) error {
				return server.Start()
			},
		},
	}

	proctor.Run(os.Args)
}
