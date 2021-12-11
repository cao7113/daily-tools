package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"unicode"
)

func asciiCmd() *cli.Command {
	return &cli.Command{
		Name:    "ascii",
		Aliases: []string{"a"},
		Usage:   "ascii-code table",
		Action: func(c *cli.Context) error {
			fmt.Printf("# ref https://zh.wikipedia.org/wiki/ASCII \n")
			for i := 0; i < unicode.MaxLatin1; i++ {
				rStart := i%0x10 == 0
				if rStart {
					fmt.Printf("%02x: ", i)
				}

				c := rune(i)
				if unicode.IsGraphic(c) {
					fmt.Printf("%c	", c)
				} else {
					fmt.Printf("%02x	", i)
				}

				if (i+1)%0x10 == 0 {
					fmt.Printf("\n")
				}
			}
			fmt.Printf("\n")
			return nil
		},
	}
}
