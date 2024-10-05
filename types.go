package main

import (
	"github.com/cilium/ebpf/ringbuf"
)

type indexedRecord struct {
	index  int
	record ringbuf.Record
}

type dataBlock struct {
	data [4096]uint8
	len  uint32
}

// type indexedDataBlock struct {
// 	index int
// 	data  dataBlock
// }
