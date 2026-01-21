package main

import (
	"flag"
	"fmt"
	"os"

	"code"
)

func main() {
	Run()
}

func Run() {
	humanFlag := flag.Bool("human", false, "human-readable sizes (auto-select unit)")
	shortHumanFlag := flag.Bool("h", false, "human-readable sizes (shorthand)")
	allFlag := flag.Bool("all", false, "include hidden files and directories")
	shortAllFlag := flag.Bool("a", false, "include hidden files and directories (shorthand)")
	recursiveFlag := flag.Bool("recursive", false, "recursive size of directories")
	shortRecursiveFlag := flag.Bool("r", false, "recursive size of directories (shorthand)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <path>\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Global options:")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\nExamples:")
		fmt.Fprintln(os.Stderr, "  hexlet-path-size data.csv")
		fmt.Fprintln(os.Stderr, "  hexlet-path-size --human data.csv")
		fmt.Fprintln(os.Stderr, "  hexlet-path-size -h data.csv")
		fmt.Fprintln(os.Stderr, "  hexlet-path-size --all project/")
		fmt.Fprintln(os.Stderr, "  hexlet-path-size -a project/")
		fmt.Fprintln(os.Stderr, "  hexlet-path-size -h -a project/")
		fmt.Fprintln(os.Stderr, "  hexlet-path-size -h -a -r project/")
		fmt.Fprintln(os.Stderr, "\nNote: Without -r flag, only top-level files are counted for directories.")
	}

	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	path := flag.Arg(0)

	human := *humanFlag || *shortHumanFlag
	all := *allFlag || *shortAllFlag
	recursive := *recursiveFlag || *shortRecursiveFlag

	result, err := code.GetPathSize(path, recursive, human, all)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(result)
}
