package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

const cmdName = "daily-tools"

func main() {
	app := &cli.App{
		Name:     cmdName,
		HelpName: cmdName,
		Usage:    `daily tools command`,
	}
	app.Commands = []*cli.Command{
		numCmd(),
		timeCmd(),
		logicCmd(),
		asciiCmd(),
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
