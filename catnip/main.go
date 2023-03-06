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
	bytesSpecify := flag.Int("r", 0, "Specify the number of bytes. '-r # of bytes'")
	flag.Parse()
	r := *bytesSpecify

	//if -b flag is not used make it 4096 bytes/second by default
	if r == 0 {
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
		data, err := input.Read(buffer)
		//error handling
		//ignoring EOF error since we are continuously reading the file
		if err != nil && err != io.EOF {
			panic(err)
		}
		if data == 0 {
			//Wait for one second if no input before continuing
			time.Sleep(time.Second)
			break
		}

		//Write to the output stream
		_, err = output.Write(buffer[:data])
		//error handling
		if err != nil {
			panic(err)
		}

		//Wait to limit the rate based off of specified bytes (default 4096)
		time.Sleep(time.Second / time.Duration(bytes))
	}
}
