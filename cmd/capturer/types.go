package main

import (
	"github.com/cilium/ebpf/ringbuf"
)

const (
	DataBlockSize uint = 4096
)

type indexedRecord struct {
	index  int
	record ringbuf.Record
}

type dataBlock struct {
	dataBuf [DataBlockSize]uint8
	dataLen uint32
}

type indexedDataBlock struct {
	index     int
	offset    int64
	dataBlock dataBlock
}

func makeIndexedDataBlock(index int, offset int64, data [DataBlockSize]uint8, dataLen uint32) indexedDataBlock {
	return indexedDataBlock{
		index:     index,
		offset:    offset,
		dataBlock: dataBlock{dataBuf: data, dataLen: dataLen},
	}
}
