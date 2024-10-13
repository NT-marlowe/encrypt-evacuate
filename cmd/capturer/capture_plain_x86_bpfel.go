// Code generated by bpf2go; DO NOT EDIT.
//go:build 386 || amd64

package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"

	"github.com/cilium/ebpf"
)

type capture_plainEncDataEventT struct {
	TimestampNs uint64
	Pid         uint32
	Tid         uint32
	Data        [4096]uint8
	DataLen     int32
	_           [4]byte
}

// loadCapture_plain returns the embedded CollectionSpec for capture_plain.
func loadCapture_plain() (*ebpf.CollectionSpec, error) {
	reader := bytes.NewReader(_Capture_plainBytes)
	spec, err := ebpf.LoadCollectionSpecFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("can't load capture_plain: %w", err)
	}

	return spec, err
}

// loadCapture_plainObjects loads capture_plain and converts it into a struct.
//
// The following types are suitable as obj argument:
//
//	*capture_plainObjects
//	*capture_plainPrograms
//	*capture_plainMaps
//
// See ebpf.CollectionSpec.LoadAndAssign documentation for details.
func loadCapture_plainObjects(obj interface{}, opts *ebpf.CollectionOptions) error {
	spec, err := loadCapture_plain()
	if err != nil {
		return err
	}

	return spec.LoadAndAssign(obj, opts)
}

// capture_plainSpecs contains maps and programs before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type capture_plainSpecs struct {
	capture_plainProgramSpecs
	capture_plainMapSpecs
}

// capture_plainSpecs contains programs before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type capture_plainProgramSpecs struct {
	ProbeEntryEVP_EncryptUpdate *ebpf.ProgramSpec `ebpf:"probe_entry_EVP_EncryptUpdate"`
}

// capture_plainMapSpecs contains maps before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type capture_plainMapSpecs struct {
	EventsRingbuf *ebpf.MapSpec `ebpf:"events_ringbuf"`
	PtrToFd       *ebpf.MapSpec `ebpf:"ptr_to_fd"`
}

// capture_plainObjects contains all objects after they have been loaded into the kernel.
//
// It can be passed to loadCapture_plainObjects or ebpf.CollectionSpec.LoadAndAssign.
type capture_plainObjects struct {
	capture_plainPrograms
	capture_plainMaps
}

func (o *capture_plainObjects) Close() error {
	return _Capture_plainClose(
		&o.capture_plainPrograms,
		&o.capture_plainMaps,
	)
}

// capture_plainMaps contains all maps after they have been loaded into the kernel.
//
// It can be passed to loadCapture_plainObjects or ebpf.CollectionSpec.LoadAndAssign.
type capture_plainMaps struct {
	EventsRingbuf *ebpf.Map `ebpf:"events_ringbuf"`
	PtrToFd       *ebpf.Map `ebpf:"ptr_to_fd"`
}

func (m *capture_plainMaps) Close() error {
	return _Capture_plainClose(
		m.EventsRingbuf,
		m.PtrToFd,
	)
}

// capture_plainPrograms contains all programs after they have been loaded into the kernel.
//
// It can be passed to loadCapture_plainObjects or ebpf.CollectionSpec.LoadAndAssign.
type capture_plainPrograms struct {
	ProbeEntryEVP_EncryptUpdate *ebpf.Program `ebpf:"probe_entry_EVP_EncryptUpdate"`
}

func (p *capture_plainPrograms) Close() error {
	return _Capture_plainClose(
		p.ProbeEntryEVP_EncryptUpdate,
	)
}

func _Capture_plainClose(closers ...io.Closer) error {
	for _, closer := range closers {
		if err := closer.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Do not access this directly.
//
//go:embed capture_plain_x86_bpfel.o
var _Capture_plainBytes []byte
