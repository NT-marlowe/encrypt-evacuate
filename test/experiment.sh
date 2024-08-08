#!/bin/bash
#

set -x


PROGRAM=ebpf-ssl
for file in $(ls ./data); do
    # echo ${PROGRAM}
    sudo ../${PROGRAM} &
    pid=$!
    sleep 1
    sudo kill -9 $pid > /dev/null 2>&1
    echo "Killed $pid"
done
