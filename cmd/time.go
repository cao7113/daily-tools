package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"strconv"
	"strings"
	"time"
)

func timeCmd() *cli.Command {
	return &cli.Command{
		Name:    "time",
		Aliases: []string{"t"},
		Usage:   "time related commands",
		Action: func(c *cli.Context) error {
			return c.App.Command("now").Run(c)
		},
		Subcommands: []*cli.Command{
			{
				Name:    "now",
				Aliases: []string{"n"},
				Usage:   "now time",
				Action: func(c *cli.Context) error {
					n := time.Now().UTC()
					nl := n.Local()
					fmt.Printf("local: %s\n", nl)
					fmt.Printf("utc  : %s\n", n)
					fmt.Printf("epoch: %d\n", n.Unix())
					fmt.Printf("nano : %d\n", n.UnixNano())
					wd := nl.Weekday()
					fmt.Printf("wkday: %s\n", wd.String())
					fmt.Printf("yrday: %d\n", nl.YearDay())
					return nil
				},
			}, // now
			{
				Name:    "timestamp",
				Aliases: []string{"at"},
				Usage:   "time on the timestamp",
				Action: func(c *cli.Context) error {
					inTs := c.Args().First()
					l := len(inTs)
					secLen := 10 // normal timestamp length in seconds
					max := 19    // max length in nano-seconds
					if l > max {
						return fmt.Errorf("invalid time: %s", inTs)
					}
					trailingZ := 0
					if l < secLen {
						trailingZ = max - secLen
					} else {
						trailingZ = max - l
					}

					fs := inTs + strings.Repeat("0", trailingZ)
					fs = strings.Repeat("0", max-len(fs)) + fs
					ns, err := strconv.ParseInt(fs, 10, 64)
					if err != nil {
						return err
					}
					tm := time.Unix(ns/1e9, ns%1e9)
					fmt.Printf("timestamp : %s\n", inTs)
					fmt.Printf("raw length: %d\n", l)
					fmt.Printf("local time: %s\n", tm.Local())
					fmt.Printf("utc   time: %s\n", tm.UTC())
					fmt.Printf("since du  : %s\n", time.Since(tm).String())
					return nil
				},
			}, // time
		},
	}
}
