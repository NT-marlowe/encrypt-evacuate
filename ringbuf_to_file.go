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

func processRingBufRecord(indexedRecordCh <-chan indexedRecord, indexedDataBlockCh chan dataBlock, file *os.File) {
	go writeFileData(indexedDataBlockCh, file)

	for i := 0; i < Parallelism; i++ {
		go decodeIndexedRecord(indexedRecordCh, indexedDataBlockCh)
	}

}

func decodeIndexedRecord(irdCh <-chan indexedRecord, indexedDataBlockCh chan<- dataBlock) {
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

		indexedDataBlockCh <- dataBlock{data: event.Data, len: uint32(event.DataLen)}
	}
}

func writeFileData(indexedDataBlockCh <-chan dataBlock, file *os.File) {
	var data dataBlock
	var ok bool
	for {
		data, ok = <-indexedDataBlockCh
		if !ok {
			log.Println("Data channel closed, exiting..")
			return
		}

		file.Write(data.data[:data.len])
	}
}
