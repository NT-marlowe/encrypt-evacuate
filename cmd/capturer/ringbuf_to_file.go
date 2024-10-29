package main

import (
	"bytes"
	"encoding/binary"

	// "fmt"
	"log"
	"os"

	"github.com/cilium/ebpf/ringbuf"
	// "time"
)

// main goroutine
//
//	--> decodeIndexedRecord (multi goroutines)
//	--> writeFileData (single goroutine)
func startProcessingStages(irdCh <-chan ringbuf.Record, idbCh chan indexedDataBlock, parallelism int) {
	// go writeFileDataSequntial(idbCh, file)
	go writeFileDataOffset(idbCh)

	for i := 0; i < parallelism; i++ {
		go decodeIndexedRecord(irdCh, idbCh)
	}

}

func decodeIndexedRecord(irdCh <-chan ringbuf.Record, idbCh chan<- indexedDataBlock) {
	var event capture_plainEncDataEventT

	// var start time.Time
	// var elapsed time.Duration
	for {
		rd, ok := <-irdCh
		// start = time.Now()
		if !ok {
			log.Println("Record channel closed, exiting..")
			return
		}

		if err := binary.Read(bytes.NewBuffer(rd.RawSample), binary.LittleEndian, &event); err != nil {
			log.Printf("parsing ringbuf event: %s", err)
			continue
		}

		// log.Printf("offset = %d, filename = %s\n", event.Offset, bytesToString(event.Filename[:]))
		// elapsed = time.Since(start)
		// fmt.Printf("binary.Read: %v\n", elapsed)

		idbCh <- makeIndexedDataBlock(0, event.Offset, event.Filename, event.Data, uint32(event.DataLen))
	}
}

func writeFileDataOffset(idbCh <-chan indexedDataBlock) {
	fileHandlerMap := make(map[string]*os.File, 0)
	for idb := range idbCh {
		filename := bytesToString(idb.filename[:])
		file, ok := fileHandlerMap[filename]
		if ok {
			file.Seek(idb.offset, 0)
			file.Write(idb.dataBlock.dataBuf[:idb.dataBlock.dataLen])
			continue
		}

		file, err := setupDataShelter(dataShelterPath, filename)
		if err != nil {
			log.Fatal("Creating file: ", err)
		}
		fileHandlerMap[filename] = file
		defer file.Close()

		file.Seek(idb.offset, 0)
		file.Write(idb.dataBlock.dataBuf[:idb.dataBlock.dataLen])
	}
}

func writeFileDataSequntial(idbCh <-chan indexedDataBlock, file *os.File) {
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
