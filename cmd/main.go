package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	// debugLogPrefix string = "[DEBUG]"
	// infoLogPrefix  string = "[INFO]"
	// warnLogPrefix  string = "[WARN]"
	errorLogPrefix string = "[ERROR]"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ./main <inputfile>")
	}

	if err := splitLogs(os.Args[1]); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

}

func splitLogs(filepath string) error {
	// Open the file and get an *os.File; don't forget to close the file after
	inputFile, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()

	// Open a file to write stdout log lines to
	stdoutFile, err := os.Create(filepath + ".stdout")
	if err != nil {
		panic(err)
	}
	defer stdoutFile.Close()

	// Open a file to write stderr log lines to
	stderrFile, err := os.Create(filepath + ".stderr")
	if err != nil {
		panic(err)
	}
	defer stderrFile.Close()

	var writer io.Writer

	// Use *os.File as an io.Reader to create Scanner
	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {

		// Check for an error on the scanner
		if err = scanner.Err(); err != nil {
			fmt.Printf("Encountered error scanning: %v", err)
			return err
		}

		line := scanner.Text()
		writer = stdoutFile
		if strings.HasPrefix(line, errorLogPrefix) {
			writer = stderrFile
		}
		if err = writeLogLine(writer, []byte((scanner.Text() + "\n"))); err != nil {
			fmt.Printf("Encountered error writing: %v", err)
			return err
		}
	}

	return nil
}

// writeLogLine takes an io.Writer and uses it to write a byte buffer b.
func writeLogLine(writer io.Writer, b []byte) error {
	if _, err := writer.Write(b); err != nil {
		return err
	}
	return nil
}
