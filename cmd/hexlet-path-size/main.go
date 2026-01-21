package main

import (
	"fmt"
	"log"
	"os"

	"code"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "hexlet-path-size",
		Usage: "print size of a file or directory; supports -r (recursive), -H (human-readable), -a (include hidden)",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "recursive",
				Aliases: []string{"r"},
				Usage:   "recursive size of directories",
			},
			&cli.BoolFlag{
				Name:    "human",
				Aliases: []string{"H"},
				Usage:   "human-readable sizes (auto-select unit)",
			},
			&cli.BoolFlag{
				Name:    "all",
				Aliases: []string{"a"},
				Usage:   "include hidden files and directories",
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() == 0 {
				return cli.ShowAppHelp(c)
			}

			path := c.Args().First()
			result, err := code.GetPathSize(
				path,
				c.Bool("recursive"),
				c.Bool("human"),
				c.Bool("all"),
			)
			if err != nil {
				return err
			}

			fmt.Println(result)
			return nil
		},
		EnableBashCompletion: true,
	}

	app.CommandNotFound = func(c *cli.Context, command string) {
		fmt.Fprintf(c.App.Writer, "No command '%s' found.\n", command)
		cli.ShowAppHelp(c)
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
