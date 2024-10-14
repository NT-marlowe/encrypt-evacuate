BIN := ebpf-ssl
NAME := capturer
TMP_FILE_NAME := tmp
USER := $(shell whoami)

GO_SRCS := $(shell find . -name '*.go')


.PHONY: all
all: ${GO_SRCS} gen
	go build -o ${BIN} ./cmd/${NAME}

# bpf_bpfel.go: ebpf_src/capture_plain.c
# 	go generate

.PHONY: gen
gen: ebpf_src/capture_plain.c chown
	go generate cmd/capturer/gen.go

.PHONE: chown
chown:
	sudo chown -R ${USER}:${USER} ./cmd

capture_plain.o: ebpf_src/capture_plain.c
	cd ebpf_src && \
		clang -O2 -g -target bpf -D__TARGET_ARCH_x86 \
		-c capture_plain.c -o capture_plain.o


.PHONY: run
run: all
	if [ -f /data_shelter/${TMP_FILE_NAME} ]; then \
		sudo rm /data_shelter/${TMP_FILE_NAME}; \
		echo "Remove old file"; \
	fi
	sudo ./${BIN} ${TMP_FILE_NAME} 4

.PHONY: dump
dump: capture_plain.o
	cd ebpf_src && \
		llvm-objdump -S capture_plain.o

.PHONY: load
load: capture_plain.o
	cd ebpf_src && \
		sudo bpftool prog load capture_plain.o /sys/fs/bpf/my_program -d

.PHONY: cat
cat:
	sudo cat /sys/kernel/debug/tracing/trace_pipe

.PHONE: clear-shelter
clear-shelter:
	sudo rm -f /data_shelter/*
