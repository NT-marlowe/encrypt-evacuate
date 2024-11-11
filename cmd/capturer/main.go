package main

import (
	"errors"
	"fmt"
	"log"

	// "time"

	"github.com/cilium/ebpf"
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
	parallelism, err := parseArgs()
	if err != nil {
		log.Fatal("Parsing args:", err)
	}

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

	upobe, err := attachUprobeProgram(&objs)
	if err != nil {
		log.Fatal("Attaching uprobe program:", err)
	}
	defer upobe.Close()

	links := attachAllTracingPrograms(&objs)
	for _, l := range links {
		defer l.Close()
	}

	// kp, err := link.Kretprobe("do_sys_openat2", objs.KretprobeOpenat2, nil)
	// if err != nil {
	// 	log.Fatal("Attaching kretprobe:", err)
	// }
	// defer kp.Close()

	rd, err := ringbuf.NewReader(objs.EventsRingbuf)
	if err != nil {
		log.Fatal("Creating ringbuf reader:", err)
	}
	defer rd.Close()

	startStopper(rd)

	recordCh := make(chan ringbuf.Record)
	defer close(recordCh)

	encDataEventCh := make(chan capture_plainEncDataEventT)
	defer close(encDataEventCh)

	// Starts decoding goroutines and a writing goroutine.
	startProcessingStages(recordCh, encDataEventCh, parallelism)

	// var start time.Time
	// var elapsed time.Duration
	// for index := 0; ; index++ {
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

		recordCh <- record
	}
}

func attachAllTracingPrograms(objs *capture_plainObjects) []link.Link {
	programs := []*ebpf.Program{
		objs.FexitDoSysOpen,
		objs.FentryKsysRead,
		objs.FexitKsysRead,
		objs.FexitKsysLseek,
	}

	var links []link.Link

	for _, prog := range programs {
		l, err := link.AttachTracing(link.TracingOptions{
			Program: prog,
		})
		if err != nil {
			log.Fatal("Attaching tracing:", err)
		}

		links = append(links, l)
	}
	return links
}

func attachUprobeProgram(objs *capture_plainObjects) (link.Link, error) {
	ex, err := link.OpenExecutable(sharedLibraryPath)
	if err != nil {
		return nil, fmt.Errorf("opening %s: %s", sharedLibraryPath, err)
	}

	uprobe, err := ex.Uprobe(symbol, objs.ProbeEntryEVP_EncryptUpdate, nil)
	if err != nil {
		return nil, fmt.Errorf("attach %s to uprobe", symbol)
	}

	return uprobe, nil
}
