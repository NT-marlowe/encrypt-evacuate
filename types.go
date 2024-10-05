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

type dataBlock [DataBlockSize]uint8

type indexedDataBlock struct {
	index     int
	dataBlock dataBlock
	dataLen   uint32
}

func makeIndexedDataBlock(index int, data [DataBlockSize]uint8, dataLen uint32) indexedDataBlock {
	return indexedDataBlock{
		index:     index,
		dataBlock: data,
		dataLen:   dataLen,
	}
}
