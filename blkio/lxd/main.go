package main

import (
    "fmt"
    "github.com/lxc/incus/client"
    "os/signal"
    "os"
    "syscall"
)

func main() {

    interrupt := make(chan os.Signal, 1)
    signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

    cleanup := func() {
          os.Exit(0)
      }


    // Create a new LXD client
    client, err := incus.ConnectIncusUnix("", nil)
    if err != nil {
        fmt.Println("Error connecting to LXD:", err)
	      cleanup()
    }

    poolNames, err := client.GetStoragePoolNames()
    if err != nil {
    fmt.Println("Error getting storage pool names:", err)
	  cleanup()
    }

    for _, poolName := range poolNames {
            volumes, err := client.GetStoragePoolVolumes(poolName)
            if err != nil {
                fmt.Println("Error getting storage pool volumes:", err)
	              cleanup()
                    }


            for _, volume := range volumes {
                fmt.Printf("Volume Name: %s\n", volume.Name)
                fmt.Printf("Type: %s\n", volume.Type)
                fmt.Printf("Size: %s\n", volume.Config["size"])
                fmt.Println()
		return
            }
        }
    }
