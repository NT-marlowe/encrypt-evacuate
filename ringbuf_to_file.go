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

func decodeIndexedRecord(irdCh <-chan indexedRecord, dataBlockCh chan<- dataBlock) {
	var event capture_sslEncDataEventT

	for {
		ird, ok := <-irdCh
		if !ok {
			log.Println("Record channel closed, exiting..")
			return
		}

		if err := binary.Read(bytes.NewBuffer(ird.record.RawSample), binary.LittleEndian, &event); err != nil {
			log.Printf("parsing ringbuf event: %s", err)
			continue
		}

		dataBlockCh <- dataBlock{data: event.Data, len: uint32(event.DataLen)}
	}
}

func processRingBufRecord(indexedRecordCh <-chan indexedRecord, file *os.File) {
	if _, ok := <-indexedRecordCh; !ok {
		log.Println("Record channel closed, exiting..")
		return
	}

	dataBlockCh := make(chan dataBlock)
	defer close(dataBlockCh)

	go writeFileData(dataBlockCh, file)

	decodeIndexedRecord(indexedRecordCh, dataBlockCh)
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
