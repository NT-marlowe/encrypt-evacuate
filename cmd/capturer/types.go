package main

import (
	"github.com/cilium/ebpf/ringbuf"
)

const (
	dataBlockSize uint = 4096
	fileNameLen   uint = 256
)

type indexedRecord struct {
	index  int
	record ringbuf.Record
}

type dataBlock struct {
	dataBuf [dataBlockSize]uint8
	dataLen uint32
}

type indexedDataBlock struct {
	index     int
	offset    int64
	filename  [fileNameLen]int8
	dataBlock dataBlock
}

func makeIndexedDataBlock(index int, offset int64, filename [fileNameLen]int8,
	data [dataBlockSize]uint8, dataLen uint32) indexedDataBlock {

	return indexedDataBlock{
		index:     index,
		offset:    offset,
		filename:  filename,
		dataBlock: dataBlock{dataBuf: data, dataLen: dataLen},
	}
}
