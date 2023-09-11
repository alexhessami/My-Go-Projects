package main

import (
    "fmt"
    "io/ioutil"
    "os"
    "os/signal"
    "syscall"
    "time"
    "strings"
    "strconv"
)

const (
    blkioPath = "/sys/fs/cgroup/blkio/"
    serviceName = "blkio.throttle.io_service_bytes"
    serviceNewName = "blkio.throttle.io_serviced"
    logInterval = 5 * time.Second
)

func main() {
    // Create a channel to handle Ctrl+C signals.
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    // Create a ticker to collect data at regular intervals.
    ticker := time.NewTicker(logInterval)

    currentbytesInt := 0
    currentrwInt := 0

    // Start the data collection loop.
    for {
        select {
        case <-ticker.C:
            timestamp := time.Now().Format(time.RFC3339)
		
	    totalReadWrite, err := readBlkioReadWrite()
	    slicedRW := totalReadWrite[6:]
	    newrwInt, _ := strconv.Atoi(slicedRW)
            if err != nil {
                fmt.Printf("[%s] Error reading blkio read writes: %v\n", timestamp, err)
	}
		
            totalBytes, err := readBlkioServiceBytes()
	    slicedBytes := totalBytes[6:]
	    newbytesInt, _ := strconv.Atoi(slicedBytes)
            if err != nil {
                fmt.Printf("[%s] Error reading blkio service bytes: %v\n", timestamp, err)
		    
            } else {
                fmt.Printf("[%s] %s bytes\n", timestamp, strings.TrimSpace(totalBytes))
		
		if newbytesInt > currentbytesInt {
		fmt.Println("Difference: +", newbytesInt - currentbytesInt,"\n")
		currentbytesInt = newbytesInt
	    	} else if newbytesInt < currentbytesInt {
		fmt.Println("Difference: -", newbytesInt - currentbytesInt,"\n")
		currentbytesInt = newbytesInt
	   	 } else if newbytesInt == currentbytesInt {
		fmt.Println("No difference.\n")
	    	}
		    
		fmt.Printf("[%s] %s read/writes\n", timestamp, strings.TrimSpace(totalReadWrite))
		    
		if newrwInt > currentrwInt {
		fmt.Println("Difference: +", newrwInt - currentrwInt,"\n")
		currentrwInt = newrwInt
	    	} else if newrwInt < currentrwInt {
		fmt.Println("Difference: -", newrwInt - currentrwInt,"\n")
		currentrwInt = newrwInt
	  	} else if newrwInt == currentrwInt {
		fmt.Println("No difference.\n")
	    }
           }

        case <-sigChan:
            return
    }
   }
  }



func readBlkioServiceBytes() (string, error) {
    filePath := fmt.Sprintf("%s/%s", blkioPath, serviceName)
    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        return "", err
    }

    lines := strings.Split(strings.TrimSpace(string(data)), "\n")
    if len(lines) > 0 {
        return lines[len(lines)-1], nil
    }

    return "", nil
}

func readBlkioReadWrite() (string, error) {
    filePath := fmt.Sprintf("%s/%s", blkioPath, serviceNewName)
    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        return "", err
    }

    lines := strings.Split(strings.TrimSpace(string(data)), "\n")
    if len(lines) > 0 {
        return lines[len(lines)-1], nil
    }

    return "", nil
}
