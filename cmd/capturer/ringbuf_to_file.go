package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"time"
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

	var start time.Time
	var elapsed time.Duration
	for {
		ird, ok := <-irdCh
		start = time.Now()
		if !ok {
			log.Println("Record channel closed, exiting..")
			return
		}

		if err := binary.Read(bytes.NewBuffer(ird.record.RawSample), binary.LittleEndian, &event); err != nil {
			log.Printf("parsing ringbuf event: %s", err)
			continue
		}
		elapsed = time.Since(start)
		fmt.Printf("binary.Read: %v\n", elapsed)

		idbCh <- makeIndexedDataBlock(ird.index, event.Data, uint32(event.DataLen))
	}
}

// slice, key: Item.index, value: time.TIme
var enqueueTime = make(map[int]time.Time)

func measureTime(index int, op string) {
	t, ok := enqueueTime[index]
	if !ok {
		fmt.Printf("No enqueue time found for index %d\n", index)
		return
	}
	elapsed := time.Since(t)

	fmt.Printf("%s: %v\n", op, elapsed)
	delete(enqueueTime, index)
}

func writeFileData(idbCh <-chan indexedDataBlock, file *os.File) {
	m := make(map[int]dataBlock)
	currentIndex := 0
	var idb indexedDataBlock
	var db dataBlock
	var ok bool

	for idb = range idbCh {
		// fmt.Printf("idb.index: %d, currentIndex: %d\n", idb.index, currentIndex)
		if idb.index == currentIndex {
			db = idb.dataBlock
			file.Write(db.dataBuf[:db.dataLen])
			fmt.Printf("from chan, idx = %d\n", idb.index)
			currentIndex++
		} else {
			m[idb.index] = idb.dataBlock
			enqueueTime[idb.index] = time.Now()
			fmt.Printf("idx = %d was enqueued\n", idb.index)
		}

		for {
			db, ok = m[currentIndex]
			if !ok {
				break
			}
			file.Write(db.dataBuf[:db.dataLen])
			delete(m, currentIndex)
			fmt.Printf("from map, idx = %d\n", currentIndex)
			measureTime(currentIndex, "writeFileData")

			currentIndex++
		}
	}
}
