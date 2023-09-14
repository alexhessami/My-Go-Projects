package main

import (
    "fmt"
    "io/ioutil"
    "os"
    "os/exec"
    "os/signal"

    "strings"
    "syscall"
    "time"
)

const blkioPath = "/sys/fs/cgroup/blkio/"

func main() {
    interrupt := make(chan os.Signal, 1)
    signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

    cleanup := func() {
        os.Exit(0)
    }

    go func() {
        <-interrupt
        cleanup()
    }()

    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            cmd := exec.Command("lsblk", "--output", "NAME,MAJ:MIN", "--noheadings")
            output, err := cmd.Output()
            if err != nil {
                fmt.Println("Error running lsblk:", err)
                cleanup()
            }

            poolInfo := parseLsblkOutput(string(output))

            for _, info := range poolInfo {
                stats, err := readBlkioStats(info.MajMin)
                if err != nil {
                    fmt.Printf("Error reading blkio stats for %s (%s): %v\n", info.Name, info.MajMin, err)
                    continue
                }
                fmt.Printf("Pool Name: %s\n", info.Name)
                fmt.Printf("MAJ:MIN Number: %s\n", info.MajMin)
                fmt.Printf("Stats:\n")
                for key, value := range stats {
                    fmt.Printf("  %s: %s\n", key, value)
                }
                fmt.Println()
            }
        }
    }
}

type PoolInfo struct {
    Name   string
    MajMin string
}

func parseLsblkOutput(output string) []PoolInfo {
    lines := strings.Split(output, "\n")
    var poolInfo []PoolInfo
    for _, line := range lines {
        fields := strings.Fields(line)
        if len(fields) >= 2 {
            name := fields[0]
            majMin := fields[1]
            poolInfo = append(poolInfo, PoolInfo{Name: name, MajMin: majMin})
        }
    }
    return poolInfo
}

func readBlkioStats(majMin string) (map[string]string, error) {
    stats := make(map[string]string)
    filePath := blkioPath + "blkio.throttle.io_service_bytes"
    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        return nil, err
    }

    lines := strings.Split(string(data), "\n")
    for _, line := range lines {
        fields := strings.Fields(line)
        if len(fields) < 3 {
            continue
        }
        currentMajMin := fields[0]
        if currentMajMin == majMin {
            stats[fields[1]] = fields[2]
        }
    }

    return stats, nil
}
