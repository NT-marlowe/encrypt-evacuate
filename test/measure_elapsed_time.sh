#!/bin/bash
#
if [ "$(id -u)" -ne 0 ]; then
    echo "Please run as root"
    exit 1
fi


set -u
set -e
# set -x

function run_measure() {
    ../ebpf-ssl tmp &
    ebpf_pid=$!

    echo "ebpf start"

    sleep 1

    data_filename=$1
    base_filename=$(basename $data_filename .data)
    cat /sys/kernel/debug/tracing/trace_pipe \
        > ./result/ebpf_time/elapsed_time_${base_filename}.txt 2>/dev/null &
    cat_pid=$!

    echo "Start to measure elapsed time: ${cat_pid}"

    sleep 1

    for i in {1..10}; do
        ./my_simple_ransomware $data_filename -e
        sleep 1
    done

    kill $cat_pid

    sleep 0.5

    kill -SIGINT $ebpf_pid

    sleep 1

    rm /usr/tmp/data_shelter/tmp
    ps auxww | grep ebpf 
}

for file in $(ls ./data | grep -v enc | egrep '^0'); do
# file=$1
    echo "Start to measure elapsed time: $file"
    run_measure ./data/$file
    # run_measure $file
    
    sleep 3
    echo "------------"
done

