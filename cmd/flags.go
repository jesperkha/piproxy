package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
)

func showHelp() {
	fmt.Println("piproxy [flags]")
	fmt.Println()
	fmt.Println("Flags")
	fmt.Println("    -h --help     Show this help message")
	fmt.Println("    -q --quiet    Disable logging")
	os.Exit(1)
}

func parseFlags() {
	for _, flag := range os.Args[1:] {
		if flag == "--help" || flag == "-h" {
			showHelp()
		} else if flag == "--quiet" || flag == "-q" {
			log.SetOutput(io.Discard)
		} else {
			fmt.Printf("unknown flag '%s'", flag)
			showHelp()
		}
	}
}
