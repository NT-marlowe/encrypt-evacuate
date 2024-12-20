package main

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"

	"github.com/cilium/ebpf/ringbuf"
)

func parseArgs() (int, error) {
	if len(os.Args) != 1 && len(os.Args) != 2 {
		log.Fatalf("Usage: %s [parallelism]", os.Args[0])
	}

	var parallelism int
	var err error
	if len(os.Args) == 1 {
		// p = 15 is the tenttavely the best value for parallelism.
		parallelism = 15
	} else {
		parallelism, err = strconv.Atoi(os.Args[1])
		if err != nil {
			return 0, err
		}
	}
	return parallelism, err
}

func makeDataShleter(dataShelterPath string) error {
	err := os.MkdirAll(dataShelterPath, 0755)
	return err
}

func createShelteredFile(dataShelterPath string, filename string) (*os.File, error) {
	fullPath := filepath.Join(dataShelterPath, filename)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		log.Fatal("Creating dir: ", err)
	}

	file, err := os.Create(fullPath)
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
