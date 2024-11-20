#!/bin/bash
#
# This script run the capturer once, and the ransomware program over all files.

set -e
set -u

# if not root user, exit
if [ "$EUID" -ne 0 ]; then
    echo "Please run as root"
    exit
fi

USER=marlowe
EBPF_PROGRAM=ebpf-ssl

parallelism=$1

cd $(git rev-parse --show-toplevel)
make all
make clear-shelter
cd test

../${EBPF_PROGRAM} ${parallelism} &
pid=$!
sleep 1

# for file in $(ls ./data | grep -v enc); do
# for file in $(ls ./data/incremental/* | grep -v enc); do
for file in $(ls ./data/power_10/* | grep -v enc); do
    # for file in $(ls ./data/power_2/* | grep -v enc); do
    # for file in $(ls ./data/incremental/10_1MB.data | grep -v enc); do
    # for file in $(ls ./data/2* | grep -v enc); do
    # file=$(basename $file)

    ./my_simple_ransomware ${file}
    echo "Ransomware ran on ${file}"

    # sleep 4

    echo -----------

done

kill -SIGINT $pid >/dev/null 2>&1
echo "Killed $pid"

chown -R ${USER}:${USER} ../cmd/*
chown -R ${USER}:${USER} ../${EBPF_PROGRAM}
chown -R ${USER}:${USER} ./data/*

# rm ./data/*.enc

echo "parallelism: ${parallelism}"
