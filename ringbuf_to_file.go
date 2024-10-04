package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"

	"github.com/cilium/ebpf/ringbuf"
)

type tmp_data struct {
	data [4096]uint8
	len  uint32
}

func processRingBufRecord(recordCh <-chan ringbuf.Record, file *os.File) {
	var event capture_sslEncDataEventT
	dataCh := make(chan tmp_data)
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

		dataCh <- tmp_data{data: event.Data, len: uint32(event.DataLen)}
		// file.Write(event.Data[:event.DataLen])

	}
}

// func writeFileData(dataCh <-chan []byte, file *os.File) {
func writeFileData(dataCh <-chan tmp_data, file *os.File) {
	// var data []byte
	var data tmp_data
	var ok bool
	for {
		data, ok = <-dataCh
		if !ok {
			log.Println("Data channel closed, exiting..")
			return
		}

		file.Write(data.data[:data.len])
	}
}
