#!/bin/bash
#
set -e


echo "filename, levenstein_distance, partial_match_ratio"
for file in $(ls ./data | grep -v enc); do
    echo -n "$file, "

    original=./data/$file
    recovered=/usr/tmp/data_shelter/$file
    python3 commands.py calc ${original} ${recovered}
done


