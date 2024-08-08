#!/bin/bash
#
set -e


for file in $(ls ./data | grep -v enc); do
    echo "Comparing $file"

    original=./data/$file
    recovered=/usr/tmp/data_shelter/$file
    python3 commands.py calc ${original} ${recovered}
done


