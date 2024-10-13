package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"

	"github.com/cilium/ebpf/ringbuf"
)

func parseArgs() (string, int) {
	if len(os.Args) != 2 && len(os.Args) != 3 {
		log.Fatalf("Usage: %s filename [parallelism]", os.Args[0])
	}

	filename := os.Args[1]
	var parallelism int
	var err error
	if len(os.Args) == 2 {
		// p = 15 is the tenttavely the best value for parallelism.
		parallelism = 15
	} else {
		parallelism, err = strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("Invalid parallelism: %s", err)
		}
	}
	return filename, parallelism
}

func makeDataShleter(dataShelterPath string) error {
	err := os.MkdirAll(dataShelterPath, 0755)
	return err
}

func createShelteredFile(dataShelterPath string, filename string) (*os.File, error) {
	file, err := os.Create(dataShelterPath + "/" + filename)
	return file, err
}

// Creates the data shelter dir and the file to store the captured data.
func setupDataShelter(dataShelterPath, filename string) (*os.File, error) {
	err := makeDataShleter(dataShelterPath)
	if err != nil {
		return nil, err
	}

	file, err := createShelteredFile(dataShelterPath, filename)
	return file, err
}

func startStopper(rd *ringbuf.Reader) {
	stopper := make(chan os.Signal, 5)
	signal.Notify(stopper, os.Interrupt)

	go func() {
		<-stopper

		if err := rd.Close(); err != nil {
			log.Fatalf("Closing ringbuf reader: %v", err)
		}

	}()
}
