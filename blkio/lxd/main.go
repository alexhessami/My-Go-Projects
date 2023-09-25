package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type Response struct {
	Type       string   `json:"type"`
	Status     string   `json:"status"`
	StatusCode int      `json:"status_code"`
	Operation  string   `json:"operation"`
	ErrorCode  int      `json:"error_code"`
	Error      string   `json:"error"`
	Metadata   []string `json:"metadata"`
}

type IOStats struct {
	Read    int `json:"read"`
	Write   int `json:"write"`
	Sync    int `json:"sync"`
	Async   int `json:"async"`
	Discard int `json:"discard"`
	Total   int `json:"total"`
}

func readIOStatistics(majMin string) (IOStats, error) {
	filePath := "/sys/fs/cgroup/blkio/blkio.throttle.io_service_bytes"
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return IOStats{}, err
	}

	stats := IOStats{}
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, majMin) {
			parts := strings.Fields(line)
			if len(parts) == 3 {
				value, err := strconv.Atoi(parts[2])
				if err != nil {
					return IOStats{}, err
				}

				switch parts[1] {
				case "Read":
					stats.Read = value
				case "Write":
					stats.Write = value
				case "Sync":
					stats.Sync = value
				case "Async":
					stats.Async = value
				case "Discard":
					stats.Discard = value
				case "Total":
					stats.Total = value
				}
			}
		}
	}

	return stats, nil
}

func normalizeDeviceName(deviceName string) string {
	// Replace extra dashes and underscores with a single dash
	deviceName = regexp.MustCompile(`[-_]+`).ReplaceAllString(deviceName, "-")

	return deviceName
}

func normalizeAPIName(apiName string) string {
	// Replace extra dashes and underscores with a single dash
	apiName = regexp.MustCompile(`[-_]+`).ReplaceAllString(apiName, "-")

	// Add four dashes between words
	apiName = strings.Replace(apiName, "-", "----", -1)

	return apiName
}

func main() {
	socketPath := "/var/snap/lxd/common/lxd/unix.socket"
	request := "GET /1.0/storage-pools/default/volumes/container HTTP/1.0\r\n\r\n"

	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		log.Fatalf("Error connecting to UNIX socket: %v\n", err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte(request))
	if err != nil {
		log.Fatalf("Error sending request: %v\n", err)
	}

	resp, err := http.ReadResponse(bufio.NewReader(conn), nil)
	if err != nil {
		log.Fatalf("Error reading HTTP response: %v\n", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v\n", err)
	}

	var jsonResponse Response
	decoder := json.NewDecoder(strings.NewReader(string(bodyBytes)))
	err = decoder.Decode(&jsonResponse)
	if err != nil {
		log.Fatalf("Error decoding JSON: %v\n", err)
	}

	//Check if the response contains metadata
	if len(jsonResponse.Metadata) == 0 {
		log.Fatalf("No storage pool metadata found in the response.\n")
	}

	//Process each storage pool path from the metadata
	for _, storagePoolPath := range jsonResponse.Metadata {
		//Removes "/1.0/storage-pools/default/volumes/container/"
		containerName := strings.TrimPrefix(storagePoolPath, "/1.0/storage-pools/default/volumes/container/")

		//Generate a standardized API name to match the lsblk format
		standardizedAPIName := fmt.Sprintf("container-containers_%s", normalizeAPIName(containerName))

		//Capture the MAJ:MIN for the device from lsblk
		output, err := exec.Command("lsblk", "--pairs", "--noheadings", "--output", "NAME,MAJ:MIN").Output()
		if err != nil {
			log.Fatalf("Error running lsblk: %v\n", err)
		}

		lsblkData := string(output)

		//Pattern to search for in lsblk output
		lsblkPattern := fmt.Sprintf("NAME=\"%s\"", standardizedAPIName)

		//Find the corresponding MAJ:MIN  
		re := regexp.MustCompile(fmt.Sprintf(`%s\s+MAJ:MIN="(\d+):(\d+)"`, lsblkPattern))
		matches := re.FindStringSubmatch(lsblkData)

		if len(matches) == 3 {
			maj, min := matches[1], matches[2]
			fmt.Printf("Storage Pool Name: %s\n", containerName)
			fmt.Printf("MAJ:MIN: %s:%s\n", maj, min)

			ioStats, err := readIOStatistics(fmt.Sprintf("%s:%s", maj, min))
			if err != nil {
				log.Fatalf("Error reading I/O statistics: %v\n", err)
			}

			fmt.Printf("I/O Statistics:\n")
			fmt.Printf("Read: %d\n", ioStats.Read)
			fmt.Printf("Write: %d\n", ioStats.Write)
			fmt.Printf("Sync: %d\n", ioStats.Sync)
			fmt.Printf("Async: %d\n", ioStats.Async)
			fmt.Printf("Discard: %d\n", ioStats.Discard)
			fmt.Printf("Total: %d\n\n", ioStats.Total)
		} else {
			fmt.Printf("No matching device found for %s\n\n", containerName)
		}
	}
}
