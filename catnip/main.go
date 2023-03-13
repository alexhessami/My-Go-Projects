package main

import (
	"flag"
	"io"
	"os"
	"time"
)

func main() {
	//flag creation
	defaultBufferSize := 4096
	bufferSize := flag.Int("r", defaultBufferSize, "Specify the buffer size. '-r #bytes'")
	flag.Parse()

	//create input and output streams
	input := os.Stdin
	output := os.Stdout

	//create buffer
	buffer := make([]byte, *bufferSize)

	for {
		//read from input stream
		var total int
		for total < *bufferSize {
			data, err := input.Read(buffer[total:*bufferSize])
			if err == io.EOF {
				break
			} else if err != nil {
				panic(err)
			}
			total += data
		}

		//write to output stream
		var err error
		for written := 0; written < total; {
			data, err := output.Write(buffer[written:total])
			if err != nil {
				panic(err)
			}
			written += data
		}
		if err != nil {
			panic(err)
		}

		//wait to limit the rate based on buffer size
		time.Sleep(time.Second / time.Duration(*bufferSize))
	}
}
