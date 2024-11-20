#!/usr/bin/env bash
#
set -eu
out_filename=$1
iter=$2

time_len=10
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
        echo $i
        local pid

        if [ "$#" -eq 3 ]; then
            echo $parallelism
            ../${EBPF_PROGRAM} ${parallelism} >stdout.log 2>stderr.log &
            pid=$!
            copied_pid=${pid}
            echo "pid is ${pid}"
            # search process with $pid
            ps auxww | grep ${pid} | grep -v grep
            sleep 1
        fi
        run_once ${out_filename}.${i}

        sleep 2

        if [ "$#" -eq 3 ]; then
            echo ${pid}
            # kill -SIGINT ${pid} >/dev/null 2>&1
            kill -SIGINT ${copied_pid}
            sleep 2
        fi
    done
}

main "$@"
