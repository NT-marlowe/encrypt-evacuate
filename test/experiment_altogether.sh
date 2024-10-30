#!/bin/bash
#
# This script run the capturer once, and the ransomware program over all files.

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

../${EBPF_PROGRAM} ${parallelism} &
pid=$!
sleep 1

# for file in $(ls ./data | grep -v enc); do
# for file in $(ls ./data/1* | grep -v enc); do
for file in $(ls ./data/2* | grep -v enc); do
    # file=$(basename $file)
    
    ./my_simple_ransomware ${file}
    echo "Ransomware ran on ${file}"

    sleep 2

    echo -----------
    
done

kill -SIGINT $pid > /dev/null 2>&1
echo "Killed $pid"
    
chown -R ${USER}:${USER} ../cmd/*
chown -R ${USER}:${USER} ../${EBPF_PROGRAM}
chown -R ${USER}:${USER} ./data/*

rm ./data/*.enc

echo "parallelism: ${parallelism}"

