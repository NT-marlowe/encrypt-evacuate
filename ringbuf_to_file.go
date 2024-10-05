package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"
)

const (
	Parallelism = 4
)

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

func processRingBufRecord(indexedRecordCh <-chan indexedRecord, dataBlockCh chan dataBlock, file *os.File) {
	go writeFileData(dataBlockCh, file)

	for i := 0; i < Parallelism; i++ {
		go decodeIndexedRecord(indexedRecordCh, dataBlockCh)
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
