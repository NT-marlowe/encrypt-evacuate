package main

import (
	"errors"
	"log"
	"os"
	"os/signal"

	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/ringbuf"
	"github.com/cilium/ebpf/rlimit"
)

const (
	sharedLibraryPath = "/lib/x86_64-linux-gnu/libcrypto.so.3"
	symbol            = "EVP_EncryptUpdate"
	dataShelterPath   = "/data_shelter"

	ChannelBufferSize = 1000
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s filename", os.Args[0])
	}
	filename := os.Args[1]

	// Remove resource limits for kernels <5.11.
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatal("Removing memlock:", err)
	}

	// Load the compiled eBPF ELF and load it into the kernel.
	var objs capture_sslObjects
	if err := loadCapture_sslObjects(&objs, nil); err != nil {
		log.Fatal("Loading eBPF objects:", err)
	}
	defer objs.Close()

	ex, err := link.OpenExecutable(sharedLibraryPath)
	if err != nil {
		log.Fatalf("Opening %s: %s", sharedLibraryPath, err)
	}

	uprobe, err := ex.Uprobe(symbol, objs.ProbeEntryEVP_EncryptUpdate, nil)
	if err != nil {
		log.Fatalf("Uprobe %s: %s", symbol, err)
	}
	defer uprobe.Close()

	rd, err := ringbuf.NewReader(objs.EventsRingbuf)
	if err != nil {
		log.Fatal("Creating ringbuf reader:", err)
	}
	defer rd.Close()

	stopper := make(chan os.Signal, 5)
	signal.Notify(stopper, os.Interrupt)

	go func() {
		<-stopper

		if err := rd.Close(); err != nil {
			log.Fatalf("Closing ringbuf reader: %v", err)
		}

	}()

	// create a file in dataShelterPath
	err = os.MkdirAll(dataShelterPath, 0766)
	if err != nil {
		log.Fatal("Creating data shelter path:", err)
	}
	file, err := os.Create(dataShelterPath + "/" + filename)
	if err != nil {
		log.Fatal("Creating file in data shelter path:", err)
	}
	defer file.Close()

	recordCh := make(chan ringbuf.Record, ChannelBufferSize)
	defer close(recordCh)
	go processRingBufRecord(recordCh, file)

	for {
		record, err := rd.Read()

		if err != nil {
			if errors.Is(err, ringbuf.ErrClosed) {
				log.Println("Ringbuf closed, exiting..")
				return
			}
			log.Printf("Reading record: %s", err)
			continue
		}

		recordCh <- record
	}
}
