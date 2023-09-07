package main

import (
    "fmt"
    "io/ioutil"
    "os"
    "os/signal"
    "syscall"
    "time"
    "strings"
)

const (
    blkioPath = "/sys/fs/cgroup/blkio/"
    serviceName = "blkio.throttle.io_service_bytes"
    serviceNewName = "blkio.throttle.io_serviced"
    logInterval = 1 * time.Minute
)

func main() {
    // Create a channel to handle Ctrl+C signals.
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    // Create a ticker to collect data at regular intervals.
    ticker := time.NewTicker(logInterval)

    // Start the data collection loop.
    for {
        select {
        case <-ticker.C:
            timestamp := time.Now().Format(time.RFC3339)
	    totalReadWrite, err := readBlkioReadWrite()
            if err != nil {
                fmt.Printf("[%s] Error reading blkio read writes: %v\n", timestamp, err)
	}
            totalBytes, err := readBlkioServiceBytes()
            if err != nil {
                fmt.Printf("[%s] Error reading blkio service bytes: %v\n", timestamp, err)
            } else {
                fmt.Printf("[%s] %s bytes\n", timestamp, strings.TrimSpace(totalBytes))
		fmt.Printf("[%s] %s read/writes\n", timestamp, strings.TrimSpace(totalReadWrite))
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
