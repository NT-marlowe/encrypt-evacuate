package main

import (
	"github.com/cilium/ebpf/ringbuf"
)

type indexedRecord struct {
	index  int
	record ringbuf.Record
}
