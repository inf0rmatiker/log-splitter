package splitter

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

const (
	errorLogSignature string = "ERR"

	stdoutFileSuffix string = ".stdout"
	stderrFileSuffix string = ".stderr"
)

// ProcessSequential processes log lines in an input file specified by filepath.
// If log lines contain an error prefix, they are written to a .stderr file, otherwise,
// they are written to a .stdout file. Line writing happens sequentially, all on one thread.
func ProcessSequential(filepath string) error {

	var inputFile, stderrFile, stdoutFile *os.File
	var err error

	// Open the input file for reading
	if inputFile, err = os.Open(filepath); err != nil {
		return err
	}
	defer inputFile.Close()

	// Open a file to write stdout log lines to
	if stdoutFile, err = os.Create(filepath + stdoutFileSuffix); err != nil {
		return err
	}
	defer stdoutFile.Close()

	// Open a file to write stderr log lines to
	if stderrFile, err = os.Create(filepath + stderrFileSuffix); err != nil {
		return err
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
		writer = stdoutFile // use *os.File as io.Writer
		if strings.Contains(line, errorLogSignature) {
			writer = stderrFile
		}
		if err = writeLogLine(writer, line); err != nil {
			fmt.Printf("Encountered error writing line: %v", err)
			return err
		}
	}

	return nil
}

// writeLogLine takes an io.Writer and uses it to write a string s as its own line.
func writeLogLine(writer io.Writer, s string) error {
	if _, err := writer.Write([]byte(s + "\n")); err != nil {
		return err
	}
	return nil
}

// ProcessConcurrent processes log lines in an input file specified by filepath.
// If log lines contain an error prefix, they are written to a .stderr file, otherwise,
// they are written to a .stdout file. Line writing happens concurrently with multiple goroutines.
func ProcessConcurrent(filepath string) error {

	var inputFile, stderrFile, stdoutFile *os.File
	var err error

	// Open input file for reading
	if inputFile, err = os.Open(filepath); err != nil {
		return err
	}
	defer inputFile.Close()

	// Open a file to write stdout log lines to
	if stdoutFile, err = os.Create(filepath + stdoutFileSuffix); err != nil {
		return err
	}
	defer stdoutFile.Close()

	// Open a file to write stderr log lines to
	if stderrFile, err = os.Create(filepath + stderrFileSuffix); err != nil {
		return err
	}
	defer stderrFile.Close()

	// Start two goroutines:
	//   1 for writing to stdout (should have its own channel)
	//   2 for writing to stderr (should have its own channel)
	// We'll need to 'wait' on these to finish writing, otherwise the main
	// thread could exit before the children have had time to write their lines.
	// So, we'll use a sync.WaitGroup for 2 goroutines.
	stderrChan, stdoutChan := make(chan string), make(chan string)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go processLogs(stdoutFile, stdoutChan, &wg)
	go processLogs(stderrFile, stderrChan, &wg)

	// Iterate over input file scanned lines
	//   Check line string prefix
	// 		Output line on appropriate channel
	// Use *os.File as an io.Reader to create Scanner
	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {

		// Check for an error on the scanner
		if err = scanner.Err(); err != nil {
			fmt.Printf("Encountered error scanning: %v", err)
			close(stdoutChan)
			close(stderrChan)
			return err
		}

		line := scanner.Text()
		if strings.Contains(line, errorLogSignature) {
			stderrChan <- line
		} else {
			stdoutChan <- line
		}
	}

	// Close channels so goroutines finish their loops, and mark wg.Done().
	close(stdoutChan)
	close(stderrChan)

	// Wait on all goroutines to finish what they were doing before exiting.
	wg.Wait()

	return nil
}

func processLogs(writer io.Writer, ch chan string, wg *sync.WaitGroup) {
	// Consume strings off the channel until it's closed,
	// writing them to the passed-in io.Writer. After the channel is closed,
	// mark the WaitGroup done for our goroutine.
	defer wg.Done()
	for s := range ch {
		writeLogLine(writer, s)
	}
}
