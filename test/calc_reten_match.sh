#!/bin/bash
#
set -e
set -u

subcom=$1

# echo "filename, levenstein_distance, partial_match_ratio"
echo "filename, ${subcom}"
# for original in $(ls ./data | grep -v enc); do

# ringbuf_filesize=$2

# for original in $(ls ./data/2* | grep -v enc); do
# for original in $(ls ./data/incremental/* | grep -v enc); do
# for original in $(ls ./data/power_10/* | grep -v enc); do
for original in $(ls ./files_sysbench/* | grep -v '\.enc'); do
    # for original in $(ls ./data/power_2/* | grep -v enc); do
    base=$(basename $original)
    echo -n "${base}, "

    # recovered=/data_shelter/${ringbuf_filesize}MiB_$original
    recovered=/data_shelter/${original}
    python3 commands.py ${subcom} ${original} ${recovered}
done
