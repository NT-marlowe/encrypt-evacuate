#!/usr/bin/env bash
#
set -eu
out_filename=$1
iter=$2
time_len=30

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

function main() {
    for i in $(seq 1 ${iter}); do
        run_once ${out_filename}.${i}
        sleep 2
    done
}

main "$@"
