package main
//test123
import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

func main() {
	//flag creation
	defaultBufferSize := 4096
	bufferSize := flag.Int("r", defaultBufferSize, "Specify the buffer size. '-r #bytes'")
	verbose := flag.Bool("v", false, "Enable verbose mode.")
	help := flag.Bool("help", false, "By default output will be 4096 bytes a second. Use -r # to specify the number of bytes. Use -v for verbose.")
	flag.Parse()

	//create input and output streams
	input := os.Stdin
	output := os.Stdout

	if err := start(input, output, *bufferSize, *verbose, *help); err != nil {
		panic(err)
	}
}

func start(input io.Reader, output io.Writer, bufferSize int, verbose bool, help bool) error {
	//create buffer
	buffer := make([]byte, bufferSize)

	for {
		if help {
			fmt.Println("By default output will be 4096 bytes a second. Use -r # to specify the number of bytes. Use -v for verbose.")
			break
		}

		//read from input stream
		var total int
		for total < bufferSize {
			data, err := input.Read(buffer[total:bufferSize])
			if err == io.EOF {
				break
			} else if err != nil {
				return err
			}
			total += data

			if verbose {
				fmt.Print(".")
			}
		}

		//write to output stream
		for written := 0; written < total; {
			var err error
			data, err := output.Write(buffer[written:total])
			if err != nil {
				return err
			}
			written += data
		}

		//wait to limit the rate based on buffer size
		time.Sleep(time.Second / time.Duration(bufferSize))
	}
	return nil

}
