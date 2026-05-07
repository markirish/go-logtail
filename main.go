package main

import (
	"bufio"
	"flag"
	"io"
	"os"
	"strings"
	"time"

	"github.com/markirish/go-logtail/internal/parser"
)

type LogEntry struct {
	File string
	Text string
}

func main() {
	// Set up command-line flags
	debug := flag.Bool("debug", false, "enable debug logging")

	flag.Parse()
	fileNames := flag.Args()

	fileNotFound := false

	// Create the channel early, so that we only have to iterate through the
	// files once.
	logEntries := make(chan LogEntry)

	for _, fileName := range fileNames {
		if *debug {
			parser.PrintLog("[DEBUG]: Checking file: %s\n", fileName)
		}

		fileInfo, err := os.Stat(fileName)
		if err != nil {
			fileNotFound = true
			if os.IsNotExist(err) {
				if *debug {
					parser.PrintLog("[DEBUG]: %s not found\n", fileName)
				}
			} else {
				if *debug {
					parser.PrintLog("[DEBUG]: could not stat %s: %v\n", fileName, err)
				}
			}
			continue
		} else {
			if *debug {
				parser.PrintLog("[DEBUG]: %s found\n", fileName)
			}
			// The file exists, so we can start tailing it.
			file, err := os.Open(fileName)
			if err != nil {
				logEntries <- LogEntry{File: fileInfo.Name(), Text: parser.GenerateLog("error opening file: %v\n", err)}
				continue
			}
			go tailFile(file, logEntries)
		}
	}

	if fileNotFound {
		parser.PrintLog("[ERROR]: one or more files were not found\n")
		os.Exit(1)
	}

	// Print log entries as they come in.
	for entry := range logEntries {
		parser.PrintLog("[%s]: %s\n", entry.File, entry.Text)
	}
}

func tailFile(file *os.File, logEntries chan<- LogEntry) {
	defer file.Close()

	_, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		logEntries <- LogEntry{File: file.Name(), Text: parser.GenerateLog("error seeking file: %v\n", err)}
		return
	}

	reader := bufio.NewReader(file)

	// In the case that a line gets split across multiple reads or writes, we
	// need to keep track of the partial line read until a new line
	// Use a strings.Builder to efficiently build the line as it is read.
	var partial strings.Builder

	for {
		// Use ReadString to read until the next newline character, or the EOF if
		// there is no newline. This way we can handle lines that are written in
		// multiple parts.
		line, err := reader.ReadString('\n')

		// Encountered a newline, so we have a complete line to process.
		if err == nil {
			partial.WriteString(line)
			text := strings.TrimRight(partial.String(), "\n")

			if text != "" {
				logEntries <- LogEntry{File: file.Name(), Text: text}
			}

			// Reset partial to potnetiall grab the next line if split across writes.
			partial.Reset()

			continue
		}

		// EOF is treated as a normal condition, even though passed as an error. It
		// can happen with multiple writes to the file without a newline. Grab the
		// partial line and store it until the next write, which may complete the
		// line.
		if err == io.EOF {
			if line != "" {
				partial.WriteString(line)
			}

			time.Sleep(500 * time.Millisecond)
			continue
		}

		// Encountered an error that wasn't EOF! Log it and exit the goroutine.
		logEntries <- LogEntry{
			File: file.Name(),
			Text: parser.GenerateLog("error reading file: %v", err),
		}
		return
	}
}
