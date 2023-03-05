package main

import (
	"flag"
	"io"
	"os"
	"time"
)

func main() {

	//flag creation
	bytes := 4096
	var x int
	bytesSpecify := flag.Int("b", x, "Specify the number of bytes. '-b # of bytes'")
	flag.Parse()

	//if -b flag is not used make it 4096 bytes/second by default
	if x == 0 {
		bytes = 4096
	} else {
		bytes = *bytesSpecify
	}

	//Created the input and output streams
	input := os.Stdin
	output := os.Stdout

	//Allocate specified bytes to read from the input stream (4096 by default)
	buffer := make([]byte, bytes)

	//Create loop to continuously output data as it comes in
	for {
		//Read from the input stream
		n, err := input.Read(buffer)
		//error handling
		//ignoring EOF error since we are continuously reading the file
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			//Wait for one second if no input before continuing
			time.Sleep(time.Second)
			continue
		}

		//Write to the output stream
		_, err = output.Write(buffer[:n])
		//error handling
		if err != nil {
			panic(err)
		}

		//Wait to limit the rate
		time.Sleep(time.Second / 4096)
	}
}
