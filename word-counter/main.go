package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	lines := flag.Bool("l", false, "Count lines")
	help := flag.Bool("help", false, "Default is count words. '-l' will count lines.")
	flag.Parse()

	if *help {
		fmt.Println("Default is count words. '-l' will count lines.")
		return
	}

	fmt.Println(count(os.Stdin, *lines))

	//if !os.Stdin {
	//	fmt.Println("Default is count words. To use tool echo the words and then pipe into this program. To count lines use the '-l' flag. Use the '-help' flag to view these instructions again.")
	//	return
	//}
}

func count(r io.Reader, countLines bool) int {
	scanner := bufio.NewScanner(r)

	if !countLines {
		scanner.Split(bufio.ScanWords)
		fmt.Println("Words:")
	} else {
		fmt.Println("Lines:")
	}

	wc := 0

	for scanner.Scan() {
		wc++
	}

	return wc

}
