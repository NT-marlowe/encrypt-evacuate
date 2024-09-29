#!/bin/bash
#
set -e
set -u


# echo "filename, levenstein_distance, partial_match_ratio"
# echo "filename, match_ratio"
# for file in $(ls ./data | grep -v enc); do

subcom=$1
# ringbuf_filesize=$2

# for file in $(ls ./data/2* | grep -v enc); do
for file in $(ls ./data/1* | grep -v enc); do
    file=$(basename $file)
    echo -n "$file, "

    original=./data/$file
    # recovered=/data_shelter/${ringbuf_filesize}MiB_$file
    recovered=/data_shelter/$file
    python3 commands.py ${subcom} ${original} ${recovered}
done


