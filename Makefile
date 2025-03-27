BIN := ebpf-ssl
NAME := capturer
TMP_FILE_NAME := tmp
USER := $(shell whoami)

GO_SRCS := $(shell find . -name '*.go')

# 実行可能なシングルバイナリを作成する．このバイナリにはeBPFプログラムがバイトコードとして埋め込まれている．
.PHONY: all
all: ${GO_SRCS} gen
	go build -o ${BIN} ./cmd/${NAME}


# eBPF Cソースファイル（capture_plain.c）をGoバインディングに変換する．
# cmd/capturer/gen.go に記述された `//go:generate` ディレクティブを実行することで，
# Cilium eBPFツールチェイン（bpf2go）を使って，eBPFプログラムとデータ構造（enc_data_event_t）に
# 対応するGoコード（capture_plain_bpfel.go）を自動生成する．
# - このGoコードは，ユーザー空間のGoアプリケーションからeBPFオブジェクトを操作するために必要．
# - ターゲットアーキテクチャはamd64に固定されている．
# - bpf2goの出力先はcmd/capturer/ディレクトリ内となる．
# また，事前に `./cmd` 以下の所有権をユーザに変更しておくことで，
# `go generate` 実行時に生成ファイルの書き込みエラーを防ぐ．
.PHONY: gen
gen: ebpf_src/capture_plain.c chown
	go generate cmd/capturer/gen.go


.PHONE: chown
chown:
	sudo chown -R ${USER}:${USER} ./cmd


# eBPFプログラム (written in C) をコンパイルしてオブジェクトファイル (eBPF bytecode) を生成する．
# syscallの引数にアクセスする際には，`-D__TARGET_ARCH_x86`をつける必要がある．レジスタ上の値の配置をx86のものに合わせるため．
capture_plain.o: ebpf_src/capture_plain.c
	cd ebpf_src && \
		clang -O2 -g -target bpf -D__TARGET_ARCH_x86 \
		-c capture_plain.c -o capture_plain.o


# コマンドライン引数においてデコード処理の並列度を指定する．
# バイナリは指定された数だけgoroutineを生成し，デコード処理を並列実行する．
# makeに対してコマンドライン引数を指定できるようにするべきだったが，手間がかかりそうだったのでやらなかった．
.PHONY: run
run: all
# if [ -f /data_shelter/${TMP_FILE_NAME} ]; then \
# 	sudo rm /data_shelter/${TMP_FILE_NAME}; \
# 	echo "Remove old file"; \
# fi
# sudo ./${BIN} ${TMP_FILE_NAME} 4
	sudo ./${BIN} 4

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
	sudo rm -rf /data_shelter/*
