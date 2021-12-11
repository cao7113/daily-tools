package main

import (
	"fmt"
	"github.com/cao7113/dailytools/lib"
	"github.com/urfave/cli/v2"
	"strconv"
	"strings"
	"unicode/utf8"
)

func numCmd() *cli.Command {
	return &cli.Command{
		Name:    "num",
		Aliases: []string{"n"},
		Usage:   "number related commands",
		Action: func(c *cli.Context) error {
			return c.App.Command("han").Run(c)
		},
		Subcommands: []*cli.Command{
			{
				Name:    "binary",
				Aliases: []string{"b", "int"},
				Usage:   "binary info for signed number(two's complement representation) e.g. num binary -- -2",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:        "base",
						Aliases:     []string{"b"},
						Value:       10,
						DefaultText: "10",
					},
					&cli.BoolFlag{
						Name:        "big-endian",
						Aliases:     []string{"e"},
						Value:       true,
						DefaultText: "true",
					},
					&cli.BoolFlag{
						Name:        "min-bytes",
						Aliases:     []string{"m"},
						Value:       true,
						DefaultText: "true",
					},
				},
				Action: func(c *cli.Context) error {
					base := c.Int("base")
					be := c.Bool("big-endian")
					mb := c.Bool("min-bytes")
					num := c.Args().First()
					if base == 10 {
						if len(num) >= 2 {
							c2 := num[:2]
							switch c2 {
							case "0b":
								num = num[2:]
								base = 2
							case "0o":
								base = 8
								num = num[2:]
							case "0x":
								base = 16
								num = num[2:]
							}
						}
					}
					if num == "" {
						return fmt.Errorf("require num")
					}

					n, err := strconv.ParseInt(num, base, 64)
					if err != nil {
						return err
					}

					ed := "little-Endian"
					if be {
						ed = "big-Endian"
					}

					str := fmt.Sprintf("num: %s (%d-based)", num, base)
					if base != 10 {
						str += fmt.Sprintf(" %d(10-based)", n)
					}
					fmt.Println(str)
					bs := lib.GetIntBinString(n, be, !mb, " ")
					fmt.Printf("%s (two's complement with %s)\n", bs, ed)

					// for negative two's complement
					if n < 0 {
						n = -n // same as n := ^(n - 0x01)
						fmt.Printf("\n## other-side num: %d\n", n)
						bs = lib.GetIntBinString(n, be, !mb, " ")
						fmt.Printf("bin(big-endian: %t): \n%s\n", be, bs)

						n1 := ^n
						bs = lib.GetIntBinString(n1, be, !mb, " ")
						fmt.Printf("bin(one's complement with big-endian: %t): \n%s\n", be, bs)

						n2 := n1 + 0x01
						bs = lib.GetIntBinString(n2, be, !mb, " ")
						fmt.Printf("bin(two's complement with big-endian: %t): \n%s\n", be, bs)
					}

					return nil
				},
			}, // binary
			{
				Name:    "unsigned-integer",
				Aliases: []string{"ub"},
				Usage:   "binary representation of unsigned number",
				Action: func(c *cli.Context) error {
					num := c.Args().First()
					un, err := strconv.ParseUint(num, 10, 64)
					if err != nil {
						return err
					}
					bs := strconv.FormatUint(un, 2)
					chunks := lib.SplitToGroupChunks(bs, 8, '0')
					hs := strconv.FormatUint(un, 0x10)
					hChunks := lib.SplitToGroupChunks(hs, 2, '0')

					fmt.Printf("num: %s\n", num)
					fmt.Printf("hex: %s\n", strings.Join(hChunks, ""))
					fmt.Printf("bin: %s\n", strings.Join(chunks, " "))
					fmt.Printf("bytes in detail:\n")
					for i, ch := range chunks {
						fmt.Printf("%0d: %s 0x%s\n", i+1, ch, hChunks[i])
					}
					return nil
				},
			}, // unsigned-integer
			{
				Name:    "han [char]",
				Aliases: []string{"zh"},
				Usage:   "chinese char table",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:        "rows",
						Aliases:     []string{"r"},
						Value:       20,
						DefaultText: "20",
					},
					&cli.StringFlag{
						Name:        "from",
						Aliases:     []string{"f"},
						DefaultText: "0x4e00",
					},
				},
				Action: func(c *cli.Context) error {
					fromCp := c.String("from")
					rows := c.Int("rows")

					fromChar := c.Args().First()
					if fromChar != "" {
						if fromCp != "" {
							return fmt.Errorf("--from conflict with arg, use arg")
						}
						r, _ := utf8.DecodeRune([]byte(fromChar))
						fromCp = strconv.FormatInt(int64(r), 10)
					}
					if fromCp == "" {
						fromCp = "0x4e00"
					}
					limit := 0x10 * rows
					base := 10
					if strings.HasPrefix(fromCp, "0x") {
						base = 0x10
						fromCp = strings.TrimLeft(fromCp, "0x")
					}
					un, err := strconv.ParseUint(fromCp, base, 64)
					if err != nil {
						return err
					}
					fmt.Printf("# %d chars from code-point: %s %d(10-based)\n", limit, fromCp, un)
					fmt.Printf("0x000: ")
					for i := 0; i < 0x10; i++ {
						fmt.Printf("0x%02x	", i)
					}
					fmt.Printf("\n")
					for i := 0; i < limit; i++ {
						r := rune(un + uint64(i))
						if i%0x10 == 0 {
							fmt.Printf("0x%03x: ", r)
						}
						fmt.Printf("%c	", r)
						if (i+1)%0x10 == 0 {
							fmt.Printf("\n")
						}
					}

					return nil
				},
			}, // han
			{
				Name:    "unicode-point",
				Aliases: []string{"p"},
				Usage:   "unicode code-point",
				Action: func(c *cli.Context) error {
					inNum := c.Args().First()
					base := 10
					num := inNum
					if strings.HasPrefix(inNum, "0x") {
						base = 0x10
						num = strings.TrimLeft(inNum, "0x")
					}
					un, err := strconv.ParseUint(num, base, 64)
					if err != nil {
						return err
					}
					fmt.Printf("%s: %d 0x%04x 0b%0b %c %b\n",
						inNum, un, un, un, rune(un), []byte(string(rune(un))))
					return nil
				},
			}, // unicode-point
			{
				Name:    "unicode",
				Aliases: []string{"u"},
				Usage:   "unicode encoding",
				Action: func(c *cli.Context) error {
					cs := c.Args().First()
					bs := []byte(cs)
					for len(bs) > 0 {
						r, size := utf8.DecodeRune(bs)
						fmt.Printf("coding: %c	%d	0x%0x	%d bytes	", r, r, r, size)
						fmt.Printf("hex: ")
						for _, b := range bs[:size] {
							fmt.Printf("x%0x ", b)
						}
						fmt.Printf("	bin: ")
						for _, b := range bs[:size] {
							fmt.Printf("%08b ", b)
						}
						fmt.Printf("\n")
						bs = bs[size:]
					}
					fmt.Printf("string: %s\n", cs)
					fmt.Printf("c-len : %d\n", utf8.RuneCount([]byte(cs)))
					fmt.Printf("b-len : %d\n", len([]byte(cs)))

					return nil
				},
			}, // unicode
			{
				Name:    "stat",
				Aliases: []string{"s"},
				Usage:   "stat a decimal number",
				Action: func(c *cli.Context) error {
					num := c.Args().First()
					chunks := lib.SplitToGroupChunks(num, 3, 0)
					sNum := strings.Join(chunks, ",")
					fmt.Printf("num: %s\n"+
						"len: %d\n"+
						"sep: %s\n",
						num, len(num), sNum)

					if strings.HasPrefix(num, "00") {
						fmt.Printf("\ntrim headling zeros\n")
						vNum := strings.TrimLeftFunc(num, func(r rune) bool {
							return r == '0'
						})
						vChunks := lib.SplitToGroupChunks(vNum, 3, 0)
						vsNum := strings.Join(vChunks, ",")
						fmt.Printf("num: %s\n"+
							"len: %d\n"+
							"sep: %s\n",
							vNum, len(vNum), vsNum)
					}
					return nil
				},
			}, // stat
		},
	}
}
