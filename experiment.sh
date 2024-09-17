#!/bin/bash
#

set -e
set -u


EBPF_PROGRAM=ebpf-ssl
DATA_SHELTER=/usr/tmp/data_shelter

# sudo rm -f ${DATA_SHELTER}/*
ringbuf_size=$1

make all

cd test
# for file in $(ls ./data | grep -v enc); do
for file in $(ls ./data/2* | grep -v enc); do
    file=$(basename $file)
    sudo ../${EBPF_PROGRAM} ${ringbuf_size}MiB_$file &
    pid=$!
    sleep 1
    
    ./my_simple_ransomware ./data/$file -e
    echo "Ransomware ran on $file"

    sleep 3

    sudo kill -SIGINT $pid > /dev/null 2>&1
    echo "Killed $pid"
    
    sleep 1
    echo -----------
done


