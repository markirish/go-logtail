package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	ANSIColorRed   = "\033[31m"
	ANSIColorGreen = "\033[32m"
	ANSIColorReset = "\033[0m"
)

func main() {
	// Set up command-line flags
	debug := flag.Bool("debug", false, "enable debug logging")

	flag.Parse()
	files := flag.Args()

	for _, file := range files {
		if *debug {
			fmt.Println("DEBUG: Checking file:", file)
		}

		_, err := os.Stat(file)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("%sERROR: file not found: %s%s\n", ANSIColorRed, file, ANSIColorReset)
				continue
			}

			fmt.Printf("%sERROR: could not stat %s: %w%s\n", ANSIColorRed, file, err, ANSIColorReset)
			continue
		} else {
			if *debug {
				fmt.Printf("DEBUG: %s found\n", file)
			}
		}
	}

}
