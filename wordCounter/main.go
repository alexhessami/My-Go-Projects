package main

import (
	"bufio"
	"bytes" 
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	lines := flag.Bool("l", false, "Count lines")
	help := flag.Bool("help", false, "Default is count words. '-l' will count lines.")
	bytes := flag.Bool("b", false, "Count the number of bytes")
	flag.Parse()

	if *help {
		fmt.Println("Default is count words. '-l' will count lines. '-b' will count bytes.")
		return
	}

	fmt.Println(count(os.Stdin, *lines, *bytes))

}

func count(r io.Reader, countLines bool, countBytes bool) int {
	scanner := bufio.NewScanner(r)

	if countBytes {

		//holds data and reads the contents of r into the reader
		var buf bytes.Buffer
		_, err := io.Copy(&buf, r)

		//error handling
		if err != nil {
			// Handle any errors that occurred while reading
			fmt.Println("Error reading from reader:", err)
		}

		// Convert the buffer contents to a string
		text := buf.String()

		fmt.Println("Bytes:")
		return (len(text))

		//will run by default
	} else if !countLines {
		scanner.Split(bufio.ScanWords)
		fmt.Println("Words:")

		//bufio.Scanner will count lines by default
	} else {
		fmt.Println("Lines:")
	}

	wc := 0

	for scanner.Scan() {
		wc++
	}

	return wc

}
