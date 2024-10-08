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
		start = time.Now()
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
		elapsed = time.Since(start)
		fmt.Printf("rd.Read: %v\n", elapsed)
	}
}

func writeFileData(idbCh <-chan indexedDataBlock, file *os.File) {
	// itemCh := make(chan priority_queue.Item)

	// go func() {
	// 	defer close(itemCh)
	// 	for idb := range idbCh {
	// 		itemCh <- priority_queue.MakeItem(idb.index, idb.dataBlock)
	// 	}
	// }()

	restoredCh := restoreOrder(idbCh)

	for item := range restoredCh {
		idb := item.dataBlock
		// log.Printf("idx: %d\n", item.GetIndex())

		file.Write(idb.dataBuf[:idb.dataLen])
	}

}
