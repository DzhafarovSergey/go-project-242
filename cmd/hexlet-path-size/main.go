package main

import (
	"context"
	"fmt"
	"os"

	"code"

	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Name:  "hexlet-path-size",
		Usage: "get size of file or directory",
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
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() == 0 {
				return cli.Exit("Error: path is required", 1)
			}

			path := cmd.Args().First()
			recursive := cmd.Bool("recursive")
			human := cmd.Bool("human")
			all := cmd.Bool("all")

			sizeStr, err := code.GetPathSize(path, recursive, human, all)
			if err != nil {
				return cli.Exit(fmt.Sprintf("Error: %v", err), 1)
			}

			fmt.Printf("%s\t%s\n", sizeStr, path)
			return nil
		},
	}

	app.UsageText = `hexlet-path-size [options] <path>`
	app.Description = `Get size of file or directory.
		If PATH is a directory, shows total size of its contents.
		Examples:
		hexlet-path-size file.txt
		hexlet-path-size -H file.txt
		hexlet-path-size -r -a -H directory/`

	app.UseShortOptionHandling = true

	if err := app.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
