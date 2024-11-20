#!/usr/bin/env bash
#
set -eu
out_filename=$1
iter=$2

time_len=30
EBPF_PROGRAM=ebpf-ssl

function run_once() {
    local filename=$1
    echo $filename

    sudo sync
    echo 3 | sudo tee /proc/sys/vm/drop_caches

    iostat -xyt -o JSON 1 ${time_len} >${filename} &
    pid=$!

    for file in $(ls ./files_sysbench | grep -v 'enc'); do
        echo "Encrypting ${file}"
        ./my_simple_ransomware ./files_sysbench/${file}
    done

    echo "Cleaning up..."
    rm -rf ./files_sysbench/*.enc
    wait $pid
}

function build() {
    cd $(git rev-parse --show-toplevel)
    make all
    make clear-shelter
    cd test
}

function main() {
    if [ "$#" -eq 3 ]; then
        echo "Building..."
        parallelism=$3
        build
    fi

    for i in $(seq 1 ${iter}); do
        if [ "$#" -eq 3 ]; then
            ../${EBPF_PROGRAM} ${parallelism} &
            pid=$!
        fi
        run_once ${out_filename}.${i}

        sleep 2

        if [ "$#" -eq 3 ]; then
            kill -SIGINT ${pid} >/dev/null 2>&1
            sleep 2
        fi
    done
}

main "$@"
