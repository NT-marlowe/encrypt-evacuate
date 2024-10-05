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

func processRingBufRecord(irdCh <-chan indexedRecord, idbCh chan indexedDataBlock, file *os.File) {
	go writeFileData(idbCh, file)

	for i := 0; i < Parallelism; i++ {
		go decodeIndexedRecord(irdCh, idbCh)
	}

}

func decodeIndexedRecord(irdCh <-chan indexedRecord, idbCh chan<- indexedDataBlock) {
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

		idbCh <- makeIndexedDataBlock(ird.index, event.Data, uint32(event.DataLen))
	}
}

func writeFileData(idbCh <-chan indexedDataBlock, file *os.File) {
	var idb indexedDataBlock
	var ok bool
	for {
		idb, ok = <-idbCh
		if !ok {
			log.Println("Data channel closed, exiting..")
			return
		}

		// log.Printf("idx: %d\n", idb.index)
		file.Write(idb.dataBlock[:idb.dataLen])
	}
}
