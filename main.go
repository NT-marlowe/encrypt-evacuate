package main

import (
	"bytes"
	"encoding/binary"
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
	// sharedLibraryPath = "/usr/lib/x86_64-linux-gnu/libssl.so"
	// sharedLibraryPath = "/usr/lib/python3.10/lib-dynload/_ssl.cpython-310-x86_64-linux-gnu.so"
	// symbol = "SSL_write"
	symbol = "EVP_EncryptUpdate"
	// symbol = "EVP_EncryptInit_ex"
)

func main() {
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

	var event capture_sslEncDataEventT
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

		if err := binary.Read(bytes.NewBuffer(record.RawSample), binary.LittleEndian, &event); err != nil {
			log.Printf("parsing ringbuf event: %s", err)
			continue
		}

		log.Printf("data: %s\n", string(event.Data[:event.DataLen]))

	}

	// log.Print("Received signal, exiting..")
}
