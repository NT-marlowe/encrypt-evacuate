#!/bin/bash
#

set -e


EBPF_PROGRAM=ebpf-ssl
DATA_SHELTER=/usr/tmp/data_shelter

sudo rm -f ${DATA_SHELTER}/*

for file in $(ls ./data | grep -v enc); do
    sudo ../${EBPF_PROGRAM} $file &
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


