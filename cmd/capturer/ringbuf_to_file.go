package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"

	"github.com/cilium/ebpf/ringbuf"
)

// main goroutine
//
//	--> decodeIndexedRecord (multi goroutines)
//	--> writeFileData (single goroutine)
func startProcessingStages(irdCh <-chan ringbuf.Record, eventCh chan capture_plainEncDataEventT, parallelism int) {
	// go writeFileDataSequntial(eventCh, file)
	go writeFileDataOffset(eventCh)

	for i := 0; i < parallelism; i++ {
		go decodeIndexedRecord(irdCh, eventCh)
	}

}

func decodeIndexedRecord(irdCh <-chan ringbuf.Record, eventCh chan<- capture_plainEncDataEventT) {
	var event capture_plainEncDataEventT

	for {
		rd, ok := <-irdCh
		if !ok {
			log.Println("Record channel closed, exiting..")
			return
		}

		if err := binary.Read(bytes.NewBuffer(rd.RawSample), binary.LittleEndian, &event); err != nil {
			log.Printf("parsing ringbuf event: %s", err)
			continue
		}

		eventCh <- event
	}
}

func writeFileDataOffset(eventCh <-chan capture_plainEncDataEventT) {
	fileHandlerMap := make(map[string]*os.File, 0)
	for event := range eventCh {
		filename := bytesToString(event.Filename[:])
		// log.Printf("pwd = %s\n", bytesToString(event.Pwd[:]))
		file, ok := fileHandlerMap[filename]
		if ok {
			file.Seek(event.Offset, 0)
			file.Write(event.Data[:event.DataLen])
			continue
		}

		// log.Printf("Creating file: %s\n", filename)

		file, err := setupDataShelter(dataShelterPath, filename)
		if err != nil {
			log.Fatal("Creating file: ", err)
		}
		fileHandlerMap[filename] = file
		defer file.Close()

		file.Seek(event.Offset, 0)
		file.Write(event.Data[:event.DataLen])
	}
}

func bytesToString(data []int8) string {
	// int8 -> byte type cast
	byteData := make([]byte, len(data))
	for i, b := range data {
		byteData[i] = byte(b)
	}

	// trim strings after null character
	n := 0
	for n < len(byteData) && byteData[n] != 0 {
		n++
	}

	return string(byteData[:n])
}

func writeFileDataSequntial(eventCh <-chan indexedDataBlock, file *os.File) {
	m := make(map[int]dataBlock)
	currentIndex := 0
	var event indexedDataBlock
	var db dataBlock
	var ok bool

	for event = range eventCh {
		if event.index == currentIndex {
			db = event.dataBlock
			file.Write(db.dataBuf[:db.dataLen])
			currentIndex++
		} else {
			m[event.index] = event.dataBlock
			// better without using continue
			// continue
			// enqueueTime[event.index] = time.Now()
		}

		for {
			db, ok = m[currentIndex]
			if !ok {
				break
			}
			file.Write(db.dataBuf[:db.dataLen])
			delete(m, currentIndex)

			currentIndex++
		}
	}
}
