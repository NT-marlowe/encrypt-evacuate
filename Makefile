BIN=ebpf-ssl

.PHONY: all
all: bpf_bpfel.go
	go build

bpf_bpfel.go: capture_ssl.c
	go generate

capture_ssl.o: capture_ssl.c
	clang -O2 -g -target bpf -c -o capture_ssl.o capture_ssl.c

.PHONY: run
run: all
	sudo ./${BIN} tmp

.PHONY: dump
dump: capture_ssl.o
	llvm-objdump -S capture_ssl.o

.PHONY: load
load: capture_ssl.o
	sudo bpftool prog load capture_ssl.o /sys/fs/bpf/my_program -d

.PHONY: cat
cat:
	sudo cat /sys/kernel/debug/tracing/trace_pipe

.PHONE: clear-shelter
clear-shelter:
	sudo rm -f /data_shelter/*
