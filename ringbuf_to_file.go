package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"

	"github.com/cilium/ebpf/ringbuf"
)

func processRingBufRecord(recordCh <-chan ringbuf.Record, file *os.File) {
	var event capture_sslEncDataEventT
	dataCh := make(chan []byte)
	defer close(dataCh)

	go writeFileData(dataCh, file)

	for {
		record, ok := <-recordCh
		if !ok {
			log.Println("Record channel closed, exiting..")
			return
		}

		if err := binary.Read(bytes.NewBuffer(record.RawSample), binary.LittleEndian, &event); err != nil {
			log.Printf("parsing ringbuf event: %s", err)
			continue
		}

		dataCh <- event.Data[:event.DataLen]
		// file.Write(event.Data[:event.DataLen])

	}
}

func writeFileData(dataCh <-chan []byte, file *os.File) {
	var data []byte
	var ok bool
	for {
		data, ok = <-dataCh
		if !ok {
			log.Println("Data channel closed, exiting..")
			return
		}

		file.Write(data)
	}
}
