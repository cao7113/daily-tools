package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

var boolOperand = map[string]uint8{"1": 1, "0": 0}

func logicCmd() *cli.Command {
	return &cli.Command{
		Name:    "logic",
		Aliases: []string{"l"},
		Usage:   "logic related commands",
		Action: func(c *cli.Context) error {
			return c.App.Command("table").Run(c)
		},
		Subcommands: []*cli.Command{
			{
				Name:        "table",
				Aliases:     []string{"tab"},
				Description: "operation table",
				Action: func(c *cli.Context) error {
					fmt.Printf("AND operation(true if all true)\n")
					for xk, xv := range boolOperand {
						for yk, yv := range boolOperand {
							fmt.Printf("%s & %s = %d\n", xk, yk, xv&yv)
						}
					}
					fmt.Printf("OR operation(true if any true)\n")
					for xk, xv := range boolOperand {
						for yk, yv := range boolOperand {
							fmt.Printf("%s | %s = %d\n", xk, yk, xv|yv)
						}
					}
					fmt.Printf("XOR(Exclusive OR) operation(true if different)\n")
					for xk, xv := range boolOperand {
						for yk, yv := range boolOperand {
							fmt.Printf("%s ^ %s = %d\n", xk, yk, xv^yv)
						}
					}
					fmt.Printf("AND NOT(Bit Clear) operation(true if not clear)\n")
					for xk, xv := range boolOperand {
						for yk, yv := range boolOperand {
							fmt.Printf("%s &^ %s = %d\n", xk, yk, xv&^yv)
						}
					}

					return nil
				},
			},
		},
	}
}
