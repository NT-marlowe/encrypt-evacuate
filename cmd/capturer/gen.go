package main

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -type enc_data_event_t -target amd64 capture_plain ../../ebpf_src/capture_plain.c
