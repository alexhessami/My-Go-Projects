package main

import (
	"io"
	"os"
	"time"
)

func main() {
	//Created the input and output streams
	in := os.Stdin
	out := os.Stdout

	//Allocate 4096 bytes to read from the input stream
	buf := make([]byte, 4096)

	//Create loop to continuously output data as it comes in
	for {
		//Read from the input stream
		n, err := in.Read(buf)
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
		_, err = out.Write(buf[:n])
		//error handling
		if err != nil {
			panic(err)
		}

		//Wait to limit the rate
		time.Sleep(time.Second / 4096)
	}
}
