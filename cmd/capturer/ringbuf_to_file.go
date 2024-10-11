package main

import (
	"bytes"
	"encoding/binary"
	// "fmt"
	"log"
	"os"
	// "time"
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

	// var start time.Time
	// var elapsed time.Duration
	for {
		ird, ok := <-irdCh
		// start = time.Now()
		if !ok {
			log.Println("Record channel closed, exiting..")
			return
		}

		if err := binary.Read(bytes.NewBuffer(ird.record.RawSample), binary.LittleEndian, &event); err != nil {
			log.Printf("parsing ringbuf event: %s", err)
			continue
		}
		// elapsed = time.Since(start)
		// fmt.Printf("binary.Read: %v\n", elapsed)

		idbCh <- makeIndexedDataBlock(ird.index, event.Data, uint32(event.DataLen))
	}
}

func writeFileData(idbCh <-chan indexedDataBlock, file *os.File) {
	restoredCh := restoreOrder(idbCh)

	// var start time.Time
	// var elapsed time.Duration
	for item := range restoredCh {
		// start = time.Now()
		idb := item.dataBlock
		// log.Printf("idx: %d\n", item.GetIndex())

		file.Write(idb.dataBuf[:idb.dataLen])
		// elapsed = time.Since(start)
		// fmt.Printf("file.Write: %v\n", elapsed)
	}

}

func map_write_tmp(idbCh <-chan indexedDataBlock, file *os.File) {
	m := make(map[int]dataBlock)
	currentIndex := 0
	var db dataBlock
	var ok bool

	for {
		select {
		case idb := <-idbCh:
			if idb.index == currentIndex {
				dataBlock := idb.dataBlock
				file.Write(dataBlock.dataBuf[:dataBlock.dataLen])
				currentIndex++
			} else {
				m[idb.index] = idb.dataBlock
			}
		default:
			if db, ok = m[currentIndex]; ok {
				file.Write(db.dataBuf[:db.dataLen])
				delete(m, currentIndex)
				currentIndex++
			}

		}
	}
}
