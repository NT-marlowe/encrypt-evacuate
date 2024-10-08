package main

import (
	"bytes"
	"ebpf-ssl/internal/priority_queue"
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

func writeFileData(idbCh <-chan indexedDataBlock, file *os.File) {
	itemCh := make(chan priority_queue.Item)

	go func() {
		defer close(itemCh)
		var start time.Time
		var elapsed time.Duration

		for idb := range idbCh {
			start = time.Now()
			itemCh <- priority_queue.MakeItem(idb.index, idb.dataBlock)
			elapsed = time.Since(start)
			fmt.Printf("MakeItem: %v\n", elapsed)
		}
	}()

	restoredCh := priority_queue.RestoreOrder(itemCh)

	var start time.Time
	var elapsed time.Duration
	for item := range restoredCh {
		start = time.Now()
		idb := item.GetValue().(dataBlock)
		// log.Printf("idx: %d\n", item.GetIndex())
		file.Write(idb.dataBuf[:idb.dataLen])
		elapsed = time.Since(start)
		fmt.Printf("file.Write: %v\n", elapsed)
	}

}
