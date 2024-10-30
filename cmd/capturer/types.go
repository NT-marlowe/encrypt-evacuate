package main

const (
	dataBlockSize uint = 4096
	fileNameLen   uint = 256
)

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
