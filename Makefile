BIN := ebpf-ssl
NAME := capturer
TMP_FILE_NAME := tmp

GO_SRCS := $(shell find . -name '*.go')


.PHONY: all
all: ${GO_SRCS} gen
	go build -o ${BIN} ./cmd/${NAME}

# bpf_bpfel.go: ebpf_src/capture_ssl.c
# 	go generate

.PHONY: gen
gen: ebpf_src/capture_ssl.c
	go generate cmd/capturer/gen.go

# capture_ssl.o: capture_ssl.c
# 	clang -O2 -g -target bpf -c -o capture_ssl.o capture_ssl.c


.PHONY: run
run: all
	if [ -f /data_shelter/${TMP_FILE_NAME} ]; then \
		sudo rm /data_shelter/${TMP_FILE_NAME}; \
		echo "Remove old file"; \
	fi
	sudo ./${BIN} ${TMP_FILE_NAME} 4

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
