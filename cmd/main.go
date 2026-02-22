package main

import (
	"fmt"
	"os"

	"github.com/inf0rmatiker/logsplitter/pkg/splitter"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ./main <inputfile>")
	}

	// if err := splitter.ProcessSequential(os.Args[1]); err != nil {
	// 	fmt.Fprintln(os.Stderr, err.Error())
	// 	os.Exit(1)
	// }

	if err := splitter.ProcessConcurrent(os.Args[1]); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
