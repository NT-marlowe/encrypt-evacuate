#!/usr/bin/env bash
#
set -eu

time_len=30
EBPF_PROGRAM=ebpf-ssl

function encrypt_sysbench_files() {
    for file in $(ls ./files_sysbench | grep -v 'enc'); do
        echo "Encrypting ${file}"
        ./my_simple_ransomware ./files_sysbench/${file}
    done
}

function clear_cache() {
    sudo sync
    echo 3 | sudo tee /proc/sys/vm/drop_caches
}

function run_once() {
    local filename=$1
    echo $filename

    clear_cache

    iostat -xyt -o JSON 1 ${time_len} >${filename} &
    pid=$!

    encrypt_sysbench_files

    echo "Cleaning up..."
    wait $pid
    rm -rf ./files_sysbench/*.enc

}

function build() {
    cd $(git rev-parse --show-toplevel)
    make all
    make clear-shelter
    cd test
}

function main() {
    if [ "$#" -lt 2 ]; then
        clear_cache
        encrypt_sysbench_files
        return
    fi

    out_filename=$1
    iter=$2

    local parallelism
    if [ "$#" -eq 3 ]; then
        if [ "$EUID" -ne 0 ]; then
            echo "Please run as root"
            exit
        fi

        echo "Building..."
        parallelism=$3
        build
    fi

    for i in $(seq 1 ${iter}); do
        local copied_pid

        if [ "$#" -eq 3 ]; then
            ../${EBPF_PROGRAM} ${parallelism} &
            local pid=$!
            copied_pid=${pid}
            sleep 1
        fi
        run_once ${out_filename}.${i}

        sleep 2

        if [ "$#" -eq 3 ]; then
            kill -SIGINT ${copied_pid} >/dev/null 2>&1
            sleep 2
            rm -rf /data_shelter/files_sysbench
        fi

        # rm ./files_sysbench/*.enc
    done
}

main "$@"
