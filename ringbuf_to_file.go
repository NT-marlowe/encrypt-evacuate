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

func processRingBufRecord(indexedRecordCh <-chan indexedRecord, indexedDataBlockCh chan indexedDataBlock, file *os.File) {
	go writeFileData(indexedDataBlockCh, file)

	for i := 0; i < Parallelism; i++ {
		go decodeIndexedRecord(indexedRecordCh, indexedDataBlockCh)
	}

}

func decodeIndexedRecord(irdCh <-chan indexedRecord, indexedDataBlockCh chan<- indexedDataBlock) {
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

		indexedDataBlockCh <- makeIndexedDataBlock(ird.index, event.Data, uint32(event.DataLen))
	}
}

func writeFileData(indexedDataBlockCh <-chan indexedDataBlock, file *os.File) {
	var idb indexedDataBlock
	var ok bool
	for {
		idb, ok = <-indexedDataBlockCh
		if !ok {
			log.Println("Data channel closed, exiting..")
			return
		}

		log.Printf("idx: %d\n", idb.index)
		file.Write(idb.data.data[:idb.data.len])
	}
}
