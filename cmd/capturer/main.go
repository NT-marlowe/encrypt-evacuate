package main

import (
	"errors"
	// "fmt"
	"log"

	// "time"

	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/ringbuf"
	"github.com/cilium/ebpf/rlimit"
)

const (
	sharedLibraryPath = "/lib/x86_64-linux-gnu/libcrypto.so.3"
	symbol            = "EVP_EncryptUpdate"
	dataShelterPath   = "/data_shelter"
)

func main() {
	filename, parallelism := parseArgs()

	// Remove resource limits for kernels <5.11.
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatal("Removing memlock:", err)
	}

	// Load the compiled eBPF ELF and load it into the kernel.
	var objs capture_plainObjects
	if err := loadCapture_plainObjects(&objs, nil); err != nil {
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

	link_read_fexit, err := link.AttachTracing(link.TracingOptions{
		Program: objs.FentryKsysRead,
	})
	if err != nil {
		log.Fatal("Attaching tracing:", err)
	}
	defer link_read_fexit.Close()

	link_read_fentry, err := link.AttachTracing(link.TracingOptions{
		Program: objs.FexitKsysRead,
	})
	if err != nil {
		log.Fatal("Attaching tracing:", err)
	}
	defer link_read_fentry.Close()

	link_openat, err := link.AttachTracing(
		link.TracingOptions{Program: objs.FexitDoSysOpen},
	)
	if err != nil {
		log.Fatal("Attaching tracing:", err)
	}
	defer link_openat.Close()

	rd, err := ringbuf.NewReader(objs.EventsRingbuf)
	if err != nil {
		log.Fatal("Creating ringbuf reader:", err)
	}
	defer rd.Close()

	startStopper(rd)

	// create a file in dataShelterPath
	file, err := setupDataShelter(dataShelterPath, filename)
	if err != nil {
		log.Fatal("Set up data shelter: ", err)
	}
	defer file.Close()

	indexedRecordCh := make(chan indexedRecord)
	defer close(indexedRecordCh)

	indexedDataBlockCh := make(chan indexedDataBlock)
	defer close(indexedDataBlockCh)

	// main goroutine: processRingBufRecord
	//		--> decodeIndexedRecord (multi goroutines)
	//		--> writeFileData (single goroutine)
	processRingBufRecord(indexedRecordCh, indexedDataBlockCh, file, parallelism)

	var index int
	// var start time.Time
	// var elapsed time.Duration
	for {
		// start = time.Now()

		record, err := rd.Read()

		// elapsed = time.Since(start)
		// fmt.Printf("rd.Read: %v\n", elapsed)

		if err != nil {
			if errors.Is(err, ringbuf.ErrClosed) {
				log.Println("Ringbuf closed, exiting..")
				return
			}
			log.Printf("Reading record: %s", err)
			continue
		}

		indexedRecordCh <- indexedRecord{index: index, record: record}
		index++

	}
}
