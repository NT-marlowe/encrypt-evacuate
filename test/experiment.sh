#!/bin/bash
#

set -e


EBPF_PROGRAM=ebpf-ssl
DATA_SHELTER=/data_shelter

sudo rm -f ${DATA_SHELTER}/*

cd .. && make && cd test

# for file in $(ls ./data | grep -v enc); do
for file in $(ls ./data/1* | grep -v enc); do
    file=$(basename $file)
    sudo ../${EBPF_PROGRAM} $file &
    pid=$!
    sleep 1
    
    ./my_simple_ransomware ./data/$file
    echo "Ransomware ran on $file"

    sleep 3

    sudo kill -SIGINT $pid > /dev/null 2>&1
    echo "Killed $pid"
    
    sleep 1
    echo -----------
    
    # rm ./data/{$file}.enc
done


