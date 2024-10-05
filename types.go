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
	data [DataBlockSize]uint8
	len  uint32
}

type indexedDataBlock struct {
	index int
	data  dataBlock
}

func makeIndexedDataBlock(index int, data [DataBlockSize]uint8, dataLen uint32) indexedDataBlock {
	return indexedDataBlock{
		index: index,
		data:  dataBlock{data: data, len: dataLen},
	}
}
