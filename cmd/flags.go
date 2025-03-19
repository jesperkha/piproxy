package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
)

type Flags struct {
	LogFile bool
}

func showHelp() {
	fmt.Println("piproxy [flags]")
	fmt.Println()
	fmt.Println("Flags")
	fmt.Println("    -h --help     Show this help message")
	fmt.Println("    -q --quiet    Disable logging")
	fmt.Println("    --logfile     Log to file (LOG_FILE in .env)")
	os.Exit(1)
}

func parseFlags() (flags Flags) {
	for _, flag := range os.Args[1:] {
		if flag == "--help" || flag == "-h" {
			showHelp()
		} else if flag == "--quiet" || flag == "-q" {
			log.SetOutput(io.Discard)
		} else if flag == "--logfile" {
			flags.LogFile = true
		} else {
			fmt.Printf("unknown flag '%s'\n", flag)
			showHelp()
		}
	}

	return flags
}
