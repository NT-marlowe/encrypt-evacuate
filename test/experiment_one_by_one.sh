#!/bin/bash
#
# This script run the capturer and the ransomware program one by one.

set -e
set -u


# if not root user, exit
if [ "$EUID" -ne 0 ]
  then echo "Please run as root"
  exit
fi

USER=marlowe
EBPF_PROGRAM=ebpf-ssl
DATA_SHELTER=/data_shelter

parallelism=$1

rm -rf ${DATA_SHELTER}/*
cd .. && make all && cd test

# for file in $(ls ./data | grep -v enc); do
# for file in $(ls ./data/1* | grep -v enc); do
for file in $(ls ./data/2* | grep -v enc); do
    # file=$(basename $file)
    ../${EBPF_PROGRAM} ${parallelism} &
    pid=$!
    sleep 1
    
    ./my_simple_ransomware ${file}
    echo "Ransomware ran on ${file}"

    sleep 3

    kill -SIGINT $pid > /dev/null 2>&1
    echo "Killed $pid"
    
    sleep 3
    echo -----------
    
    # rm ./data/{$file}.enc
done

chown -R ${USER}:${USER} ../cmd/*
chown -R ${USER}:${USER} ../${EBPF_PROGRAM}
chown -R ${USER}:${USER} ./data/*

echo "parallelism: ${parallelism}"

