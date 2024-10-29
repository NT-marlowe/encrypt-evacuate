package main

import (
	"bytes"
	"encoding/binary"

	// "fmt"
	"log"
	"os"
	// "time"
)

// main goroutine
//
//	--> decodeIndexedRecord (multi goroutines)
//	--> writeFileData (single goroutine)
func startProcessingStages(irdCh <-chan indexedRecord, idbCh chan indexedDataBlock, file *os.File, parallelism int) {
	go writeFileData(idbCh, file)

	for i := 0; i < parallelism; i++ {
		go decodeIndexedRecord(irdCh, idbCh)
	}

}

func decodeIndexedRecord(irdCh <-chan indexedRecord, idbCh chan<- indexedDataBlock) {
	var event capture_plainEncDataEventT

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

		// log.Printf("offset = %d, filename = %s\n", event.Offset, bytesToString(event.Filename[:]))
		// elapsed = time.Since(start)
		// fmt.Printf("binary.Read: %v\n", elapsed)

		idbCh <- makeIndexedDataBlock(ird.index, event.Offset, event.Data, uint32(event.DataLen))
	}
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
			currentIndex++
		} else {
			m[idb.index] = idb.dataBlock
			// better without using continue
			// continue
			// enqueueTime[idb.index] = time.Now()
		}

		for {
			db, ok = m[currentIndex]
			if !ok {
				break
			}
			file.Write(db.dataBuf[:db.dataLen])
			delete(m, currentIndex)
			// measureTime(currentIndex, "writeFileData")

			currentIndex++
		}
	}
}

func writeFileDataOffset(idbCh <-chan indexedDataBlock, file *os.File) {
	for idb := range idbCh {
		file.Seek(idb.offset, 0)
		file.Write(idb.dataBlock.dataBuf[:idb.dataBlock.dataLen])
	}
}

func bytesToString(data []int8) string {
	// int8 -> byteの型変換
	byteData := make([]byte, len(data))
	for i, b := range data {
		byteData[i] = byte(b)
	}

	// 0バイトで終端されているので、それ以降をトリム
	n := 0
	for n < len(byteData) && byteData[n] != 0 {
		n++
	}

	return string(byteData[:n])
}

// slice, key: Item.index, value: time.TIme
// var enqueueTime = make(map[int]time.Time)

// func measureTime(index int, op string) {
// 	t, ok := enqueueTime[index]
// 	if !ok {
// 		// fmt.Printf("No enqueue time found for index %d\n", index)
// 		return
// 	}
// 	elapsed := time.Since(t)

// 	fmt.Printf("%s: %v\n", op, elapsed)
// 	delete(enqueueTime, index)
// }
