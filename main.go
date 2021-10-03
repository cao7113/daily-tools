package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	app := &cli.App{}
	app.Commands = []*cli.Command{
		{
			Name:    "binary",
			Aliases: []string{"b"},
			Usage:   "binary string",
			Action: func(c *cli.Context) error {
				num := c.Args().First()
				un, err := strconv.ParseUint(num, 10, 64)
				if err != nil {
					return err
				}
				g := 8
				bs := strconv.FormatUint(un, 2)
				l := len(bs)
				r := l % g
				if r > 0 && r < g {
					bs = strings.Repeat("0", g-r) + bs
				}
				fmt.Printf("num: %s\n"+
					"bin(high -> low):\n",
					num)
				chunks := groupChunks(bs, g)
				for i, ch := range chunks {
					fmt.Printf("%0d: %s\n", i+1, ch)
				}
				return nil
			},
		},
		{
			Name:    "stat",
			Aliases: []string{"s"},
			Usage:   "stat a decimal number",
			Action: func(c *cli.Context) error {
				num := c.Args().First()
				chunks := groupChunks(num, 3)
				sNum := strings.Join(chunks, ",")
				fmt.Printf("num: %s\n"+
					"len: %d\n"+
					"sep: %s\n",
					num, len(num), sNum)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func groupChunks(num string, g int) []string {
	var chunks []string
	l := len(num)
	r := l % g
	if r != 0 {
		chunks = append(chunks, num[:r])
	}
	for i := 0; i < l/g; i++ {
		fi := r + i*g
		chunks = append(chunks, num[fi:fi+g])
	}
	return chunks
}
