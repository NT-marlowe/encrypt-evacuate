#!/bin/bash

set -e
set -u


# echo "filename, retention_rate"

ringbuf_filesize=$1

# for file in $(ls ./data | grep -v enc); do
for file in $(ls ./data/1* | grep -v enc); do
    file=$(basename $file)
    echo -n "$file, "
    original=./data/$file
    recovered=/usr/tmp/data_shelter/$file
    
    recovered=/usr/tmp/data_shelter/${ringbuf_filesize}MiB_$file
    
    python3 commands.py reten ${original} ${recovered}
done
