package main

import (
	"flag"
	"fmt"
	"os"

	"code"
)

func main() {
	recursive := flag.Bool("r", false, "recursive size of directories")
	recursiveLong := flag.Bool("recursive", false, "recursive size of directories")
	human := flag.Bool("H", false, "human-readable sizes (auto-select unit)")
	humanLong := flag.Bool("human", false, "human-readable sizes (auto-select unit)")
	all := flag.Bool("a", false, "include hidden files and directories")
	allLong := flag.Bool("all", false, "include hidden files and directories")
	help := flag.Bool("h", false, "show help")
	helpLong := flag.Bool("help", false, "show help")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <path>\n\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Options:")
		fmt.Fprintln(os.Stderr, "  -r, --recursive    recursive size of directories")
		fmt.Fprintln(os.Stderr, "  -H, --human        human-readable sizes (auto-select unit)")
		fmt.Fprintln(os.Stderr, "  -a, --all          include hidden files and directories")
		fmt.Fprintln(os.Stderr, "  -h, --help         show help")
		fmt.Fprintln(os.Stderr, "\nExamples:")
		fmt.Fprintln(os.Stderr, "  hexlet-path-size file.txt")
		fmt.Fprintln(os.Stderr, "  hexlet-path-size -H file.txt")
		fmt.Fprintln(os.Stderr, "  hexlet-path-size -r -a -H directory/")
	}

	flag.Parse()

	if *help || *helpLong {
		flag.Usage()
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Error: path is required")
		flag.Usage()
		os.Exit(1)
	}

	path := args[0]

	isRecursive := *recursive || *recursiveLong
	isHuman := *human || *humanLong
	isAll := *all || *allLong

	result, err := code.GetPathSizeWithPath(path, isRecursive, isHuman, isAll)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(result)
}
