package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"
)

type dataBlock struct {
	data [4096]uint8
	len  uint32
}

func processRingBufRecord(indexedRecordCh <-chan indexedRecord, file *os.File) {
	var event capture_sslEncDataEventT
	dataBlockCh := make(chan dataBlock)
	defer close(dataBlockCh)

	go writeFileData(dataBlockCh, file)

	for {
		val, ok := <-indexedRecordCh
		if !ok {
			log.Println("Record channel closed, exiting..")
			return
		}

		if err := binary.Read(bytes.NewBuffer(val.record.RawSample), binary.LittleEndian, &event); err != nil {
			log.Printf("parsing ringbuf event: %s", err)
			continue
		}

		log.Printf("%d ", val.index)

		dataBlockCh <- dataBlock{data: event.Data, len: uint32(event.DataLen)}
		// file.Write(event.Data[:event.DataLen])

	}
}

func writeFileData(dataBlockCh <-chan dataBlock, file *os.File) {
	var data dataBlock
	var ok bool
	for {
		data, ok = <-dataBlockCh
		if !ok {
			log.Println("Data channel closed, exiting..")
			return
		}

		file.Write(data.data[:data.len])
	}
}
